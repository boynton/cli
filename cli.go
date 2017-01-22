// A simple command line argument parsing lib
package cli

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Params []string
type Options map[string]interface{}

type Option struct {
	Name         string
	Def          interface{}
	Descr        string
	Type         string
	intTarget    *int
	stringTarget *string
	boolTarget   *bool
}

type Command struct {
	Name        string
	Description string
	Options     map[string]*Option
}

func New(name string, descr string) *Command {
	cmd := &Command{
		Name:        name,
		Description: descr,
		Options:     make(map[string]*Option),
	}
	return cmd
}

func (cmd *Command) sortedOptions() []*Option {
	var names []string
	for k, _ := range cmd.Options {
		names = append(names, k)
	}
	sort.Strings(names)
	var lst []*Option
	for _, name := range names {
		lst = append(lst, cmd.Options[name])
	}
	return lst
}

func (cmd *Command) Usage() string {
	//i.e. "usage: test-cli [--option1 <string>] [--option2 <string>] [--option3 <int>] [--option4 [<bool>]] param ..."
	s := "usage: " + cmd.Name
	for _, opt := range cmd.sortedOptions() {
		s += fmt.Sprintf(" [--%s <%s>]", opt.Name, opt.Type)
	}
	s += " [param ...]\n"
	if cmd.Description != "" {
		s += "\n" + cmd.Description
	}
	return s
}

func (cmd *Command) StringOption(target *string, name string, def string, descr string) {
	cmd.Options[name] = &Option{Name: name, Def: def, Descr: descr, Type: "string", stringTarget: target}
}

func (cmd *Command) IntOption(target *int, name string, def int, descr string) {
	cmd.Options[name] = &Option{Name: name, Def: def, Descr: descr, Type: "int", intTarget: target}
}

func (cmd *Command) BoolOption(target *bool, name string, def bool, descr string) {
	cmd.Options[name] = &Option{Name: name, Def: def, Descr: descr, Type: "bool", boolTarget: target}
}

//Parse - parse the command line, return a list of params and a map of options
func (cmd *Command) Parse() (Params, Options) {
	args := os.Args
	options := make(map[string]interface{})
	for _, opt := range cmd.Options {
		options[opt.Name] = opt.Def
	}
	max := len(args)
	var expectedOpt *Option
	var params []string
	for i := 1; i < max; i++ {
		var param string
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			name := arg[2:]
			if opt, ok := cmd.Options[name]; ok {
				if expectedOpt != nil {
					if expectedOpt.Type == "bool" {
						options[expectedOpt.Name] = true
					} else {
						Fatal("Missing value for " + arg)
					}
				}
				expectedOpt = opt
			} else {
				Fatal("Unknown option: " + arg)
			}
		} else if expectedOpt != nil {
			var val interface{}
			switch expectedOpt.Type {
			case "string":
				val = arg
			case "int":
				n, err := strconv.Atoi(arg)
				if err != nil {
					Fatal("Bad int: " + arg)
				}
				val = n
				options[expectedOpt.Name] = val
			case "bool":
				if strings.ToLower(arg) == "true" {
					val = true
				} else if strings.ToLower(arg) == "false" {
					val = false
				} else {
					val = true
					param = arg
				}
			}
			options[expectedOpt.Name] = val
			expectedOpt = nil
		} else {
			param = arg
		}
		if param != "" {
			params = append(params, param)
		}
	}
	if expectedOpt != nil {
		if expectedOpt.Type == "bool" {
			options[expectedOpt.Name] = true
		} else {
			Fatal("Missing value for " + expectedOpt.Name)
		}
	}
	for name, val := range options {
		opt := cmd.Options[name]
		switch opt.Type {
		case "string":
			s, _ := val.(string)
			if opt.stringTarget != nil {
				*opt.stringTarget = s
			}
		case "int":
			i, _ := val.(int)
			if opt.intTarget != nil {
				*opt.intTarget = i
			}
		case "bool":
			b, _ := val.(bool)
			if opt.boolTarget != nil {
				*opt.boolTarget = b
			}
		default:
		}
	}
	return params, options
}

func Fatal(msg interface{}) {
	fmt.Printf("*** %v\n", msg)
	os.Exit(1)
}

//--- the following is a different dynamic approach. No spec, just parse what you find

func Parse() (Params, Options) {
	args := os.Args[1:]
	options := make(map[string]interface{})
	params := make([]string, 0)
	var expectedOption []string
	for _, token := range args {
		if expectedOption == nil {
			if strings.HasPrefix(token, "--") {
				name := token[2:]
				expectedOption = strings.Split(name, ".")
			} else {
				params = append(params, token)
			}
		} else {
			if strings.HasPrefix(token, "--") {
				Fatal("Missing value for " + strings.Join(expectedOption, "."))
			}
			put(options, expectedOption, token)
			expectedOption = nil
		}
	}
	return params, options
}

func put(obj map[string]interface{}, path []string, tokenValue interface{}) {
	if len(path) == 1 {
		obj[path[0]] = tokenValue
	} else {
		var next map[string]interface{}
		tmp, ok := obj[path[0]]
		if ok {
			next, ok = tmp.(map[string]interface{})
			if !ok {
				Fatal("Malformed object reference: " + fmt.Sprint(path))
			}
		} else {
			next = make(map[string]interface{})
			obj[path[0]] = next
		}
		put(next, path[1:], tokenValue)
	}
}

func (opt Options) Get(path string) interface{} {
	o, err := get(opt, strings.Split(path, "."))
	if err != nil {
		Fatal(err.Error())
	}
	return o
}

func (opt Options) GetObject(path string) map[string]interface{} {
	o := opt.Get(path)
	if o != nil {
		if m, ok := o.(map[string]interface{}); ok {
			return m
		}
		Fatal("Not an object: " + fmt.Sprint(o))
	}
	return nil
}

func (opt Options) GetInt(path string, defaultValue int) int {
	o := opt.Get(path)
	if o != nil {
		if n, ok := o.(int); ok {
			return n
		}
		Fatal("Not an int: " + fmt.Sprint(o))
	}
	return defaultValue
}

func (opt Options) GetString(path string, defaultValue string) string {
	o := opt.Get(path)
	if o != nil {
		if s, ok := o.(string); ok {
			return s
		}
		Fatal("Not a string: " + fmt.Sprint(o))
	}
	return defaultValue
}

func (opt Options) GetBool(path string, defaultValue bool) bool {
	o := opt.Get(path)
	if o != nil {
		if s, ok := o.(bool); ok {
			return s
		}
		Fatal("Not a bool: " + fmt.Sprint(o))
	}
	return defaultValue
}

func get(obj map[string]interface{}, path []string) (interface{}, error) {
	if len(path) == 1 {
		o, ok := obj[path[0]]
		if !ok {
			return nil, nil
		}
		return o, nil
	} else {
		var next map[string]interface{}
		tmp, ok := obj[path[0]]
		if ok {
			next, ok = tmp.(map[string]interface{})
			if ok {
				return get(next, path[1:])
			}
			return nil, fmt.Errorf("structure incongruence for %v", path)
		}
		return nil, nil
	}
}
