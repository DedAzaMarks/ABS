package scraper

import (
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser"
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser/utils"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	URL := "%E3%E0%F0%F0%E8%20%EF%EE%F2%F2%E5%F0"
	utf8, err := url.QueryUnescape(URL)
	require.NoError(t, err)
	res, err := utils.WIN2UTF(utf8)
	require.NoError(t, err)
	log.Println(res)
}

func TestURL(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{
			"как приручить дракона",
			"https://3b5a02883www.lafa.site/torrentz/search/%EA%E0%EA%20%EF%F0%E8%F0%F3%F7%E8%F2%FC%20%E4%F0%E0%EA%EE%ED%E0/",
		},
		{
			"гарри поттер",
			"https://3b5a02883www.lafa.site/torrentz/search/%E3%E0%F0%F0%E8%20%EF%EE%F2%F2%E5%F0/",
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			win, err := utils.UTF2WIN(test.title)
			URLEncoded := url.QueryEscape(win)
			replacedSpaces := strings.ReplaceAll(URLEncoded, "+", "%20")
			require.NoError(t, err)
			URL := "https://3b5a02883www.lafa.site/torrentz/search/" + replacedSpaces + "/"
			require.Equal(t, test.expected, URL)
		})
	}
}

func TestSearch(t *testing.T) {
	win, err := utils.UTF2WIN("как приручить дракона")
	URLEncoded := url.QueryEscape(win)
	replacedSpaces := strings.ReplaceAll(URLEncoded, "+", "%20")
	requestURL := "https://3b5a02883www.lafa.site/torrentz/search/" + replacedSpaces + "/"
	resp, err := http.Get(requestURL)
	if err != nil {
		require.NoError(t, err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	res, err := parser.ParseSearch(resp.Body)
	require.NoError(t, err)
	log.Println(res)
}
