package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
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

	ticker := time.NewTicker(time.Minute)

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

	currency, err := db.GetCurrencyByCode(context.Background(), c.Code)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("Error fetching currency by code", err)
		}
	}

	if currency.Code == "" {
		cc, err := db.CreateCurrencies(context.Background(), database.CreateCurrenciesParams{
			ID:        uuid.New(),
			Name:      c.Name,
			Code:      c.Code,
			Price:     c.Price,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Fatal("Error creating new currency", err)
		}

		fmt.Println(cc)
		return
	}

	if c.Price == currency.Price {
		return
	}

	db.UpdateCurrencyPrice(context.Background(), database.UpdateCurrencyPriceParams{
		ID:    currency.ID,
		Price: currency.Price,
	})
	fmt.Println(currency)
}

func formatPrice(p string) string {
	price := strings.Replace(p, ",", "", 1)
	return strings.Replace(price, "$", "", 1)
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
				Price: formatPrice(price),
			})
		}
	})

	fmt.Println(currencies)
	return currencies
}
