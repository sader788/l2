package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getElement(link string) ([]byte, error) {
	get, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func downloadSite(link string, maxDeep int) {
	parse, err := url.Parse(link)
	if err != nil {
		return
	}

	os.Mkdir(parse.Host, os.ModePerm)

	downloadSiteRec(link, 0, maxDeep)
}

func downloadSiteRec(link string, deep int, maxDeep int) {
	if maxDeep <= deep {
		return
	}

	data, err := getElement(link)
	if err != nil {
		return
	}

	links := getLinks(bytes.NewReader(data), link)

	for _, link := range links {
		err = os.MkdirAll(link.Host+link.Path, os.ModePerm)
		if err != nil {
			continue
		}
	}

	for _, link := range links {
		li := strings.LastIndex(link.Path, "/")
		if li < 0 {
			continue
		}

		el, err := getElement(link.String())
		if err != nil {
			continue
		}

		str := fmt.Sprintf("%s%s%s.html", link.Host, link.Path, link.Path[li:])

		writeToFile(el, str)

		downloadSiteRec(link.String(), deep+1, maxDeep)
	}

}

func getLinks(reader io.Reader, parentURL string) []*url.URL {
	tkn := html.NewTokenizer(reader)

	var vals []*url.URL
L:
	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			break L
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "a" || t.Data == "link" {
				attr := t.Attr
				for _, a := range attr {
					if a.Key == "href" {
						if !strings.HasPrefix(a.Val, `/`) {
							continue
						}
						parse, err := url.Parse(parentURL + a.Val)
						if err != nil {
							continue
						}

						vals = append(vals, parse)
						//fmt.Println(a.Val)
					}
				}
			}
		}
	}

	return vals
}

// Функция записывает данные в файл
func writeToFile(data []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	parse, err := url.Parse(`https://pingvinus.ru/`)
	if err != nil {
		return
	}
	fmt.Println(parse)

	downloadSite(`https://pingvinus.ru`, 2)

}
