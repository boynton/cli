package main

import (
	"encoding/json"
	"fmt"
	"github.com/boynton/cli"
)

func main() {
	var option1 string
	var option2 string
	var option3 int
	var option4 bool
	var help bool

	cmd := cli.New("test-cli", "A test CLI program. Try out various options and params to see the effect.")
	cmd.StringOption(&option1, "option1", "default string", "This is the first option")
	cmd.StringOption(&option2, "option2", "another default", "This is the second option")
	cmd.IntOption(&option3, "option3", 23, "An int option")
	cmd.BoolOption(&option4, "option4", false, "A bool option")
	cmd.StringOption(nil, "option5", "bletch", "An option with no static variable defined")
	cmd.BoolOption(&help, "help", false, "Show help")

	params, options := cmd.Parse()

	if help {
		fmt.Println(cmd.Usage())
	} else {
		fmt.Println("non-option params:", Pretty(params)) //the non-option params in a single Array
		fmt.Println("options:", Pretty(options)) //options in a single JSON object, note default values

		//dynamic lookup, like the other example uses
		fmt.Println("option3:", options.GetInt("option3", 0))

		//type-safe version
		fmt.Println("option1:", option1)
		fmt.Println("option2:", option2)
		fmt.Println("option3:", option3)
		fmt.Println("option4:", option4)

	}
}

func Pretty(obj interface{}) string {
	b, _ := json.MarshalIndent(obj, "", "    ")
	return string(b)
}
