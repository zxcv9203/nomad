package main

import (
	"fmt"
	"log"

	"github.com/zxcv9203/nomad/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "first word"}
	find, err := dictionary.Search("first")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(find)
	}
}
