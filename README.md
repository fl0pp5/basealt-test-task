# repocmp #
Utility to find difference between two branches

## Building ##
```shell
$ go build cmd/repocmp.go
```

## Usage ##
```shell
$ ./repocmp -h
Usage: ./repocmp <option> [args]
Example: ./repocmp -b1 <branch-name> -fb2 <filename> -cache <file-prefix> -pretty
  -b1 value
        first branch name
  -b2 value
        second branch name
  -cache string
        save downloaded branches. Usage: [...] -cache <file-prefix> [...]
  -fb1 value
        first branch filename
  -fb2 value
        second branch filename
  -pretty
        enable formatting
  -split string
        split output by files. Usage: [...] -split <file-prefix> [...]

```

## Examples ##

#### Print pretty to stdout ####
```shell
$ ./repocmp -b1 p10 -b2 p9 -pretty
{
    "unique_packages_1": {...}
    "unique_packages_2": {...}
    "newer_packages_from_1": {...}
}
```
#### With origin save ####
```shell
$ ./repocmp -b1 p10 -b2 p9 -cache myprefix > diff.json
$ ls
... diff.json  myprefix_branch1_963172800.json myprefix_branch2_963172800.json ...
```
#### With json file & save ####
```shell
$ ./repocmp -fb1 myprefix_branch1_963172800.json -b2 p9 -cache p9
{
    "unique_packages_1": {...}
    "unique_packages_2": {...}
    "newer_packages_from_1": {...}
}
$ ls
... p9_branch2_963173340 ...
```

#### Split output ####
```shell
$ ./repocmp -b1 p10 -b2 sisyphus -split p10_sisyphus
$ ls
... p10_sisyphus_newer_1668096715.json  p10_sisyphus_unique1_963173748.json  p10_sisyphus_unique2_963173748.json ...
```

## Output schema ##
```go
type BranchDiff struct {
	UniquePackages1    Branch `json:"unique_packages_1"`
	UniquePackages2    Branch `json:"unique_packages_2"`
	NewerPackagesFrom1 Branch `json:"newer_packages_from_1"`
}
```
