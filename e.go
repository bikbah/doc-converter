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

func (e *E) addDependence(dep Dependence) {
	e.dependence = append(e.dependence, dep)
}

func (e *E) addChildE(child *E) {
	e.e = append(e.e, child)
}

var linkCont E

func (e *E) getText() string {

	if id, ok := e.dollar["id"]; ok && strings.Contains(id, "linkContainer") {
		linkCont = *e
	}

	ddNumeration := 0
	if oldNumber, ok := e.dollar["number"]; ok {
		switch oldNumber {
		case "auto":
			ddNumeration = -1
		case "none":
			ddNumeration = 0
		case "clear":
			ddNumeration = 0
		case "parent":
			ddNumeration = -1
		default:
			log.Printf("Unknown numeration type: %s\n", oldNumber)
		}
	}

	if resetNumber, ok := e.dollar["resetnumber"]; ok && resetNumber == "true" {
		ddNumeration = 1
	}

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
		c.Attr = append(c.Attr, html.Attribute{Key: "dd-numeration", Val: strconv.Itoa(ddNumeration)})

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

	if len(resultDeps) > 0 {
		if strings.Contains(e.dollar["id"], "linkContainer") {
			log.Println(e.dollar["id"])
		}
		groupID := "g" + getRandomID(8)

		eID, ok := e.dollar["id"]
		if !ok {
			panic("E has no id..")
		}

		linkID, _ := e.dollar["linkId"]

		gr := Group{
			Id:         groupID,
			Dependency: getDependenciesText(resultDeps),
			Text:       eID,
			Link:       linkID,
		}
		addGroup(gr)
		allText = `<a id="` + groupID + `/">` + allText + `<a/>`
	}

	return allText
}
