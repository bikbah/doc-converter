package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

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

func getCSS(m map[string]string) string {
	res := ""
	for k, v := range m {
		res += k + ":" + v + ";"
	}
	return res
}

func parseCSS(s string) map[string]string {

	res := make(map[string]string)
	props := make([]string, 1)
	if strings.Contains(s, ";") {
		props = strings.Split(s, ";")
	} else {
		props[0] = s
	}
	for _, prop := range props {
		pair := strings.Split(prop, ":")
		if len(pair) == 2 {
			res[pair[0]] = pair[1]
		}
	}

	return res
}
