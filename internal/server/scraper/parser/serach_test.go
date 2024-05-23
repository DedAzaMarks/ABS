package parser

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)
import _ "embed"

//go:embed "test/Фильмы скачать торрент в хорошем качестве.html"
var src []byte

func TestParseSearch(t *testing.T) {
	buf := bytes.NewReader(src)
	titles, err := ParseSearch(io.NopCloser(buf))
	require.NoError(t, err)
	require.Equal(t, []SearchResult{
		{
			Title: "Как приручить дракона (2010)",
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona.htm"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona-3.htm",
			Title: "Как приручить дракона 3 (2019)"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona-2.htm",
			Title: "Как приручить дракона 2 (2014)"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona-dilogija.htm",
			Title: "Как приручить дракона: Дилогия (2010)"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona-vozvrashchenie-domoy.htm",
			Title: "Как приручить дракона: Возвращение (2019)"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kak-priruchit-drakona-zhurnal-snogltoga.htm",
			Title: "Как приручить дракона: Журнал Сноглтога (2019)"},
		{
			Href:  "https://3b5a02883www.lafa.site/multfilm/Zarubezhnie-multfilmi/kniga-drakonov.htm",
			Title: "	 Как приручить дракона: Книга драконов (2011)"},
	},
		titles,
	)
}
