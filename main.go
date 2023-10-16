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
	None      int = 0
	Essential     = 1
	Scenario      = 2
)

const (
	ArmstrongPowerhouse int = 0
	Steam                   = 1
	JustTrains              = 2
	FastlineSimulation      = 3
	ATS                     = 4
	Other                   = 5
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
		log.Fatal("Unable to establish db connection")
	}
	c := colly.NewCollector(colly.AllowedDomains("www.armstrongpowerhouse.com", "www.justtrains.net", "alanthomsonsim.com"))

	sitemapRegex, _ := regexp.Compile(`^https://www.armstrongpowerhouse.com/(enhancements|rolling_stock|routes|scenarios|sounds)/?.*$`)
	addToCartRegex, _ := regexp.Compile(`^\s*addToCart\(\'(\d*)\'\);$`)
	requirementsRegex, err := regexp.Compile(`^\s*(?:AP|DTG|ATS|JT|Fastline Simulation) (.*) -[ \xa0]More Information`)
	urlRegex, _ := regexp.Compile(`(www.armstrongpowerhouse.com|store.steampowered.com|www.justtrains.net|www.fastline-simulation.co.uk|sites.fastspring.com|alanthomsonsim.com)`)
	apProdIdRegex, _ := regexp.Compile(`product_id=(\d*)`)
	steamUrlRegex, _ := regexp.Compile(`^https?://store.steampowered.com/(app/\d*)`)
	atsUrlRegex, _ := regexp.Compile(`https://alanthomsonsim.com/(product/[^/]*)`)
	jtUrlRegex, _ := regexp.Compile("https?://www.justtrains.net/(product/[^/]*)")
	jtUrlPriceRegex, _ := regexp.Compile(`^US\$(\d*\.\d*)$`)

	if err != nil {
		log.Fatal(err.Error())
	}
	c.OnHTML(".sitemap-info a[href]", func(e *colly.HTMLElement) {
		if sitemapRegex.MatchString(e.Attr("href")) {
			e.Request.Visit(e.Attr("href"))
		}
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
			_, err = db.Exec("UPDATE Product SET Current_Price = $1 WHERE URL = $2", jtPrice, jtUrl)
			if err != nil {
				log.Fatal(err)
			}

		}
	})

	// ATS Site Price Parsing
	c.OnHTML(".product_title+.price>span>bdi", func(e *colly.HTMLElement) {
		atsPrice := strings.ReplaceAll(e.Text, ".", "")[2:]

		atsUrl := atsUrlRegex.FindStringSubmatch(e.Request.URL.String())[1]
		_, err = db.Exec("UPDATE Product SET Current_Price = $1 WHERE URL = $2", atsPrice, atsUrl)
		if err != nil {
			log.Fatal(err)
		}

	})

	c.OnHTML(".product-list", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, e *colly.HTMLElement) {
			priceRaw := e.ChildText(".price")
			if len(priceRaw) > 0 {
				name := e.ChildText(".name")
				price, err := decimal.NewFromString(priceRaw[2:])
				if err != nil {
					log.Fatal("Unable to parse addon price")
				}
				urlArray := strings.Split(e.ChildAttr("a", "href"), "/")
				url := urlArray[len(urlArray)-1]
				idRaw := e.ChildAttr(".controls>.cart>a", "onclick")
				id, _ := strconv.Atoi(addToCartRegex.FindStringSubmatch(idRaw)[1])

				products = append(products, Product{id, name, price, url})
				_, err = db.Exec(`
				INSERT INTO Product (ProductID, Name, URL, Current_Price, Company) 
				VALUES($1, $2, $3, $4, $5)`, id, name, url, price.Mul(decimal.NewFromInt32(100)), ArmstrongPowerhouse)
				if err != nil {
					log.Fatal("DB ERROR: " + err.Error())
				}

				productPageUrls = append(productPageUrls, homepage+"/"+url)
			}
		})
	})

	c.OnHTML(".product-info", func(e *colly.HTMLElement) {
		fixedUrlArray := strings.Split(e.Request.URL.String(), "/")
		fixedUrl := fixedUrlArray[len(fixedUrlArray)-1]
		result := db.QueryRow("SELECT ROWID FROM Product WHERE URL = ?", fixedUrl)
		var rowIndex int
		if err = result.Scan(&rowIndex); err != nil {
			log.Fatal("Unable to scan row")
		}
		requirementType := None
		e.ForEach("p", func(_ int, e *colly.HTMLElement) {
			if e.ChildText("u>b") == "Essential Requirements" {
				requirementType = Essential
			} else if e.ChildText("u>b") == "Scenario Requirements" {
				requirementType = Scenario
			}
			var name string
			var href string
			if len(e.ChildAttr("b>a,strong>a,a:has(strong),a:has(b)", "href")) > 0 && requirementsRegex.MatchString(e.Text) {
				nameRaw := e.Text
				name = requirementsRegex.FindStringSubmatch(nameRaw)[1]
				href = e.ChildAttr("b>a,strong>a,a:has(strong),a:has(b)", "href")
				sourcetype := -1
				source := urlRegex.FindStringSubmatch(href)
				if len(source) > 1 {
					switch source[1] {
					case "www.armstrongpowerhouse.com":
						sourcetype = ArmstrongPowerhouse
					case "store.steampowered.com":
						sourcetype = Steam
					case "www.justtrains.net":
						sourcetype = JustTrains
					case "sites.fastspring.com":
						sourcetype = FastlineSimulation
					case "www.fastline-simulation.co.uk":
						sourcetype = FastlineSimulation
					case "alanthomsonsim.com":
						sourcetype = ATS
					}
				}

				var requirementID int
				if sourcetype == ArmstrongPowerhouse {
					matchResult := apProdIdRegex.FindStringSubmatch(href)
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
						result := db.QueryRow("SELECT ROWID FROM Product WHERE ProductID = ?", productId)
						if err = result.Scan(&requirementID); err != nil {
							log.Fatal("Unable to scan row")
						}

					} else {
						productUrlArray := strings.Split(href, "/")
						productUrl := productUrlArray[len(productUrlArray)-1]
						if productUrl == "fsa-fta-wagon-pack" { // BUG ON THE AP SITE Class 317 Vol 1
							productUrl = strings.Replace(productUrl, "-", "_", -1)
						}
						productUrl = strings.Split(productUrl, "?")[0] // Other link bug ica d wagon pack for tda d link
						result := db.QueryRow("SELECT ROWID FROM Product WHERE URL = ?", productUrl)
						if err = result.Scan(&requirementID); err != nil {
							log.Fatal("Unable to scan row")
						}
					}
				} else {
					if steamUrlRegex.MatchString(href) {
						href = steamUrlRegex.FindStringSubmatch(href)[1]
					}
					if jtUrlRegex.MatchString(href) {
						href = jtUrlRegex.FindStringSubmatch(href)[1]
					}
					if atsUrlRegex.MatchString(href) {
						href = atsUrlRegex.FindStringSubmatch(href)[1]
					}
					result = db.QueryRow("SELECT ROWID from Product where Name = $1 OR ($2 <> 3 AND URL = $3);", name, sourcetype, href)
					if err = result.Scan(&requirementID); err != nil {
						_, err = db.Exec("INSERT INTO Product (Name, URL, Company) VALUES(?, ?, ?);", name, href, sourcetype)
						if err != nil {
							log.Fatal("Unable to insert new product")
						}
						result = db.QueryRow("SELECT last_insert_rowid()")
						err = result.Scan(&requirementID)
						if err != nil {
							log.Fatal("Unable to get recently inserted id")
						}
					}
				}
				if requirementType == Essential {
					_, err = db.Exec("INSERT INTO EssentialJoin (ProductID, EssentialID) VALUES(?, ?);", rowIndex, requirementID)
				} else if requirementType == Scenario {
					_, err = db.Exec("INSERT INTO ScenarioJoin (ProductID, ScenarioID) VALUES(?, ?);", rowIndex, requirementID)
				}
			}
		})
	})

	log.Println("Scraping AP site...")
	c.Visit(fmt.Sprintf("%s/index.php?route=information/sitemap", homepage))
	for _, url := range productPageUrls {
		c.Visit(url)
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
}

func reportDailyPrice(db *sql.DB) {
	log.Println("Creating reporting data...")
	atsRows, err := db.Query("SELECT Name, Current_Price FROM Product")
	if err != nil {
		log.Fatal("Error getting ats ids")
	}
	var name string
	var price int
	namePriceMap := make(map[string]int)
	for atsRows.Next() {
		err := atsRows.Scan(&name, &price)
		if err != nil {
			log.Fatal(err)
		}
		namePriceMap[name] = price
	}
	_, err = db.Exec("INSERT OR REPLACE INTO PriceReporting (Date) VALUES(CURRENT_DATE);")
	for k, v := range namePriceMap {
		if err != nil {
			log.Fatal(err)
		}
		query := fmt.Sprintf("UPDATE PriceReporting SET [%s] = ? WHERE [Date] = (CURRENT_DATE);", k)
		_, err := db.Exec(query, v)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// This function exists due to missing entries in IsThereAnyDeal api and fastline because it's a dead site now
func addMissingPrices(db *sql.DB) {
	log.Println("Adding missing prices...")
	db.Exec("UPDATE Product SET Current_Price = 999 WHERE URL = 'app/24083'")
	db.Exec("UPDATE Product SET Current_Price = 2499 WHERE URL = 'app/222554'")
	db.Exec("UPDATE Product SET Current_Price = 350 WHERE Name = '102t GLW Bogie Tanks'")
	db.Exec("UPDATE Product SET Current_Price = 350 WHERE Name = 'HEA Hoppers - Post BR'")
	db.Exec("UPDATE Product SET Current_Price = 350 WHERE Name = 'ZCA Sea Urchins'")
}

func fetchAtsPrices(c *colly.Collector, db *sql.DB) {
	log.Println("Retrieving ATS Prices...")
	atsRows, err := db.Query("SELECT URL FROM Product where Company = 4")
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
		c.Visit("https://alanthomsonsim.com/" + u)
	}
}

func fetchJtPrices(c *colly.Collector, db *sql.DB) {
	log.Println("Retrieving JT Prices...")
	jtRows, err := db.Query("SELECT URL FROM Product where Company = 2")
	if err != nil {
		log.Fatal("Error getting jt ids")
	}

	var jtUrl string
	var jtUrls []string
	for jtRows.Next() {
		err := jtRows.Scan(&jtUrl)
		if err != nil {
			log.Fatal(err)
		}
		jtUrls = append(jtUrls, jtUrl)
	}
	for _, u := range jtUrls {
		c.Visit("https://www.justtrains.net/" + u)
	}
}

func fetchSteamPrices(db *sql.DB) {
	log.Println("Retrieving Steam Prices...")
	isadKey := os.Getenv("isadKey")

	// HTTP endpoint
	steamRows, err := db.Query("SELECT URL FROM Product where Company = 1")
	if err != nil {
		log.Fatal("Error getting steam ids")
	}

	var steamid string
	var appids []string
	for steamRows.Next() {
		err := steamRows.Scan(&steamid)
		if err != nil {
			log.Fatal(err)
		}
		appids = append(appids, steamid)
	}

	plainsUrl := fmt.Sprintf("https://api.isthereanydeal.com/v01/game/plain/id/?key=%s&shop=steam&country=GB", isadKey)

	plainsBodyBytes, err := json.Marshal(&appids)
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
		log.Fatal(err)
	}
	type PlainsResponseType struct {
		Data map[string]string `json:"data"`
	}

	plainsResponseObj := PlainsResponseType{Data: make(map[string]string)}

	json.Unmarshal(plainsResponseBody, &plainsResponseObj)

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
		log.Fatal(err)
	}

	type ProductDetails struct {
		PriceNew float64 `json:"price_new"`
		PriceOld float64 `json:"price_old"`
		PriceCut float64 `json:"price_cut"`
		Url      float64 `json:"url"`
	}

	type ProductDetailsList struct {
		List []ProductDetails `json:"list"`
	}

	type PriceResponseType struct {
		Data map[string]ProductDetailsList `json:"data"`
	}

	//type ResponseType struct {
	//	Data map[string]string `json:"data"`
	//}
	//
	priceResponseObj := PriceResponseType{Data: make(map[string]ProductDetailsList)}

	json.Unmarshal(priceResponseBody, &priceResponseObj)

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
		db.Exec("UPDATE Product SET Current_Price = $1 WHERE URL = $2", v, k)
	}
}
