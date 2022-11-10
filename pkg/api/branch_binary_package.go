package api

type Package struct {
	Name      string `json:"name"`
	Epoch     int    `json:"epoch"`
	Version   string `json:"version"`
	Release   string `json:"release"`
	Arch      string `json:"arch"`
	DistTag   string `json:"disttag"`
	BuildTime int    `json:"buildtime"`
	Source    string `json:"source"`
}

type Branch struct {
	RequestArgs struct{}  `json:"request_args"`
	Length      int       `json:"length"`
	Packages    []Package `json:"packages"`
}

type BranchDiff struct {
	UniquePackages1    []Package `json:"unique_packages_1"`
	UniquePackages2    []Package `json:"unique_packages_2"`
	NewerPackagesFrom1 []Package `json:"newer_packages_from_1"`
}
