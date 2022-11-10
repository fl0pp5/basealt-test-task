package branch

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"repocmp/pkg/api"
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
