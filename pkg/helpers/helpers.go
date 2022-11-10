package helpers

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"repocmp/pkg/api"
)

func ValidateBranchName(name string) error {
	if _, ok := api.AllowedBranches[name]; !ok {
		return fmt.Errorf("invalid branch name: `%s` is not allowed\n", name)
	}
	return nil
}

func GetMarshaller(v any, indent bool) func() ([]byte, error) {
	return func() ([]byte, error) {
		if indent {
			return json.MarshalIndent(v, "", "    ")
		}
		return json.Marshal(v)
	}
}

func WriteJsonToFile(v any, filename string, indent bool) error {
	data, err := GetMarshaller(v, indent)()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

func FatalIf(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func HasFlag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
func HasAnyRequiredFlags(required []string) bool {
	for _, item := range required {
		if HasFlag(item) {
			return true
		}
	}
	return false
}
