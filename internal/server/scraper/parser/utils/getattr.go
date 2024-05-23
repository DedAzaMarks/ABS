package utils

import "golang.org/x/net/html"

func GetAttr(n *html.Node, key string) (val string, ok bool) {
	for _, a := range n.Attr {
		if a.Key == key {
			val = a.Val
			ok = true
			return
		}
	}
	return
}
