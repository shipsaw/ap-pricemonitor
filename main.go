package main

import (
	"database/sql"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
	"log"
	_ "modernc.org/sqlite"
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
	db, err := sql.Open("sqlite", "products.db")
	if err != nil {
		log.Fatal("Unable to establish db connection")
	}
	c := colly.NewCollector(colly.AllowedDomains("www.armstrongpowerhouse.com"))

	sitemapRegex, _ := regexp.Compile(`^https://www.armstrongpowerhouse.com/(enhancements|rolling_stock|routes|scenarios|sounds)/?.*$`)
	addToCartRegex, _ := regexp.Compile(`^\s*addToCart\(\'(\d*)\'\);$`)
	requirementsRegex, err := regexp.Compile(`^\s*(?:AP|DTG|ATS|JT|Fastline Simulation) (.*) -[ \xa0]More Information`)
	urlRegex, _ := regexp.Compile(`(www.armstrongpowerhouse.com|store.steampowered.com|www.justtrains.net|www.fastline-simulation.co.uk|sites.fastspring.com)`)
	apProdIdRegex, _ := regexp.Compile(`product_id=(\d*)`)
	steamUrlRegex, _ := regexp.Compile(`^https?://(store.steampowered.com/app/\d*)`)

	if err != nil {
		log.Fatal(err.Error())
	}
	c.OnHTML(".sitemap-info a[href]", func(e *colly.HTMLElement) {
		if sitemapRegex.MatchString(e.Attr("href")) {
			e.Request.Visit(e.Attr("href"))
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

	c.Visit(fmt.Sprintf("%s/index.php?route=information/sitemap", homepage))
	for _, url := range productPageUrls {
		c.Visit(url)
	}

	var priceSum decimal.Decimal
	for _, p := range products {
		priceSum = priceSum.Add(p.price)
	}
}
