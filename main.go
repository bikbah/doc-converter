package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bikbah/go-php-serialize/phpserialize"
)

func main() {
	dat, err := ioutil.ReadFile("1")
	if err != nil {
		log.Fatalf("Read file error: %v", err)
	}

	// fmt.Println(string(dat))
	var res interface{}
	res, err = phpserialize.Decode(string(dat))
	if err != nil {
		log.Fatalf("Deserialize data error: %v", err)
	}
	m, _ := res.(map[interface{}]interface{})
	// fmt.Println(m, casted)
	for k, _ := range m {
		fmt.Println(m[k])
		return
	}
	for k, _ := range res {
		fmt.Println(res[k])
		return
	}
	// fmt.Printf("%+v", res)
	// fmt.Println(m)
	// fmt.Println(reflect.TypeOf(res))
	// fmt.Println(res['$'])
	// m := res.(map[string]interface{})
	// fmt.Println(m)

}
