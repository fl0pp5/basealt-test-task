package helpers

import (
	"fmt"
	"repocmp/pkg/api"
)

func ValidateBranchName(name string) error {
	if _, ok := api.AllowedBranches[name]; !ok {
		return fmt.Errorf("invalid branch name: `%s` is not allowed\n", name)
	}
	return nil
}
