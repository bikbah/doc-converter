// TODO: repcale <br> with <p>
// TODO: put valid numeration to p-hs
// TODO: exclude spreadSheetTable shit from tables
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

	rootE := E{}
	rootE.parse(res)
	log.Println(rootE.e[0].e[0].e[0].text)

	// Parse raw data and get text from it
	// m, ok := res.(map[interface{}]interface{})
	// checkAssert(ok)
	// convertedMap := parseRawMap(m)
	// wholeText := getText(convertedMap, 0)
	// saveData(wholeText)
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
	dollarMap, ok := getElemByKey("$", m)
	if !ok {
		log.Fatal("No $ field of e element..")
	}
	if id, ok := dollarMap["id"]; ok {
		// res += fmt.Sprintf("id--->%v/ ", id)
		eID = fmt.Sprintf("%v", id)
	}

	// Check if e is some variant
	if _, ok := dollarMap["var"]; ok || strings.Contains(eID, "linkContainer") {
		paragraphDepth--
	}

	// Get attributes of e
	attrMapStr := make(map[string]string)
	attrMap, ok := getElemByKey("attribute", m)
	if ok {
		for i := 0; i < len(attrMap); i++ {
			attr, ok := getElemByKey(strconv.Itoa(i), attrMap)
			if !ok {
				log.Fatal("Invalid index in attribute array of span")
			}
			attrDollar, ok := getElemByKey("$", attr)
			if !ok {
				log.Fatal("Attribute element has no $ field..")
			}
			attrName, ok := attrDollar["name"]
			if !ok {
				log.Fatal("Attribute element has no name..")
			}
			attrValue, ok := attrDollar["value"]
			if !ok {
				log.Fatal("Attribute element has no value..")
			}

			attrNameStr := fmt.Sprintf("%v", attrName)
			attrValueStr := fmt.Sprintf("%v", attrValue)
			attrMapStr[attrNameStr] = attrValueStr
		}
	}

	// Appending own texts first
	textMap, ok := getElemByKey("text", m)
	if ok {
		res += getParagraphText(textMap, eID, paragraphDepth, attrMapStr)
	}

	// Then append texts of children
	children, ok := getElemByKey("e", m)
	if !ok {
		return
	}
	for i := 0; i < len(children); i++ {
		child, ok := getElemByKey(strconv.Itoa(i), children)
		if !ok {
			log.Fatal("Invalid index in children e of e..")
		}
		res += getText(child, paragraphDepth+1)
	}
	return
}

//func getParagraphText(textMap map[string]interface{}, eID string, paragraphDepth int) string {
func getParagraphText(textMap map[string]interface{}, eID string, paragraphDepth int, attrMap map[string]string) string {

	textBody := ""
	for i := 0; i < len(textMap); i++ {
		textChild, ok := getElemByKey(strconv.Itoa(i), textMap)
		if !ok {
			log.Fatal("getText error")
		}
		if textBodyRaw, ok := textChild["_"]; ok {
			textBody += fmt.Sprintf("%v", textBodyRaw)
			continue
		}
		textBody += getSpan(textChild)
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

		styleDefault := make(map[string]string)
		styleCustom := make(map[string]string)
		for _, attr := range c.Attr {
			if attr.Key == "style" {
				styleDefault = parseCSS(attr.Val)
			}
		}
		for k, v := range attrMap {
			if k == "style" {
				styleCustom = parseCSS(v)
				for prop, val := range styleDefault {
					styleCustom[prop] = val
				}
				c.Attr = append(c.Attr,
					html.Attribute{
						Key: "style",
						Val: getCSS(styleCustom),
					})
				continue
			}
			c.Attr = append(c.Attr,
				html.Attribute{
					Key: k,
					Val: v,
				})
		}

		html.Render(&buf, c)
		buf.WriteString("\n")
		paragraphID = "e" + getRandomID(8)
	}

	return buf.String()
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

func getSpan(m map[string]interface{}) string {
	dollarMap, ok := getElemByKey("$", m)
	if !ok {
		return ""
	}
	if tag, ok := dollarMap["tag"]; !(ok && tag.(string) == "span") {
		return ""
	}

	attrMap, ok := getElemByKey("attribute", m)
	if !ok {
		log.Fatal(`Span element has no "attribute" array..`)
	}

	rootNode, err := html.Parse(strings.NewReader("<span></span>"))
	if err != nil {
		log.Fatalf("Parse span element error: %v", err)
	}
	spanNode := rootNode.FirstChild.LastChild.FirstChild
	attrMapStr := make(map[string]string)
	for i := 0; i < len(attrMap); i++ {
		attr, ok := getElemByKey(strconv.Itoa(i), attrMap)
		if !ok {
			log.Fatal("Invalid index in attribute array of span")
		}
		attrDollar, ok := getElemByKey("$", attr)
		if !ok {
			log.Fatal("Attribute element has no $ field..")
		}
		attrName, ok := attrDollar["name"]
		if !ok {
			log.Fatal("Attribute element has no name..")
		}
		attrValue, ok := attrDollar["value"]
		if !ok {
			log.Fatal("Attribute element has no value..")
		}

		attrNameStr := fmt.Sprintf("%v", attrName)
		attrNameStr = strings.ToLower(attrNameStr)
		attrValueStr := fmt.Sprintf("%v", attrValue)
		attrMapStr[attrNameStr] = attrValueStr
	}

	datafld, ok := attrMapStr["datafld"]
	if !ok {
		log.Fatal("Span has no datafld attribute..")
	}
	spanNode.Attr = append(spanNode.Attr,
		html.Attribute{
			Key: "dd-field",
			Val: datafld,
		})

	id, ok := attrMapStr["id"]
	if !ok {
		id = "df" + getRandomID(16)
	}
	spanNode.Attr = append(spanNode.Attr,
		html.Attribute{
			Key: "id",
			Val: id,
		})

	if source, ok := attrMapStr["source"]; ok {
		spanNode.Attr = append(spanNode.Attr,
			html.Attribute{
				Key: "dd-source",
				Val: source[1:],
			})
	}

	if hint, ok := attrMapStr["hint"]; ok {
		spanNode.Attr = append(spanNode.Attr,
			html.Attribute{
				Key: "dd-hint",
				Val: hint,
			})
	}

	if replText, ok := attrMapStr["replacementtext"]; ok {
		spanNode.Attr = append(spanNode.Attr,
			html.Attribute{
				Key: "dd-replacementtext",
				Val: replText,
			})
	}

	for k, v := range attrMapStr {
		if k == "rtl" || k == "id" || k == "datafld" || k == "source" || k == "hint" || k == "replacementtext" {
			continue
		}
		spanNode.Attr = append(spanNode.Attr,
			html.Attribute{
				Key: k,
				Val: v,
			})
	}

	var buf bytes.Buffer
	html.Render(&buf, spanNode)

	return buf.String()
}

// getElementByKey gets element by key from m
// and tries to assert it to map[string]interface{}.
func getElemByKey(key string, m map[string]interface{}) (map[string]interface{}, bool) {
	elemRaw, ok := m[key]
	if !ok {
		return nil, ok
	}

	elem, ok := elemRaw.(map[string]interface{})
	checkAssert(ok)
	return elem, ok
}
