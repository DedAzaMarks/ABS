package parser

import (
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser/utils"
	"golang.org/x/net/html"
	"io"
	"log"
)

func traverseTable(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "table" {
		class, ok := utils.GetAttr(node, "class")
		if ok && class == "ts_film" {
			return node
		}
	}
	log.Print(node.Data)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if tbl := traverseTable(c); tbl != nil {
			return tbl
		}
	}
	return nil
}

func traverseLines(table *html.Node) []FilmResult {
	var results []FilmResult
	var c int
	for tr := table.FirstChild; tr != nil; tr = tr.NextSibling {
		if len(tr.Attr) != 0 {
			continue
		}
		var result FilmResult
		var err error
		td := tr.FirstChild.NextSibling
		for qualityTd := td.FirstChild.NextSibling; qualityTd != nil; qualityTd = qualityTd.NextSibling {
			if qualityTd.Type == html.ElementNode && qualityTd.Data == "span" {
				result.Quality, err = utils.WIN2UTF(qualityTd.FirstChild.Data)
				if err != nil {
					log.Println(err)
				}
				break
			}
		}
		td = td.NextSibling.NextSibling
		for translationTd := td.FirstChild; translationTd != nil; translationTd = translationTd.NextSibling {
			if translationTd.Type == html.ElementNode && translationTd.Data == "span" {
				result.TranslationVoiceover, err = utils.WIN2UTF(translationTd.FirstChild.Data)
				if err != nil {
					log.Println(err)
				}
			}
			if translationTd.Type == html.ElementNode && translationTd.Data == "i" {
				var author string
				author, err = utils.WIN2UTF(translationTd.FirstChild.Data)
				if err != nil {
					log.Println(err)
				}
				result.Author = author
			}
		}
		td = td.NextSibling.NextSibling
		if td.FirstChild != nil {
			result.FileFormat = td.FirstChild.Data
		}
		td = td.NextSibling.NextSibling
		result.Size = td.FirstChild.Data
		td = td.NextSibling.NextSibling
		for aTd := td.FirstChild; aTd != nil; aTd = aTd.NextSibling {
			if aTd.Type == html.ElementNode && aTd.Data == "a" {
				if href, ok := utils.GetAttr(aTd, "href"); ok {
					result.Magnet = href
					break
				}
			}
		}
		results = append(results, result)
		c++
	}
	return results
}

type FilmResult struct {
	Quality              string
	TranslationVoiceover string
	Author               string
	FileFormat           string
	Size                 string
	Magnet               string
}

func ParseFilm(r io.Reader) ([]FilmResult, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	table := traverseTable(node)
	if table == nil {
		return nil, myerrors.NotAFilm
	}
	print()
	tableBody := table.LastChild // ignore table header
	res := traverseLines(tableBody)
	return res, nil
}
