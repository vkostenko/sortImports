## SortImports

Joins imports in file into single group before running goimports command

## Requirements

[goimports](https://godoc.org/golang.org/x/tools/cmd/goimports): `go get golang.org/x/tools/cmd/goimports`

## Usage in Goland IDE

`go get github.com/vkostenko/sortImports`

In Goland: Preferences -> FileWatchers -> Add (+) -> goimports -> Fill:
* Program: `sortImports`, 
* Arguments: `-local somelocal,packages -w -srcdir $FilePath$`

, other options left as is -> Ok

## Fix all files in entire catalog recursively

`find . | grep "\.go" | cut -c 3- | xargs -L1 -t sortImports -local somelocal,packages -w -srcdir`
