package utils

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func UTF2WIN(utf8str string) (string, error) {
	decoder := charmap.Windows1251.NewEncoder()

	win1251str, _, err := transform.String(decoder, utf8str)
	if err != nil {
		fmt.Printf("Ошибка преобразования: %v\n", err)
		return "", err
	}
	return win1251str, nil
}
