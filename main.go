package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bikbah/go-php-serialize/phpserialize"
)

const (
	MaxDepth int = 3
)

var (
	keysCount map[string]int
)

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
	wholeText := getText(convertedMap)
	saveData(wholeText)
	// fmt.Println(wholeText)

	// See key types
	// keysCount = make(map[string]int)
	// tmpFunc(convertedMap)
	// fmt.Println(keysCount)

}

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

func getText(m map[string]interface{}) (res string) {

	res = ""
	eID := ""

	// Append ID of e
	dollarMapRaw, ok := m["$"]
	if !ok {
		log.Fatal("No $ field of e element..")
	}
	dollarMap, ok := dollarMapRaw.(map[string]interface{})
	checkAssert(ok)
	if id, ok := dollarMap["id"]; ok {
		// res += fmt.Sprintf("id--->%v/ ", id)
		eID = fmt.Sprintf("id--->%v/ ", id)
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

		// Parse raw text of each element
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(textBody))
		if err != nil {
			log.Fatalf("Getting goquery document from string error, %v", err)
		}
		// eID = eID
		selection := doc.SetAttr("id", eID)
		fmt.Println(selection.Html())
		// textBody = doc.Text()
		// if doc.Size() > 1 {
		// 	fmt.Println("More than one elem")
		// }

		res += textBody
	}

	// if res != "" {
	// 	fmt.Printf("%s\n\n", res)
	// }

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
		res += getText(child)
	}
	return
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

func checkAssert(ok bool) {
	if !ok {
		panic("assert")
		// log.Fatal("Assertion error..")
	}
}

func saveData(text string) {
	html := `<html><head><meta charset="utf8"></head><body>` + text + `</body></html>`
	if err := ioutil.WriteFile("out.html", []byte(html), 0644); err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
