package main

import (
	"fmt"
	"go-php-serialize/phpserialize"
	"io/ioutil"
	"log"
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
	m := res.(map[string]interface{})
	fmt.Println(m)
	// fmt.Printf("%+v", res)
	// fmt.Println(m)
	// fmt.Println(reflect.TypeOf(res))
	// fmt.Println(res['$'])
	// m := res.(map[string]interface{})
	// fmt.Println(m)

}
