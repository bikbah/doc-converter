package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/bikbah/go-php-serialize/phpserialize"
)

const (
	MaxDepth int = 3
)

var (
	keysCount map[string]int
)

func main() {
	// Read php serialized file
	dat, err := ioutil.ReadFile("1")
	if err != nil {
		log.Fatalf("Read file error: %v", err)
	}

	// Try to deserialize data
	res, err := phpserialize.Decode(string(dat))
	if err != nil {
		log.Fatalf("Deserialize data error: %v", err)
	}

	m, ok := res.(map[interface{}]interface{})
	if !ok {
		log.Fatal("interface{} to map[interface{}]interface{} assertion fail..")
	}
	convertedMap := parseRawMap(m)
	// convertedMap["e"] = convertedMap["e"].(map[string]interface{})
	// tmp := convertedMap["e"].(map[string]interface{})
	fmt.Println(convertedMap["e"])
	// fmt.Println(tmp["0"])
	// fmt.Println(reflect.TypeOf(convertedMap))

	// See key types
	// keysCount = make(map[string]int)
	// tmpFunc(convertedMap)
	// fmt.Println(keysCount)

}

func parseRawMap(m map[interface{}]interface{}) map[string]interface{} {
	// if depth > MaxDepth {
	// 	return
	// }
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

func tmpFunc(m map[string]interface{}) {
	for _, v := range m {
		// log.Println(reflect.TypeOf(v).String())
		strKey := reflect.TypeOf(v).String()
		_, ok := keysCount[strKey]
		if ok {
			keysCount[strKey]++
		} else {
			keysCount[strKey] = 0
		}

		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			vv, ok := v.(map[string]interface{})
			if !ok {
				log.Fatal("interface{} to map[interface{}]interface{} assertion fail..")
			}
			tmpFunc(vv)
		}
	}
}
