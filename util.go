package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// checkAssert panics if not ok.
func checkAssert(ok bool) {
	if !ok {
		panic("Assert error")
	}
}

// saveData wraps text with html tags and writes it to out.html file.
func saveData(text string) {
	html := `<html><head><meta charset="utf8"></head><body>` + text + `</body></html>`
	if err := ioutil.WriteFile("out.html", []byte(html), 0644); err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}

// getRandomID returns random hex ID with length=count.
func getRandomID(count int) string {
	res := ""
	for i := 0; i < count; i++ {
		res += fmt.Sprintf("%x", r.Intn(16))
	}

	return strings.ToUpper(res)
}
