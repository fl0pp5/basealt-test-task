package main

import (
	"flag"
	"fmt"
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
	return nil
}

func (f *BranchInfoFlag) fromHttp(value string) error {
	f.name = value
	f.get = func() (*api.Branch, error) {
		return branch.FromHttp(f.name)
	}
	return nil
}

var (
	branch1     BranchInfoFlag
	branch2     BranchInfoFlag
	savePrefix  string
	splitPrefix string
	pretty      bool
)

var usage = map[string]string{
	"b1":     "first branch name",
	"b2":     "second branch name",
	"fb1":    "first branch filename",
	"fb2":    "second branch filename",
	"cache":  "save downloaded branches",
	"split":  "split output by files",
	"pretty": "enable formatting",
}

func getBranches() (*api.Branch, *api.Branch) {
	a, err := branch1.get()
	helpers.FatalIf(err)
	b, err := branch2.get()
	helpers.FatalIf(err)

	return a, b
}

func makeFilename(prefix string, packageName string) string {
	return fmt.Sprintf("%s_%s_%d\n", prefix, packageName, time.Now().Unix())
}

func printBranchDiff(diff *api.BranchDiff) {
	if splitPrefix == "" {
		data, err := helpers.GetMarshaller(diff, pretty)()
		helpers.FatalIf(err)
		fmt.Println(string(data))
		return
	}
	fmt.Println(pretty)
	helpers.FatalIf(helpers.WriteJsonToFile(diff.UniquePackages1,
		makeFilename(splitPrefix, "unique1"), pretty))
	helpers.FatalIf(helpers.WriteJsonToFile(diff.UniquePackages2,
		makeFilename(splitPrefix, "unique2"), pretty))
	helpers.FatalIf(helpers.WriteJsonToFile(diff.NewerPackagesFrom1,
		makeFilename(splitPrefix, "newer"), pretty))
}

func saveBranch(a, b *api.Branch) {
	if savePrefix == "" {
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

func initFlags() {
	flag.Func("b1", usage["b1"], branch1.fromHttp)
	flag.Func("b2", usage["b2"], branch2.fromHttp)
	flag.Func("fb1", usage["fb1"], branch1.fromFile)
	flag.Func("fb2", usage["fb1"], branch2.fromFile)
	flag.StringVar(&savePrefix, "cache", "", usage["cache"])
	flag.StringVar(&splitPrefix, "split", "", usage["split"])
	flag.BoolVar(&pretty, "pretty", false, usage["pretty"])
	flag.Parse()

}

func main() {
	initFlags()
	a, b := getBranches()
	saveBranch(a, b)
	adiff := branch.Diff(a, b)
	bdiff := branch.Diff(b, a)
	newer := branch.Newer(a, b)

	diff := api.BranchDiff{
		UniquePackages1:    *adiff,
		UniquePackages2:    *bdiff,
		NewerPackagesFrom1: *newer,
	}

	printBranchDiff(&diff)
}
