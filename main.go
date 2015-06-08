// TODO: put valid levels to paragraphs
// TODO: put attributes to p-hs from attribute array
// TODO: put valid numeration to p-hs
// TODO: exclude spreadSheetTable shit from tables
// TODO: put into p-hs <span> elements with their attributes
// TODO: make valid groups
// TODO: implement dependencies

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/bikbah/go-php-serialize/phpserialize"
)

const (
	MaxDepth int = 3
)

var (
	keysCount map[string]int
	r         *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	// Read php serialized file
	dat, err := ioutil.ReadFile("2")
	if err != nil {
		log.Fatalf("Read file error: %v", err)
	}

	// Try to deserialize data
	res, err := phpserialize.Decode(string(dat))
	if err != nil {
		log.Fatalf("Deserialize data error: %v", err)
	}

	// Parse raw data and get text from it
	m, ok := res.(map[interface{}]interface{})
	checkAssert(ok)
	convertedMap := parseRawMap(m)
	wholeText := getText(convertedMap, 0)
	saveData(wholeText)
	// fmt.Println(wholeText)

	// See key types
	// keysCount = make(map[string]int)
	// tmpFunc(convertedMap)
	// fmt.Println(keysCount)

}

// Recursive function parseRawMap parses raw map.
func parseRawMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	for k, v := range m {
		strKey := fmt.Sprintf("%v", k)

		switch reflect.TypeOf(v).String() {
		case "map[interface {}]interface {}":
			vv, ok := v.(map[interface{}]interface{})
			if !ok {
				log.Fatal("interface{} to map[interface{}]interface{} assertion fail..")
			}
			res[strKey] = parseRawMap(vv)
			continue
		case "string":
			res[strKey] = v
			continue
		default:
			log.Fatal("There is another type!!!")
		}

	}

	return res
}

// Recursive function getText gets whole text from e, parse it as HTML,
// add ID attribute to children paragraphs,
// then calls getText recursively to its children e.
// Returns all paragraphs as string.
func getText(m map[string]interface{}, depth int) (res string) {

	res = ""
	eID := ""
	paragraphDepth := depth

	// Get ID of e
	dollarMapRaw, ok := m["$"]
	if !ok {
		log.Fatal("No $ field of e element..")
	}
	dollarMap, ok := dollarMapRaw.(map[string]interface{})
	checkAssert(ok)
	if id, ok := dollarMap["id"]; ok {
		// res += fmt.Sprintf("id--->%v/ ", id)
		eID = fmt.Sprintf("%v", id)
	}

	// Check if e is some variant
	if _, ok := dollarMap["var"]; ok {
		paragraphDepth--
	}

	// Appending own texts first
	textMapRaw, ok := m["text"]
	if ok {
		textMap, ok := textMapRaw.(map[string]interface{})
		checkAssert(ok)
		textBody := ""
		for i := 0; i < len(textMap); i++ {
			textChildRaw, ok := textMap[strconv.Itoa(i)]
			if !ok {
				log.Fatal("getText error")
			}
			textChild, ok := textChildRaw.(map[string]interface{})
			checkAssert(ok)
			if textBodyRaw, ok := textChild["_"]; ok {
				textBody += fmt.Sprintf("%v", textBodyRaw)
			}
		}

		// Parse raw text and get HTML nodes
		rootNode, err := html.Parse(strings.NewReader(textBody))
		if err != nil {
			log.Fatalf("Parse html nodes error: %v", err)
		}
		htmlBody := rootNode.FirstChild.LastChild

		var buf bytes.Buffer
		paragraphID := eID
		for c := htmlBody.FirstChild; c != nil; c = c.NextSibling {
			c.Attr = append(c.Attr, html.Attribute{Key: "id", Val: paragraphID})
			c.Attr = append(c.Attr, html.Attribute{Key: "dd-level", Val: strconv.Itoa(paragraphDepth)})
			html.Render(&buf, c)
			buf.WriteString("\n")
			paragraphID = getRandomID()
		}

		res += buf.String()
		// res += textBody
	}

	// Then append texts of children
	childrenRaw, ok := m["e"]
	if !ok {
		return
	}
	children, ok := childrenRaw.(map[string]interface{})
	checkAssert(ok)
	for i := 0; i < len(children); i++ {
		childRaw, ok := children[strconv.Itoa(i)]
		if !ok {
			log.Fatal("getText error")
		}
		child, ok := childRaw.(map[string]interface{})
		checkAssert(ok)
		res += getText(child, depth+1)
	}
	return
}

// getRandomID returns random hex ID with length=8.
func getRandomID() string {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := ""
	for i := 0; i < 8; i++ {
		res += fmt.Sprintf("%x", r.Intn(16))
	}

	return "e" + strings.ToUpper(res)
}

func tmpFunc(m map[string]interface{}) {
	for _, v := range m {
		// log.Println(reflect.TypeOf(v).String())
		strKey := reflect.TypeOf(v).String()
		_, ok := keysCount[strKey]
		if ok {
			keysCount[strKey]++
		} else {
			keysCount[strKey] = 0
		}

		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			vv, ok := v.(map[string]interface{})
			if !ok {
				log.Fatal("interface{} to map[interface{}]interface{} assertion fail..")
			}
			tmpFunc(vv)
		}
	}
}

// checkAssert panics if not ok.
func checkAssert(ok bool) {
	if !ok {
		panic("assert")
		// log.Fatal("Assertion error..")
	}
}

// saveData wraps text with html tags and writes it to out.html file.
func saveData(text string) {
	html := `<html><head><meta charset="utf8"></head><body>` + text + `</body></html>`
	if err := ioutil.WriteFile("out.html", []byte(html), 0644); err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
