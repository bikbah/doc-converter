package main

import "fmt"

type Dependence struct {
	option    string
	condition string
	value     string
	operation string
}

var invalidDepsErrors []string

func init() {
	invalidDepsErrors = make([]string, 0)
}

func (dep *Dependence) parse(i interface{}) {
	m, ok := i.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in Dependence.parse()")
	}

	dollar, ok := m["$"]
	if !ok {
		panic("Dependence has no $ field..")
	}

	dollarMap, ok := dollar.(map[interface{}]interface{})
	if !ok {
		panic("Assert error in Dependence.parse()")
	}

	depOption, ok := dollarMap["option"]
	if !ok {
		panic("Dependence has no OPTION field..")
	}

	depCondition, ok := dollarMap["condition"]
	if !ok {
		panic("Dependence has no CONDITION field..")
	}

	depValue, ok := dollarMap["value"]
	if !ok {
		panic("Dependence has no VALUE field..")
	}

	if depOperation, ok := dollarMap["operation"]; ok {
		dep.operation = fmt.Sprint(depOperation)
	}

	dep.option = fmt.Sprint(depOption)
	dep.condition = fmt.Sprint(depCondition)
	dep.value = fmt.Sprint(depValue)
}

func getDependenciesText(deps []Dependence, e E) string {
	res := ""
	for k, dep := range deps {

		res += dep.operation

		res += dep.option

		cond := dep.condition
		if cond == "=" {
			cond = "=="
		}

		res += cond
		res += dep.value

		if k > 0 && dep.operation == "" {
			if id, ok := e.dollar["id"]; ok && id != "" {
				invalidDepsErrors = append(invalidDepsErrors, id)
			}
		}
	}

	return res
}
