package main

import "fmt"

type Attribute struct {
	name  string
	value string
}

func (attr *Attribute) parse(i interface{}) {
	m, ok := i.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in Attribute.parse()")
	}

	dollar, ok := m["$"]
	if !ok {
		panic("Attribute has no $ field..")
	}

	dollarMap, ok := dollar.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in Attribute.parse()")
	}

	attrName, ok := dollarMap["name"]
	if !ok {
		panic("Attribute has no NAME field..")
	}

	attrValue, ok := dollarMap["value"]
	if !ok {
		panic("Attribute has no VALUE field..")
	}

	attr.name = fmt.Sprint(attrName)
	attr.value = fmt.Sprint(attrValue)
}
