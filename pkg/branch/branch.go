package branch

import (
	"encoding/json"
	version "github.com/knqyf263/go-rpm-version"
	"io"
	"net/http"
	"os"
	"repocmp/pkg/api"
	"repocmp/pkg/helpers"
)

func FromFile(filename string) (*api.Branch, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var branch api.Branch
	if err = json.Unmarshal(data, &branch); err != nil {
		return nil, err
	}

	return &branch, nil
}

func FromHttp(branchName string) (*api.Branch, error) {
	if err := helpers.ValidateBranchName(branchName); err != nil {
		return nil, err
	}

	url := api.BaseUrl + api.ExportBranchBinaryPackages + "/" + branchName

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var branch api.Branch
	if err := json.Unmarshal(data, &branch); err != nil {
		return nil, err
	}

	return &branch, nil
}

func Diff(a, b *api.Branch) *api.Branch {
	tmp := make(map[string]struct{}, b.Length)
	for _, item := range b.Packages {
		tmp[item.Name+item.Arch] = struct{}{}
	}

	var diff api.Branch
	for _, item := range a.Packages {
		if _, ok := tmp[item.Name+item.Arch]; !ok {
			diff.Packages = append(diff.Packages, item)
		}
	}

	diff.Length = len(diff.Packages)

	return &diff
}

func Newer(a, b *api.Branch) *api.Branch {
	tmp := make(map[string]*api.Package, a.Length)
	for i := range a.Packages {
		itemA := a.Packages[i]
		tmp[itemA.Name+itemA.Arch] = &itemA
	}

	var diff api.Branch
	for _, itemB := range b.Packages {
		if itemA, ok := tmp[itemB.Name+itemB.Arch]; ok {
			versionA := version.NewVersion(itemA.Version + "-" + itemA.Release)
			versionB := version.NewVersion(itemB.Version + "-" + itemB.Release)

			if versionA.GreaterThan(versionB) {
				diff.Packages = append(diff.Packages, *itemA)
			}
		}
	}

	diff.Length = len(diff.Packages)

	return &diff
}
