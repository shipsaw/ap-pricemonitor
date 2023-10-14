package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
	"log"
	"regexp"
	"strconv"
)

const (
	None      int = 0
	Essential     = 1
	Scenario      = 2
)

type Product struct {
	id    int
	name  string
	price decimal.Decimal
	link  string
}

var products []Product
var tempConversionRate decimal.Decimal = decimal.NewFromFloat(0.82)

const homepage = "https://www.armstrongpowerhouse.com"

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.armstrongpowerhouse.com"))

	sitemapRegex, _ := regexp.Compile(`^https://www.armstrongpowerhouse.com/(enhancements|rolling_stock|routes|scenarios|sounds)/?.*$`)
	addToCartRegex, _ := regexp.Compile(`^\s*addToCart\(\'(\d*)\'\);$`)
	_ = sitemapRegex
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
				url := e.ChildAttr("a", "href")
				idRaw := e.ChildAttr(".controls>.cart>a", "onclick")
				fmt.Println("idRaw ", idRaw)
				id, _ := strconv.Atoi(addToCartRegex.FindStringSubmatch(idRaw)[1])
				fmt.Println("id ", id)

				products = append(products, Product{id, name, price, url})
			}
		})
	})

	/*
		c.OnHTML(".product-info", func(e *colly.HTMLElement) {
			e.ForEach("p", func(_ int, e *colly.HTMLElement) {
				requirementType := None
				if e.ChildText("u>b") == "Essential Requirements" {
					requirementType = Essential
				} else if e.ChildText("u>b") == "Scenario Requirements" {
					requirementType = Scenario
				}
				if len(e.ChildAttr("b>a", "href")) > 0 {
					name := e.Text
					href := e.ChildAttr("b>a", "href")

				}
			})
		})
	*/
	fmt.Println("Calculating...")
	c.Visit(fmt.Sprintf("%s/index.php?route=information/sitemap", homepage))

	//for _, p := range products {
	//	c.Visit(p.link)
	//}

	var priceSum decimal.Decimal
	for _, p := range products {
		fmt.Printf("Id: %v, Name: %s, Price: %v, Link: %v\n", p.id, p.name, p.price, p.link)
		priceSum = priceSum.Add(p.price)
	}
	fmt.Println("\nTotal cost in pounds: $" + priceSum.Truncate(2).String())
	fmt.Println("Total cost in dollars: $" + priceSum.Div(tempConversionRate).Truncate(2).String())

}
