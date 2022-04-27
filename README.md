# lineremover
A simple tool for removing (repetitive and irrelevant) lines of text from a text file!

The initial intent of this was to simply process logs to get rid of some of the noise when trying to debug customer issues. I've used this plenty of times (as have my peers), so I might as well make it public!

By default, it will generate a new file with any lines containing `keys` stripped, and leave the original untouched.

# How to run
It's probably best to build a binary and run that, but it can be run as follows:

`go run main.go -file=<path/to/file> -keys="<keys to search for|with multiple separated by|pipes>"`

Optionally, you can provide a `-inplace=true` flag to have the processed file overwrite the original. This can be useful if you're in discovery mode and actively figuring out which lines are irrelevant to whatever it is you're trying to do.

# How to build

`go build`

and then run as:

`./lineremover -file=<path/to/file> -keys="<keys to search for|with multiple separated by|pipes>"`

