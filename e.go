package main

import "fmt"

type E struct {
	dollap    map[string]string
	attribute []Attribute
	e         []*E
	text      []Text
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
			e.dollap = make(map[string]string)

			for dollarK, dollarV := range vv {
				dollarKStr := fmt.Sprint(dollarK)
				dollarVStr := fmt.Sprint(dollarV)

				e.dollap[dollarKStr] = dollarVStr
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
		}
	}
}

func (e *E) addAttribute(attr Attribute) {
	e.attribute = append(e.attribute, attr)
}

func (e *E) addChildE(child *E) {
	e.e = append(e.e, child)
}
