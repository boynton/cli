package main

import (
	"encoding/json"
	"fmt"
	"github.com/boynton/cli"
)

// i.e. test-cli --entity.age 23 --entity.name Joe --foo bar --blah true --glorp 100
func main() {
	params, options := cli.Parse()
	fmt.Println("Params:", Pretty(params))
	fmt.Println("Options:", Pretty(options))

	age := options.GetInt("entity.age", 100)
	fmt.Println("age:", age)
	name := options.GetString("entity.name", "anonymous")
	fmt.Println("name:", name)
	blah := options.GetBool("blah", false)
	fmt.Println("blah:", blah)
	foo := options.GetString("foo", "")
	fmt.Println("foo:", foo)
	fmt.Println("entity:", Pretty(options.GetObject("entity")))
}

func Pretty(obj interface{}) string {
	b, _ := json.MarshalIndent(obj, "", "    ")
	return string(b)
}
