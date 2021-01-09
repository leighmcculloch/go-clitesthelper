# go-testcli
Go package containing helpers for testing CLIs.

## Usage

Define your main function so that it matches the signature of `testcli.MainFunc`. In tests, call your main function using `testcli.Main` which will capture stdout, stderr, and the exit code. Use the other helper functions to create a working directory that suites the use cases that need testing.

For an example of this used in the wild, see the tests of [gas](https://github.com/leighmcculloch/gas).
