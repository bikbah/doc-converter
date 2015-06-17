package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type E struct {
	dollar     map[string]string
	attribute  []Attribute
	e          []*E
	text       []Text
	level      int
	dependence []Dependence
	rawData    map[string]interface{}
	parent     *E
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
				childE.rawData = make(map[string]interface{})
				childE.parent = e

				childE.parse(childInterface)
				e.addChildE(&childE)
			}
			break

		case "dependence":
			depMap := make(map[string]interface{})

			for dK, dV := range vv {
				dKStr := fmt.Sprint(dK)
				depMap[dKStr] = dV
			}

			for i := 0; i < len(depMap); i++ {
				dKey := fmt.Sprint(i)
				childInterface, ok := depMap[dKey]
				if !ok {
					panic("Invalid index in e dependencies..")
				}

				dep := Dependence{}
				dep.parse(childInterface)
				e.addDependence(dep)
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

		default:
			e.rawData[strKey] = v
		}
	}

	e.setLevels()
}

func (e *E) setLevels() {
	childrenLevel := e.level + 1

	_, inVar := e.dollar["varcont"]
	if inVar {
		childrenLevel--
	}

	var linkCont E
	if eID, ok := e.dollar["id"]; ok && strings.Contains(eID, "linkContainer") {
		childrenLevel--
		linkCont = *e
	}

	for _, child := range e.e {

		if isVar, ok := child.dollar["var"]; (!ok || isVar != "1") && inVar {
			childrenLevel++
			inVar = false

			child.level = childrenLevel
			child.setLevels()
			continue
		}

		linkContLinkID, _ := linkCont.dollar["linkId"]
		childLinkID, _ := child.dollar["linkId"]
		if linkContLinkID != "" && childLinkID != linkContLinkID {
			childrenLevel++
			linkCont = E{}

			child.level = childrenLevel
			child.setLevels()
			continue
		}

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

func (e *E) addDependence(dep Dependence) {
	e.dependence = append(e.dependence, dep)
}

func (e *E) addChildE(child *E) {
	e.e = append(e.e, child)
}

func (e *E) getNumeration() string {
	ddNumeration := "0"
	if oldNumber, ok := e.dollar["number"]; ok {
		switch oldNumber {
		case "auto":
			ddNumeration = "-1"
		case "none":
			ddNumeration = "0"
		case "clear":
			ddNumeration = "0"
		case "parent":
			ddNumeration = "-1"
		default:
			log.Printf("Unknown numeration type: %s\n", oldNumber)
		}
	}

	if resetNumber, ok := e.dollar["resetnumber"]; ok && resetNumber == "true" {
		ddNumeration = "1"
	}

	return ddNumeration
}

func (e *E) getOwnText() string {
	ddNumeration := e.getNumeration()

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
		c.Attr = append(c.Attr, html.Attribute{Key: "dd-numeration", Val: ddNumeration})

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

var (
	linkCont E
	varText  string
)

func (e *E) getText() string {

	if id, ok := e.dollar["id"]; ok && strings.Contains(id, "linkContainer") {
		linkCont = *e

	}

	if varCont, ok := e.dollar["varcont"]; ok && varCont == "1" {
		varText = "v" + getRandomID(8)
	}

	_, okVarcont := e.dollar["varcont"]
	_, okVar := e.dollar["var"]

	if !okVarcont && !okVar && varText != "" {
		varText = ""
	}

	ownText := e.getOwnText()
	childrenText := ""
	for k, child := range e.e {

		// check if there is opened linkContainer group
		aCloseTag := ""
		contLink, _ := linkCont.dollar["linkId"]
		childLink, _ := child.dollar["linkId"]

		// close opened group if child is not logic child of linkCont or child is last
		if (contLink != "" && childLink != contLink) || (contLink != "" && k == len(e.e)-1) {
			aCloseTag = `<a/>`
			linkCont = E{}
		}

		childrenText += child.getText() + aCloseTag
	}

	allText := ownText + childrenText

	deps := make([]Dependence, 0)
	if stID, ok := e.dollar["storedId"]; ok {
		if depRaw, ok := linkCont.rawData["dependence_"+stID]; ok {

			depRawMap, ok := depRaw.(map[interface{}]interface{})
			if !ok {
				panic("Assert error in dependencies..")
			}

			depStrMap := make(map[string]interface{})
			for dK, dV := range depRawMap {
				dKstr := fmt.Sprint(dK)
				depStrMap[dKstr] = dV
			}

			for i := 0; i < len(depStrMap); i++ {
				childDep, ok := depStrMap[strconv.Itoa(i)]
				if !ok {
					panic("Invalid index in dependencies..")
				}

				dep := Dependence{}
				dep.parse(childDep)
				deps = append(deps, dep)
			}
		}
	}

	var resultDeps []Dependence
	if len(deps) > 0 {
		resultDeps = deps
	} else {
		resultDeps = e.dependence
	}

	// groups for linkContainer`s
	if eID, ok := e.dollar["id"]; ok && strings.Contains(eID, "linkContainer") {

		groupID := "g" + getRandomID(8)
		link, _ := e.dollar["link"]

		gr := Group{
			Id: groupID,
			// Dependency: getDependenciesText(resultDeps),
			Link: link,
		}

		addGroup(gr)
		allText = `<a id="` + groupID + `"/>` + allText

		return allText
	}

	// groups for var`s or e`s with dependencies
	if len(resultDeps) > 0 {

		groupID := "g" + getRandomID(8)
		varName, _ := e.dollar["variantName"]

		gr := Group{
			Id:         groupID,
			Name:       varText,
			Text:       varName,
			Dependency: getDependenciesText(resultDeps, *e),
		}

		addGroup(gr)
		allText = `<a id="` + groupID + `"/>` + allText + `<a/>`
	}

	return allText
}
