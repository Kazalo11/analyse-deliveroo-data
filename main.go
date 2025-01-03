package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type OrderSummary struct {
	Restaurant string
	Price      float64
	Date       time.Time
}

var (
	restaurants = make(map[string]int)
	total_price = 0.0
	r           = regexp.MustCompile(`\d\d.\d\d`)
)

func main() {
	file, err := os.Open("deliveroo_info.html")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	htmlBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	html := string(htmlBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	doc.Find("div.OrderSummary-c96f3428b2ccedb7 p.ccl-a396bc55704a9c8a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		_, ok := restaurants[text]
		if !ok {
			restaurants[text] = 1
		} else {
			restaurants[text] += 1
		}
	})

	doc.Find("div.OrderSummary-c96f3428b2ccedb7 p.ccl-6f43f9bb8ff2d712").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		numString := r.FindString(text)

		number, err := strconv.ParseFloat(numString, 64)
		if err != nil {
			fmt.Println("Can't convert string to number")
			return
		}
		total_price += number

	})

	type kv struct {
		Key   string
		Value int
	}

	var sortedRestaurants []kv
	for k, v := range restaurants {
		sortedRestaurants = append(sortedRestaurants, kv{k, v})
	}

	sort.Slice(sortedRestaurants, func(i, j int) bool {
		return sortedRestaurants[i].Value > sortedRestaurants[j].Value
	})

	fmt.Println("Restaurants sorted by frequency:")
	for _, entry := range sortedRestaurants {
		fmt.Printf("%s: %d\n", entry.Key, entry.Value)
	}

	fmt.Printf("Total amount spent: £%.2f\n", total_price)

}
