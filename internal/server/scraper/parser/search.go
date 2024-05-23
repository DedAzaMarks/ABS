package parser

import (
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser/utils"
	"golang.org/x/net/html"
	"io"
	"log"
)

func traverseCTitle(node *html.Node, cTitles *[]*html.Node) {
	if node.Type == html.ElementNode && node.Data == "div" {
		class, ok := utils.GetAttr(node, "class")
		if ok && class == "c_title" {
			*cTitles = append(*cTitles, node)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseCTitle(c, cTitles)
	}
}

type SearchResult struct {
	Title string
	Href  string
}

func ParseSearch(r io.Reader) ([]SearchResult, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var searchResults []SearchResult

	var cTitles []*html.Node
	traverseCTitle(node, &cTitles)
	for _, title := range cTitles {
		a := title.FirstChild
		if a.Data != "a" {
			continue
		}
		var searchResult SearchResult
		for _, attr := range a.Attr {
			if attr.Key == "href" {
				searchResult.Href = attr.Val
			}
			if attr.Key == "alt" {
				searchResult.Title, err = utils.WIN2UTF(attr.Val)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if searchResult.Title != "" && searchResult.Href != "" {
			searchResults = append(searchResults, searchResult)
		}
	}
	return searchResults, nil
}
