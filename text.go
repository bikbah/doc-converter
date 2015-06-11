package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Text struct {
	attribute    []Attribute
	text         string
	isInputField bool
}

func (text *Text) parse(i interface{}) {
	m, ok := i.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in Text.parse()")
	}

	text.isInputField = false
	text.attribute = []Attribute{}
	text.text = ""

	dollar, ok := m["$"]
	if ok {
		dollarMap, ok := dollar.(map[interface{}]interface{})
		if !ok {
			panic("Assert error in Text.parse()")
		}

		if tag, ok := dollarMap["tag"]; ok && fmt.Sprint(tag) == "span" {
			text.isInputField = true

			attrInterface, ok := m["attribute"]
			if !ok {
				panic("Input field text has no key ATTRIBUTE")
			}

			attrMap, ok := attrInterface.(map[interface{}]interface{})
			if !ok {
				panic("Assert error in Text.parse()")
			}

			for _, v := range attrMap {
				attr := Attribute{}
				attr.parse(v)
				text.addAttribute(attr)
			}
		}
	}

	if !text.isInputField {
		underscore, ok := m["_"]
		if !ok {
			panic("Text is not input field and does not have _ attribute")
		}

		text.text = fmt.Sprint(underscore)
	}
}

func (text *Text) addAttribute(attr Attribute) {
	text.attribute = append(text.attribute, attr)
}

func (text *Text) getText() string {
	if !text.isInputField {
		return text.text
	}

	rootNode, err := html.Parse(strings.NewReader("<span></span>"))
	if err != nil {
		panic("Parse span element error")
	}

	spanNode := rootNode.FirstChild.LastChild.FirstChild
	attrMapStr := make(map[string]string)
	for _, attr := range text.attribute {
		attrNameStr := strings.ToLower(attr.name)
		attrMapStr[attrNameStr] = attr.value
	}

	datafld, ok := attrMapStr["datafld"]
	if !ok {
		panic("Span has no datafld attribute..")
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
