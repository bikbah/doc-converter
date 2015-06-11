package main

type Text struct {
	attribute    []Attribute
	text         string
	isInputField bool
func (text *Text) addAttribute(attr Attribute) {
	text.attribute = append(text.attribute, attr)
}
