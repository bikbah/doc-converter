package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type E struct {
	dollar    map[string]string
	attribute []Attribute
	e         []*E
	text      []Text
	level     int
}

func (e *E) parse(i interface{}) {
	m, ok := i.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in E.parse()")
	}

	for k, v := range m {
		strKey := fmt.Sprintf("%v", k)

		vv, ok := v.(map[interface{}]interface{})
		if !ok {
			panic("Assert error in E.parse()")
		}

		switch strKey {
		case "$":
			e.dollar = make(map[string]string)

			for dollarK, dollarV := range vv {
				dollarKStr := fmt.Sprint(dollarK)
				dollarVStr := fmt.Sprint(dollarV)

				e.dollar[dollarKStr] = dollarVStr
			}
			break
		case "attribute":
			for _, attrV := range vv {
				attr := Attribute{}
				attr.parse(attrV)
				e.addAttribute(attr)
			}
			break
		case "e":
			eMap := make(map[string]interface{})

			for eK, eV := range vv {
				eKStr := fmt.Sprint(eK)
				eMap[eKStr] = eV
			}

			for i := 0; i < len(eMap); i++ {
				eKey := fmt.Sprint(i)
				childInterface, ok := eMap[eKey]
				if !ok {
					panic("Invalid index in children e..")
				}

				childE := E{}
				childE.parse(childInterface)
				e.addChildE(&childE)
			}
			break

		case "text":
			textMap := make(map[string]interface{})

			for textK, textV := range vv {
				textKStr := fmt.Sprint(textK)
				textMap[textKStr] = textV
			}

			for i := 0; i < len(textMap); i++ {
				textKStr := fmt.Sprint(i)
				textInterface, ok := textMap[textKStr]
				if !ok {
					panic("Invalid index in text array of e..")
				}

				text := Text{}
				text.parse(textInterface)
				e.addText(text)
			}
			break
		}
	}

	e.setLevels()
}

func (e *E) setLevels() {
	childrenLevel := e.level + 1
	if _, ok := e.dollar["varcont"]; ok || strings.Contains(e.dollar["id"], "linkContainer") {
		childrenLevel--
	}

	for _, child := range e.e {
		child.level = childrenLevel
		child.setLevels()
	}
}

func (e *E) addAttribute(attr Attribute) {
	e.attribute = append(e.attribute, attr)
}

func (e *E) addText(text Text) {
	e.text = append(e.text, text)
}

func (e *E) addChildE(child *E) {
	e.e = append(e.e, child)
}

func (e *E) getText() string {
	attrMap := make(map[string]string)
	for _, attr := range e.attribute {
		attrMap[attr.name] = attr.value
	}

	textBody := ""
	for _, text := range e.text {
		textBody += text.getText()
	}

	rootNode, err := html.Parse(strings.NewReader(textBody))
	if err != nil {
		panic("Parse all text elements error")
	}

	var buf bytes.Buffer
	htmlBody := rootNode.FirstChild.LastChild
	paragraphID, ok := e.dollar["id"]
	if !ok {
		paragraphID = "e" + getRandomID(8)
	}

	for c := htmlBody.FirstChild; c != nil; c = c.NextSibling {

		// replace or <br> with <p></p>
		if c.Data == "br" {
			c.Data = "p"
		}

		c.Attr = append(c.Attr, html.Attribute{Key: "id", Val: paragraphID})
		c.Attr = append(c.Attr, html.Attribute{Key: "dd-level", Val: strconv.Itoa(e.level)})

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

	ownText := buf.String()
	childrenText := ""
	for _, child := range e.e {
		childrenText += child.getText()
	}

	return ownText + childrenText
}
