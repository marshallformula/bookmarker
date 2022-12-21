package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type bookmark struct {
	name string
	url  string
}

func parse(bookmarkFile string) []bookmark {

	fileReader, err := os.Open(bookmarkFile)

	if err != nil {
		log.Fatalln("Fails reading bookmarkfile", err)
	}

	tokenizer := html.NewTokenizer(fileReader)
	foundAnchor := false
	var bookmarks []bookmark
	var b bookmark

	for {

		tt := tokenizer.Next()

		if !foundAnchor {
			b = bookmark{}
		}

		switch tt {

		case html.ErrorToken:

			return bookmarks

		case html.StartTagToken:

			token := tokenizer.Token()
			if token.Data == "a" {

				for _, attr := range token.Attr {
					if attr.Key == "href" {
						b.url = attr.Val
					}
				}

				foundAnchor = true

			}

		case html.TextToken:

			if foundAnchor {
				token := tokenizer.Token()
				// fmt.Println(token.Data)
				b.name = token.Data
				bookmarks = append(bookmarks, b)
				foundAnchor = false
			}

		}

	}

}

func main() {

	bookmarkFile := "/home/nate/Documents/bookmarks.html"
	bookmarks := parse(bookmarkFile)

	var output string

	for _, bm := range bookmarks {
		// rofi pango markup vomit on the `&` char so just display the url without any querystring params
		pathParts := strings.Split(bm.url, "?")
		output += fmt.Sprint("<b>", bm.name, "</b> <small><i><u>", pathParts[0], "</u></i></small>\n")
	}

	rofi_cmd := fmt.Sprint("printf '%s\n' \"", output, "\" | rofi -dmenu -i -format i -markup-rows")
	answer, err := exec.Command("bash", "-c", rofi_cmd).Output()

	if err != nil {
		log.Fatal(err)
	}

	answer_idx, err := strconv.Atoi(strings.TrimSpace(string(answer)))

	if err != nil {
		log.Fatal(err)
	}

	_, err = exec.Command("xdg-open", bookmarks[answer_idx].url).CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

}
