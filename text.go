package main

import "fmt"

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
