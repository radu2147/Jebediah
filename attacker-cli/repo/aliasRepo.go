package repo

import "fmt"

type AliasRepo struct {
	Aliases map[string]string
}

func (repo *AliasRepo) PrintAll() {
	for k, v := range repo.Aliases {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func (repo *AliasRepo) Add(key, value string) {
	repo.Aliases[key] = value
}

func (repo *AliasRepo) GetVictim(victim string) string {
	if repo.Aliases[victim] == "" {
		return victim
	}
	return repo.Aliases[victim]
}
