# cli
A simple command line option parsing library

## Summary

A simple library that allows the options (long form only, i.e. --option) to be placed anywhere in the command line. String
and Int type options expect a value after the option, bool types canhave the value omitted to simply specify the flag.

Multiple options not define, the the value of the option is the last such specified option in the command line.


See the [example program](test-cli/main.go) for usage.
