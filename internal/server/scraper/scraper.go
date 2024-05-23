package scraper

import (
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser"
	"io"
	"net/http"
)

func Search(url string) ([]parser.SearchResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	res, err := parser.ParseSearch(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse search: %w", err)
	}
	return res, nil
}

func Film(url string) ([]parser.FilmResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	res, err := parser.ParseFilm(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse list: %w", err)
	}
	return res, nil
}
