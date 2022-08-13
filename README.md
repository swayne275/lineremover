# lineremover
A simple tool for removing (repetitive and irrelevant) lines of text from a text file!

The initial intent of this was to simply process logs to get rid of some of the noise when trying to debug customer issues. I've used this plenty of times (as have my peers), so I might as well make it public!

By default, it will generate a new file with any lines containing `keys` or matching `pattern` stripped, and leave the original untouched.

`keys` or `pattern` can be supplied individually, or in combination.

Note that due to limitations in Go's standard library `regexp`, some lookarounds are not supported.

# How to run
It's probably best to build a binary and run that, but it can be run as follows:

`go run main.go -file=<path/to/file> -keys="<keys to search for|with multiple separated by|pipes>" -pattern="<regular expression to match>"`

Optionally, you can provide a `-inplace=true` flag to have the processed file overwrite the original. This can be useful if you're in discovery mode and actively figuring out which lines are irrelevant to whatever it is you're trying to do.

# How to build

`go build`

and then run as:

`./lineremover -file=<path/to/file> -keys="<keys to search for|with multiple separated by|pipes>" -pattern="<regular expression to match>"`

# Examples

The following examples work on `input.txt` in the `examples/` folder.

The following code snipets expect to be executed from the directory containing `main.go`.

Note that `go run main.go` can be replaced with `./lineremover` if you've built a binary
(called `lineremover`).

They can also be modified by appending `-inplace` if you'd like to overwrite `examples/input.txt`.

## Removing everything prefixed with "hello"

`go run main.go -file="example/input.txt" -keys="hello"`

## Removing everything prefixed with "hello small"

`go run main.go -file="example/input.txt" -keys="hello small"`

## Removing everything containing "world"

`go run main.go -file="example/input.txt" -keys="world"`

## Removing everything containing "big", "brig", "bight", or "bright"

`go run main.go -file="example/input.txt" -pattern=".*b([r]?)ig([ht]?).*"`
