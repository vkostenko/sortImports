## SortImports

Joins imports in file into single group before running goimports command

## Requirements

[goimports](https://godoc.org/golang.org/x/tools/cmd/goimports): `go get golang.org/x/tools/cmd/goimports`

## Usage in Goland IDE

`go get github.com/vkostenko/sortImports`

In Goland: Preferences -> FileWatchers -> Add (+) -> goimports -> Fill:
* Program: `sortImports`, 
* Arguments: `-local _local,packages_ -w -srcdir $FilePath$`

, other options left as is -> Ok
