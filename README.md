# EIC: Ensure // import "comment" using `go/ast` 
 
## Usage

```
$go get github.com/anarcher/eic
$eic
Ensure import comment

Usage:
  eic [flags]

Flags:
  -d, --dir string    transfer directory
  -n, --dryrun        show what would have been transferred
  -f, --file string   transfer a file
  -h, --help          help for eic
(go:anarcher) anarch@vacuum-2 ~/go/anarch
```

## Example

```
$cat ./main.go | head -n 1
package main

$eic -n -f ./main.go  | head -n 1
package main // import "github.com/anarcher/eic"
```

