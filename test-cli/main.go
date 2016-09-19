package main

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/boynton/cli"
)

// i.e. test-cli --entity.age 23 --entity.name Joe --foo bar --blah true --glorp 100
func main() {
	ctx := cli.Parse(os.Args)
	fmt.Println(Pretty(ctx))

	age := ctx.GetInt("entity.age", 100)
	fmt.Println("age:", age)
	name := ctx.GetString("entity.name", "anonymous")
	fmt.Println("name:", name)
	blah := ctx.GetBool("blah", false)
	fmt.Println("blah:", blah)
	foo := ctx.GetString("foo", "")
	fmt.Println("foo:", foo)
	fmt.Println("entity:", Pretty(ctx.GetObject("entity")))
}

func Pretty(obj interface{}) string {
	b, _ := json.MarshalIndent(obj, "", "    ")
	return string(b)
}
