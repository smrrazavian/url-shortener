package main

import (
	"fmt"

	"github.com/smrrazavian/url-shortener/pkg/idgen"
)

func main() {
	id, _ := idgen.GenerateID()
	fmt.Println(id)
}
