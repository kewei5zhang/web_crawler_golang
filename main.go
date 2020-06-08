package main

import "github.com/kewei5zhang/web_crawler_golang/src/scraper"

const ozBargainDealsLink = "https://www.ozbargain.com.au/deals"

func main() {
	scraper.DailyDealScraper(ozBargainDealsLink, "08/06/2020")
}
