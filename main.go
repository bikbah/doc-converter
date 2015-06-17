// TODO: make valid groups
// TODO: implement dependencies
// TODO: implement variants
// TODO: exclude spreadSheetTable shit from tables

package main

import (
	"io/ioutil"
	"log"

	"github.com/bikbah/go-php-serialize/phpserialize"
)

func main() {
	// Read php serialized file
	dat, err := ioutil.ReadFile("2")
	if err != nil {
		log.Fatalf("Read file error: %v", err)
	}

	// Try to deserialize data
	res, err := phpserialize.Decode(string(dat))
	if err != nil {
		log.Fatalf("Deserialize data error: %v", err)
	}

	rootE := E{}
	rootE.rawData = make(map[string]interface{})
	rootE.parse(res)

	saveData(rootE.getText())
	saveGroups()
	saveErrors()
}
