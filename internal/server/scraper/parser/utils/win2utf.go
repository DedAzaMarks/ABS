package utils

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
)

func WIN2UTF(utf8Str string) (string, error) {
	encoder := charmap.Windows1251.NewDecoder()

	win1251Str, _, err := transform.String(encoder, utf8Str)
	if err != nil {
		log.Printf("Ошибка преобразования: %v\n", err)
		return "", err
	}
	return win1251Str, nil
}
