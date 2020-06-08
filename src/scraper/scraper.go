package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper interface {
	DailyDealScraper(dealsLink string) ([]Deal, error)
	DealScraper(dealsLink string) ([]Deal, error)
}

// Deal is a single deal entry on ozbargain site
type Deal struct {
	Title   string
	NodeID  string
	Date    string
	Content string
	Link    string
}

// DealLinksGenerator Generate ozbargain Deals links upto page=5
func DealLinksGenerator(basicURL string) ([]string, error) {
	dealLinks := make([]string, 0)
	dealLinks = append(dealLinks, basicURL)
	for i := 1; i < 3; i++ {
		linkGen := basicURL + "?page" + "=" + strconv.Itoa(i)
		dealLinks = append(dealLinks, linkGen)
	}
	return dealLinks, nil
}

// to do daily scraper
// DailyDealScraper takes basic URL and current date to generate a list of Deals from the current date
func (s *Deal) DailyDealScraper(dealsLink string, Date string) ([]Deal, error) {
	dealLinks, _ := DealLinksGenerator(dealsLink)
	fmt.Println("dealLinks:", dealLinks)
	dailyDeals := make([]Deal, 0)
	for _, dealLinks := range dealLinks {
		// fmt.Println(dealLinks)
		bufferDeals, _ := s.DealScraper(dealLinks)
		// fmt.Println("bufferDeals:", bufferDeals)
		bufferDealDate, _ := FlatDealAttrs(bufferDeals, "Date")
		// fmt.Println("bufferDealDate:", bufferDealDate)
		uniqueDateInDeals, _ := uniqueSlice(bufferDealDate)
		if len(uniqueDateInDeals) == 1 && stringInSlice(Date, uniqueDateInDeals) {
			dailyDeals = append(dailyDeals, bufferDeals...)
		} else if len(uniqueDateInDeals) != 1 && stringInSlice(Date, uniqueDateInDeals) {
			for _, deal := range bufferDeals {
				if deal.Date == Date {
					dailyDeals = append(dailyDeals, deal)
				}
			}
		}
		// fmt.Println(uniqueDateInDeals)
		fmt.Println(dailyDeals, len(dailyDeals))
	}
	return dailyDeals, nil
}

// FlatDealAttrs take Deal list and flat the given attributes and return a flatten list of deals attributes
func FlatDealAttrs(deals []Deal, baseAttr string) ([]string, error) {
	flatAttrs := make([]string, 0)
	for _, deal := range deals {
		switch baseAttr := baseAttr; baseAttr {
		case "Title":
			flatAttrs = append(flatAttrs, deal.Title)
		case "NodeID":
			flatAttrs = append(flatAttrs, deal.NodeID)
		case "Date":
			flatAttrs = append(flatAttrs, deal.Date)
		case "Content":
			flatAttrs = append(flatAttrs, deal.Content)
		case "Link":
			flatAttrs = append(flatAttrs, deal.Link)
		}
	}
	return flatAttrs, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func uniqueSlice(slice []string) ([]string, error) {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list, nil
}

// DealScraper returns a list of Deals from a dealsLinks. bool return true or false, true means only unique dates exist in the Date values in the Deals slice, 1 means vice versa.
func (s *Deal) DealScraper(dealsLink string) ([]Deal, error) {
	deals := make([]Deal, 0)
	c := colly.NewCollector()

	re := regexp.MustCompile(`\d{2}/\d{2}/\d{4}`)

	c.OnHTML(".node", func(e *colly.HTMLElement) {
		deals = append(deals, Deal{
			Title:   e.ChildText(".title"),
			NodeID:  strings.Replace(e.ChildAttr("a", "href"), "/goto/", "", -1),
			Date:    re.FindString(e.ChildText(".submitted")),
			Content: e.ChildText(".content"),
			Link:    e.ChildAttr("a", "title"),
		})
	})
	c.Visit(dealsLink)
	// to do return bool unique list value
	return deals, nil
}
