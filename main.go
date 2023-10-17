package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	None        int = 0
	Essential       = 1
	Scenario        = 2
	Recommended     = 3
)

const (
	ArmstrongPowerhouse int = 0
	Steam                   = 1
	JustTrains              = 2
	FastlineSimulation      = 3
	ATS                     = 4
)

type Product struct {
	id    int
	name  string
	price decimal.Decimal
	link  string
}

var products []Product
var productPageUrls []string

const homepage = "https://www.armstrongpowerhouse.com"

func main() {
	log.Println("Starting")
	db, err := sql.Open("sqlite", "products.db")
	if err != nil {
		log.Fatal("Unable to establish db connection: " + err.Error())
	}

	log.Println("Building DB Tables")
	dbBuildScript, err := os.ReadFile("./CreateDb.sql")
	if err != nil {
		log.Fatal("Unable to read db build script: ", err)
	}
	if _, err := db.Exec(string(dbBuildScript)); err != nil {
		log.Fatal("Unable to execute db build script: ", err)
	}

	c := colly.NewCollector(colly.AllowedDomains("www.armstrongpowerhouse.com", "www.justtrains.net", "alanthomsonsim.com"))
	registerCollyHandlers(c, db)

	log.Println("Scraping AP site...")
	err = c.Visit(fmt.Sprintf("%s/index.php?route=information/sitemap", homepage))
	if err != nil {
		log.Fatal("Unable to visit ap sitemap: " + err.Error())
	}
	for _, pageUrl := range productPageUrls {
		err = c.Visit(pageUrl)
		if err != nil {
			log.Fatal("Unable to visit ap product page: " + pageUrl + "," + err.Error())
		}
	}

	var priceSum decimal.Decimal
	for _, p := range products {
		priceSum = priceSum.Add(p.price)
	}

	fetchSteamPrices(db)
	fetchJtPrices(c, db)
	fetchAtsPrices(c, db)
	addMissingPrices(db)
	reportDailyPrice(db)

	log.Println("Swapping DB Tables")
	dbSwapScript, err := os.ReadFile("./SwapTables.sql")
	if err != nil {
		log.Fatal("Unable to read db swap script: ", err)
	}
	if _, err := db.Exec(string(dbSwapScript)); err != nil {
		log.Fatal("Unable to execute db swap script: ", err)
	}
}

func registerCollyHandlers(c *colly.Collector, db *sql.DB) {
	sitemapRegex := regexp.MustCompile(`^https://www.armstrongpowerhouse.com/(enhancements|rolling_stock|routes|scenarios|sounds)/?.*$`)
	addToCartRegex := regexp.MustCompile(`^\s*addToCart\('(\d*)'\);$`)
	requirementsRegex := regexp.MustCompile(`^\s*(?:AP|DTG|ATS|JT|Fastline Simulation) (.*) -[ \xa0]More Information`)
	urlRegex := regexp.MustCompile(`(www.armstrongpowerhouse.com|store.steampowered.com|www.justtrains.net|www.fastline-simulation.co.uk|sites.fastspring.com|alanthomsonsim.com)`)
	apProdIdRegex := regexp.MustCompile(`product_id=(\d*)`)
	steamUrlRegex := regexp.MustCompile(`^https?://store.steampowered.com/(app/\d*)`)
	atsUrlRegex := regexp.MustCompile(`https://alanthomsonsim.com/(product/[^/]*)`)
	jtUrlRegex := regexp.MustCompile("https?://www.justtrains.net/(product/[^/]*)")
	jtUrlPriceRegex := regexp.MustCompile(`^US\$(\d*\.\d*)$`)

	// AP Sitemap Parsing
	c.OnHTML(".sitemap-info a[href]", func(e *colly.HTMLElement) {
		if sitemapRegex.MatchString(e.Attr("href")) {
			err := e.Request.Visit(e.Attr("href"))
			if err != nil {
				log.Fatalf("Unable to parse ap sitemap: %s\n", err.Error())
			}
		}
	})

	// AP Product List Page
	c.OnHTML(".product-list", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, e *colly.HTMLElement) {
			priceRaw := e.ChildText(".price")
			if len(priceRaw) > 0 {
				name := e.ChildText(".name")
				price, err := decimal.NewFromString(priceRaw[2:]) // Trim Pound sign
				if err != nil {
					log.Fatalf("Unable to parse addon price for %s\n", name)
				}
				urlArray := strings.Split(e.ChildAttr("a", "href"), "/")
				apProductUrl := urlArray[len(urlArray)-1]            // Last section of ap product url is saved as url identifier
				idRaw := e.ChildAttr(".controls>.cart>a", "onclick") // AP-assigned product id
				id, _ := strconv.Atoi(addToCartRegex.FindStringSubmatch(idRaw)[1])

				products = append(products, Product{id, name, price, apProductUrl})
				_, err = db.Exec(`
				INSERT INTO NewProduct (ProductID, Name, URL, Current_Price, Company) 
				VALUES($1, $2, $3, $4, $5)`, id, name, apProductUrl, price.Mul(decimal.NewFromInt32(100)), ArmstrongPowerhouse)
				if err != nil {
					log.Fatalf("Error adding ap product to db: %s,%s\n", name, err.Error())
				}

				productPageUrls = append(productPageUrls, homepage+"/"+apProductUrl)
			}
		})
	})

	// AP Product page
	c.OnHTML(".product-info", func(e *colly.HTMLElement) {
		// Convert to db-style url
		fixedUrlArray := strings.Split(e.Request.URL.String(), "/")
		fixedUrl := fixedUrlArray[len(fixedUrlArray)-1]

		productRowIds := db.QueryRow("SELECT ROWID FROM NewProduct WHERE URL = ?", fixedUrl)
		var rowIndex int
		if err := productRowIds.Scan(&rowIndex); err != nil {
			log.Fatal("Unable to scan row for getting ap product rowids: " + err.Error())
		}

		// Iterate through product requirements
		requirementType := None
		e.ForEach("p", func(_ int, e *colly.HTMLElement) {

			RequirementsHeader := e.ChildText("u>b")
			switch RequirementsHeader {
			case "Essential Requirements":
				requirementType = Essential
			case "Scenario Requirements":
				requirementType = Scenario
			case "Recommended Scenarios Requirement":
				requirementType = Recommended
			}

			// Set requirement name and href
			var requirementName string
			requirementHref := e.ChildAttr("b>a,strong>a,a:has(strong),a:has(b)", "href")
			if len(requirementHref) > 0 && requirementsRegex.MatchString(e.Text) {
				nameRaw := e.Text
				requirementName = requirementsRegex.FindStringSubmatch(nameRaw)[1]
				sourceType := -1
				source := urlRegex.FindStringSubmatch(requirementHref)
				if len(source) > 1 {
					switch source[1] {
					case "www.armstrongpowerhouse.com":
						sourceType = ArmstrongPowerhouse
					case "store.steampowered.com":
						sourceType = Steam
						requirementHref = steamUrlRegex.FindStringSubmatch(requirementHref)[1]
					case "www.justtrains.net":
						sourceType = JustTrains
						requirementHref = jtUrlRegex.FindStringSubmatch(requirementHref)[1]
					case "sites.fastspring.com", "www.fastline-simulation.co.uk":
						sourceType = FastlineSimulation
					case "alanthomsonsim.com":
						sourceType = ATS
						requirementHref = atsUrlRegex.FindStringSubmatch(requirementHref)[1]
					}
				}

				var requirementID int
				if sourceType == ArmstrongPowerhouse {
					// If possible, figure out ap product from productId in url
					matchResult := apProdIdRegex.FindStringSubmatch(requirementHref)
					if len(matchResult) > 0 {
						productId, err := strconv.Atoi(matchResult[1])
						if err != nil {
							log.Fatal("Unable to get ap product id from url")
						}
						if productId == 155 { // Fixes site bug on Class 37 Vol. 2 where old wherry is listed
							productId = 227
						}
						if productId == 197 { // Fixes site bug on Class 390 where old sky & weather is listed
							productId = 241
						}
						result := db.QueryRow("SELECT ROWID FROM NewProduct WHERE ProductID = ?", productId)
						if err = result.Scan(&requirementID); err != nil {
							log.Fatal("Unable to scan row: " + err.Error())
						}

						// If no id in url, attempt url product name match
					} else {
						productUrlArray := strings.Split(requirementHref, "/")
						productUrl := productUrlArray[len(productUrlArray)-1]
						if productUrl == "fsa-fta-wagon-pack" { // BUG ON THE AP SITE Class 317 Vol 1
							productUrl = strings.Replace(productUrl, "-", "_", -1)
						}
						productUrl = strings.Split(productUrl, "?")[0] // Other link bug ica d wagon pack for tda d link
						result := db.QueryRow("SELECT ROWID FROM NewProduct WHERE URL = ?", productUrl)
						if err := result.Scan(&requirementID); err != nil {
							log.Fatal("Unable to scan row: " + err.Error())
						}
					}
				} else { // If requirement is not made by AP
					result := db.QueryRow("SELECT ROWID from NewProduct where Name = $1 OR ($2 <> 3 AND URL = $3);", requirementName, sourceType, requirementHref)
					if err := result.Scan(&requirementID); err != nil {
						// If this requirement is not yet in the db, add it
						_, err = db.Exec("INSERT INTO NewProduct (Name, URL, Company) VALUES(?, ?, ?);", requirementName, requirementHref, sourceType)
						if err != nil {
							log.Fatalf("Unable to insert new requirement entry: %s,%s\n", requirementName, err.Error())
						}
						result = db.QueryRow("SELECT last_insert_rowid()")
						err = result.Scan(&requirementID)
						if err != nil {
							log.Fatalf("Unable to get recently inserted id: %s,%s\n", requirementName, err.Error())
						}
					}
				}
				switch requirementType {
				case Essential:
					_, err := db.Exec("INSERT INTO NewEssentialJoin (ProductID, EssentialID) VALUES(?, ?);", rowIndex, requirementID)
					if err != nil {
						log.Fatalf("Unable to insert requirment into Essentials table: %s,%s\n", requirementName, err.Error())
					}
				case Scenario:
					_, err := db.Exec("INSERT INTO NewScenarioJoin (ProductID, ScenarioID) VALUES(?, ?);", rowIndex, requirementID)
					if err != nil {
						log.Fatalf("Unable to insert requirment into Scenario table: %s,%s\n", requirementName, err.Error())
					}
				case Recommended:
					_, err := db.Exec("INSERT INTO NewRecommendedJoin (ProductID, ScenarioID) VALUES(?, ?);", rowIndex, requirementID)
					if err != nil {
						log.Fatalf("Unable to insert requirment into Recommended table: %s,%s\n", requirementName, err.Error())
					}
				}
			}
		})
	})

	// JT Site Price Parsing
	c.OnHTML("#lblPrice2", func(e *colly.HTMLElement) {
		if jtUrlPriceRegex.MatchString(e.Text) {
			jtPrice, err := decimal.NewFromString(jtUrlPriceRegex.FindStringSubmatch(e.Text)[1])
			if err != nil {
				log.Fatal("Unable to parse JT Price: ", err)
			}
			// TODO: Update with proper exchange value
			jtPrice = jtPrice.Mul(decimal.NewFromFloat(0.82)).Mul(decimal.NewFromInt32(100)).Truncate(0)

			jtUrl := jtUrlRegex.FindStringSubmatch(e.Request.URL.String())[1]
			_, err = db.Exec("UPDATE NewProduct SET Current_Price = $1 WHERE URL = $2", jtPrice, jtUrl)
			if err != nil {
				log.Fatalf("Unable to update JT Projuct Price: %s,%s\n", jtUrl, err.Error())
			}

		}
	})

	// ATS Site Price Parsing
	c.OnHTML(".product_title+.price>span>bdi", func(e *colly.HTMLElement) {
		atsPrice := strings.ReplaceAll(e.Text, ".", "")[2:]

		atsUrl := atsUrlRegex.FindStringSubmatch(e.Request.URL.String())[1]
		_, err := db.Exec("UPDATE NewProduct SET Current_Price = $1 WHERE URL = $2", atsPrice, atsUrl)
		if err != nil {
			log.Fatalf("Unable to parse ATS product price: %s, %s\n", atsUrl, err.Error())
		}
	})

}

func reportDailyPrice(db *sql.DB) {
	log.Println("Creating reporting data...")
	productRows, err := db.Query("SELECT Name, Current_Price FROM NewProduct")
	if err != nil {
		log.Fatal("Error getting product rows: ", err)
	}

	var name string
	var price int
	namePriceMap := make(map[string]int)
	for productRows.Next() {
		err := productRows.Scan(&name, &price)
		if err != nil {
			log.Fatal(err)
		}
		namePriceMap[name] = price
	}
	_, err = db.Exec("INSERT OR REPLACE INTO PriceReporting (Date) VALUES(CURRENT_DATE);")
	for k, v := range namePriceMap {
		if err != nil {
			log.Fatal("Error inserting date into reporting: ", err)
		}
		query := fmt.Sprintf("UPDATE PriceReporting SET [%s] = ? WHERE [Date] = (CURRENT_DATE);", k)
		_, err := db.Exec(query, v)
		if err != nil {
			log.Fatal("Error updating reporting: ", err)
		}
	}

}

// This function exists due to missing entries in IsThereAnyDeal api and fastline because it's a dead site now
func addMissingPrices(db *sql.DB) {
	log.Println("Adding missing prices...")
	_, err := db.Exec("UPDATE NewProduct SET Current_Price = 999 WHERE URL = 'app/24083'")
	if err != nil {
		log.Fatal("Unable to update app/24083 price")
	}
	_, err = db.Exec("UPDATE NewProduct SET Current_Price = 2499 WHERE URL = 'app/222554'")
	if err != nil {
		log.Fatal("Unable to update app/222554 price")
	}
	_, err = db.Exec("UPDATE NewProduct SET Current_Price = 350 WHERE Name = '102t GLW Bogie Tanks'")
	if err != nil {
		log.Fatal("Unable to update 102t GLW Bogie Tanks price")
	}
	_, err = db.Exec("UPDATE NewProduct SET Current_Price = 350 WHERE Name = 'HEA Hoppers - Post BR'")
	if err != nil {
		log.Fatal("Unable to update HEA Hoppers - Post BR price")
	}
	_, err = db.Exec("UPDATE NewProduct SET Current_Price = 350 WHERE Name = 'ZCA Sea Urchins'")
	if err != nil {
		log.Fatal("Unable to update ZCA Sea Urchins price")
	}
}

func fetchAtsPrices(c *colly.Collector, db *sql.DB) {
	log.Println("Retrieving ATS Prices...")
	atsRows, err := db.Query("SELECT URL FROM NewProduct where Company = 4")
	if err != nil {
		log.Fatal("Error getting ats ids")
	}

	var atsUrl string
	var atsUrls []string
	for atsRows.Next() {
		err := atsRows.Scan(&atsUrl)
		if err != nil {
			log.Fatal(err)
		}
		atsUrls = append(atsUrls, atsUrl)
	}
	for _, u := range atsUrls {
		err := c.Visit("https://alanthomsonsim.com/" + u)
		if err != nil {
			log.Fatalf("Unable to visit ats product site: %s, %s\n", u, err.Error())
		}
	}
}

func fetchJtPrices(c *colly.Collector, db *sql.DB) {
	log.Println("Retrieving JT Prices...")
	jtRows, err := db.Query("SELECT URL FROM NewProduct where Company = 2")
	if err != nil {
		log.Fatal("Error getting jt ids")
	}

	var jtUrl string
	var jtUrls []string
	for jtRows.Next() {
		err := jtRows.Scan(&jtUrl)
		if err != nil {
			log.Fatalf("Unable to scan JT row: %s, %s\n", jtUrl, err)
		}
		jtUrls = append(jtUrls, jtUrl)
	}
	for _, u := range jtUrls {
		err := c.Visit("https://www.justtrains.net/" + u)
		if err != nil {
			log.Fatalf("Error getting to visit jt product page: %s, %s\n", u, err.Error())
		}
	}
}

func fetchSteamPrices(db *sql.DB) {
	log.Println("Retrieving Steam Prices...")
	isadKey := os.Getenv("isadKey")
	plainsUrl := fmt.Sprintf("https://api.isthereanydeal.com/v01/game/plain/id/?key=%s&shop=steam&country=GB", isadKey)

	steamRows, err := db.Query("SELECT URL FROM NewProduct where Company = 1")
	if err != nil {
		log.Fatal("Error getting steam ids for adding prices: ", err)
	}

	// Get list of steam app ids
	var steamId string
	var appIds []string
	for steamRows.Next() {
		err := steamRows.Scan(&steamId)
		if err != nil {
			log.Fatal(err)
		}
		appIds = append(appIds, steamId)
	}

	plainsBodyBytes, err := json.Marshal(&appIds)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(plainsBodyBytes)

	// Make HTTP POST request
	plainsResp, err := http.Post(plainsUrl, "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := plainsResp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Read response body
	plainsResponseBody, err := io.ReadAll(plainsResp.Body)
	if err != nil {
		log.Fatal("Error reading plains HTTP response: ", err)
	}

	type PlainsResponseType struct {
		Data map[string]string `json:"data"`
	}

	plainsResponseObj := PlainsResponseType{Data: make(map[string]string)}

	err = json.Unmarshal(plainsResponseBody, &plainsResponseObj)
	if err != nil {
		log.Fatalf("Unable to unmarshal plains json response: %s\n", err.Error())
	}

	if plainsResp.StatusCode >= 400 && plainsResp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", plainsResp.StatusCode)
	}

	// Use the plains data to retrieve the price data
	var priceRequestParams []string
	for _, v := range plainsResponseObj.Data {
		priceRequestParams = append(priceRequestParams, v)
	}
	priceUrl := fmt.Sprintf("https://api.isthereanydeal.com/v01/game/prices/?key=%s&country=GB&plains=", isadKey)
	priceResp, err := http.Get(priceUrl + url.QueryEscape(strings.Join(priceRequestParams, ",")))
	if err != nil {
		log.Fatal(err)
	}

	// Close response body
	defer func() {
		err := priceResp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	priceResponseBody, err := io.ReadAll(priceResp.Body)
	if err != nil {
		log.Fatal("Unable to read steam price api response: ", err)
	}

	type ProductDetails struct {
		PriceNew float64 `json:"price_new"`
		PriceOld float64 `json:"price_old"`
		PriceCut float64 `json:"price_cut"`
		Url      string  `json:"url"`
	}

	type ProductDetailsList struct {
		List []ProductDetails `json:"list"`
	}

	type PriceResponseType struct {
		Data map[string]ProductDetailsList `json:"data"`
	}

	priceResponseObj := PriceResponseType{Data: make(map[string]ProductDetailsList)}
	err = json.Unmarshal(priceResponseBody, &priceResponseObj)
	if err != nil {
		log.Fatal("Unable to unmarshal steam price api data: ", err)
	}

	if priceResp.StatusCode >= 400 && priceResp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", priceResp.StatusCode)
	}

	urlPriceMap := make(map[string]decimal.Decimal)
	for k, v := range plainsResponseObj.Data {
		if len(priceResponseObj.Data[v].List) > 0 {
			urlPriceMap[k] = decimal.NewFromFloat(priceResponseObj.Data[v].List[0].PriceNew).Truncate(2).Mul(decimal.NewFromInt32(100))
		}
	}

	for k, v := range urlPriceMap {
		_, err := db.Exec("UPDATE NewProduct SET Current_Price = $1 WHERE URL = $2", v, k)
		if err != nil {
			log.Fatalf("Error setting steam price: %s,%s, %s\n", k, v, err.Error())
		}
	}
}
