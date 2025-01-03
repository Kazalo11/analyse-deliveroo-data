package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type OrderSummary struct {
	Restaurant string
	Price      float64
	Date       time.Time
}

var (
	orderSummary []string
	summaryClass = "OrderSummary-c96f3428b2ccedb7"
)

func extractTextFromPTags(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "p" {
		text := getTextContent(n)
		if text != "" {
			orderSummary = append(orderSummary, text)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextFromPTags(c)
	}
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(getTextContent(c))
	}
	return sb.String()
}

func getOrderSummary(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, summaryClass) {

				extractTextFromPTags(n)
				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getOrderSummary(c)
	}
}

func main() {
	file, err := os.Open("deliveroo_info.html")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	getOrderSummary(doc)
	for _, text := range orderSummary {
		fmt.Println(text)
	}

}
