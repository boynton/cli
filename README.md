# cli
A simple command line option parsing library for Go

## Summary

A simple library that allows the options (long form only, i.e. --option) to be placed anywhere in the command line. String
and Int type options expect a value after the option, bool types can have the value omitted to simply specify the flag.

If an option name contains a ".", then a structured object is formed (of type map[string]string), i.e. `--foo.one blah --foo.two.three bletch` will produce an option
`foo` with the value `{"two": {"three": "bletch"},"one": "blah"}`.

Multiple identical options are not allowed, the the value of the option is the last such specified option in the command line.

The API supports incremental specification: You can use it without specifying any options, resulting in a dynamic API, or you
can declare options, in which case a usage line can be generated, and typed variables for each option are supported.

See [this example program](test-cli/main.go) for the first dynamic example, and [this one](test2-cli/main.go) for an example
of declaring things. In the latter example, `nil` may be passed as the variable address to avoid use of variables (i.e. to be
compatible with the dynamic usage pattern).

In any case, the returned Params and Options types support typed fetches with defaults, i.e. `options.GetInt("argname", 23)`, so
a program can start with no spec, and specs can be added without having to otherwise change the code.

