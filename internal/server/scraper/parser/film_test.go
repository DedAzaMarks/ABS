package parser

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"testing"
)
import _ "embed"

//go:embed "test/Скачать мультфильм Как приручить дракона через торрент бесплатно - Как приручить дракона скачать торрент хорошее качество.html"
var film []byte

func TestParseFilm(t *testing.T) {
	buf := bytes.NewReader(film)
	filmVariants, err := ParseFilm(io.NopCloser(buf))
	jsn, _ := json.Marshal(filmVariants)
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, jsn, "", " ")
	require.NoError(t, err)
	log.Println(prettyJSON.String())
}
