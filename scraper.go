package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vitorqueirosz/gwallet/internal/database"
)

func startCryptoScrapper(db *database.Queries) {
	url := "https://coinmarketcap.com/"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making HTTP request:", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	ticker := time.NewTicker(time.Duration(10) * time.Second)

	for ; ; <-ticker.C {
		wg := &sync.WaitGroup{}
		currencies := cryptoScrapper(doc)

		for _, c := range currencies {
			wg.Add(1)

			go updateCurrency(db, wg, c)
		}
		wg.Wait()
	}
}

func updateCurrency(db *database.Queries, wg *sync.WaitGroup, c Currency) {
	defer wg.Done()
}

func cryptoScrapper(doc *goquery.Document) []Currency {
	currencies := make([]Currency, 0)

	doc.Find("table.cmc-table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// Extract currency name, code, and price for each <tr> element
		currencyName := s.Find("p.kKpPOn").Text()
		currencyCode := s.Find("p.coin-item-symbol").Text()
		price := s.Find("a.cmc-link > span").Text()

		if currencyName != "" {
			currencies = append(currencies, Currency{
				Name:  currencyName,
				Code:  currencyCode,
				Price: strings.Replace(price, ",", "", 1),
			})
		}
	})

	fmt.Println(currencies)
	return currencies
}
