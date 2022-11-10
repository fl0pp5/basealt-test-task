package main

import (
	"flag"
	"fmt"
	"os"
	"repocmp/pkg/api"
	"repocmp/pkg/branch"
	"repocmp/pkg/helpers"
	"time"
)

type BranchInfoFlag struct {
	name   string
	isFile bool
	get    func() (*api.Branch, error)
}

func (f *BranchInfoFlag) fromFile(value string) error {
	f.name = value
	f.isFile = true
	f.get = func() (*api.Branch, error) {
		return branch.FromFile(f.name)
	}
	return helpers.ValidateFilename(f.name)
}

func (f *BranchInfoFlag) fromHttp(value string) error {
	f.name = value
	f.get = func() (*api.Branch, error) {
		return branch.FromHttp(f.name)
	}
	return helpers.ValidateBranchName(f.name)
}

var (
	branch1     BranchInfoFlag
	branch2     BranchInfoFlag
	savePrefix  string
	splitPrefix string
	pretty      bool
)

func getBranches() (*api.Branch, *api.Branch) {
	a, err := branch1.get()
	helpers.FatalIf(err)
	b, err := branch2.get()
	helpers.FatalIf(err)

	return a, b
}

func makeFilename(prefix string, packageName string) string {
	return fmt.Sprintf("%s_%s_%d.json", prefix, packageName, time.Now().Unix())
}

func printBranchDiff(diff *api.BranchDiff) {
	if !helpers.HasFlag("split") {
		data, err := helpers.GetMarshaller(diff, pretty)()
		helpers.FatalIf(err)
		fmt.Println(string(data))
		return
	}

	helpers.FatalIf(helpers.WriteJsonToFile(diff.UniquePackages1,
		makeFilename(splitPrefix, "unique1"), pretty))
	helpers.FatalIf(helpers.WriteJsonToFile(diff.UniquePackages2,
		makeFilename(splitPrefix, "unique2"), pretty))
	helpers.FatalIf(helpers.WriteJsonToFile(diff.NewerPackagesFrom1,
		makeFilename(splitPrefix, "newer"), pretty))
}

func saveBranches(a, b *api.Branch) {
	if !helpers.HasFlag("cache") {
		return
	}

	if !branch1.isFile {
		helpers.FatalIf(helpers.WriteJsonToFile(a,
			makeFilename(savePrefix, "branch1"), pretty))
	}

	if !branch2.isFile {
		helpers.FatalIf(helpers.WriteJsonToFile(b,
			makeFilename(savePrefix, "branch2"), pretty))
	}
}

func checkRequiredFlags() {
	count := 0

	if !helpers.HasAnyRequiredFlags([]string{"b1", "fb1"}) {
		count++
		fmt.Println("option: `[f]b1` is required")
	}

	if !helpers.HasAnyRequiredFlags([]string{"b2", "fb2"}) {
		count++
		fmt.Println("option: `[f]b2` is required")
	}

	if count > 0 {
		os.Exit(0)
	}
}

func initFlags() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s <option> [args]\n"+
				"Example: %s -b1 <branch-name> -fb2 <filename> -cache <file-prefix> -pretty\n",
			os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Func("b1", "first branch name", branch1.fromHttp)
	flag.Func("b2", "second branch name", branch2.fromHttp)
	flag.Func("fb1", "first branch filename", branch1.fromFile)
	flag.Func("fb2", "second branch filename", branch2.fromFile)
	flag.StringVar(&savePrefix, "cache", "",
		"save downloaded branches. Usage: [...] -cache <file-prefix> [...]")
	flag.StringVar(&splitPrefix, "split", "",
		"split output by files. Usage: [...] -split <file-prefix> [...]")
	flag.BoolVar(&pretty, "pretty", false, "enable formatting")
	flag.Parse()

}

func main() {
	initFlags()
	checkRequiredFlags()
	a, b := getBranches()
	saveBranches(a, b)
	diff := api.BranchDiff{
		UniquePackages1:    *branch.Diff(a, b),
		UniquePackages2:    *branch.Diff(b, a),
		NewerPackagesFrom1: *branch.Newer(a, b),
	}

	printBranchDiff(&diff)
}
