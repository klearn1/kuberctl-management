/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"k8s.io/publishing-bot/cmd/publishing-bot/config"
)

const (
	stagingDirectory    = "staging/"
	rulesFile           = stagingDirectory + "publishing/rules.yaml"
	componentsDirectory = stagingDirectory + "src/k8s.io/"
)

func getGoModDependencies(dir string) (map[string][]string, error) {
	allDependencies := make(map[string][]string)
	components, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, component := range components {
		componentName := component.Name()
		if !component.IsDir() {
			// currently there is no hard check that the staging directory should not contain
			// other files
			continue
		}
		gomodFilePath := filepath.Join(dir, componentName, "go.mod")
		gomodFileContent, err := os.ReadFile(gomodFilePath)
		if err != nil {
			return nil, err
		}

		fmt.Printf("%s dependencies", componentName)

		gomodFile, err := modfile.ParseLax(gomodFilePath, gomodFileContent, nil)
		if err != nil {
			return nil, err
		}
		// get all the other dependencies from within staging, i.e all the modules in replace
		// section
		for _, module := range gomodFile.Replace {
			dep := strings.TrimPrefix(module.Old.Path, "k8s.io/")
			allDependencies[componentName] = append(allDependencies[componentName], dep)
		}
	}
	return allDependencies, nil
}

// diffSlice returns the difference of s1-s2
func diffSlice(s1, s2 []string) []string {
	var diff []string
	set := make(map[string]struct{}, len(s2))
	for _, s := range s2 {
		set[s] = struct{}{}
	}
	for _, s := range s1 {
		if _, ok := set[s]; !ok {
			diff = append(diff, s)
		}
	}
	return diff
}

func getKeys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	rules, err := config.LoadRules(rulesFile)
	if err != nil {
		os.Exit(1)
	}

	gomodDependencies, err := getGoModDependencies(stagingDirectory)

	var processedRepos []string
	for _, rule := range rules.Rules {
		mainBranch := rule.Branches[0]
		// CHeck 1
		// if this no longer exists in master
		if _, ok := gomodDependencies[rule.DestinationRepository]; !ok {
			// make sure we dont include a rule to publish it from master
			for _, branch := range rule.Branches {
				if branch.Name == "master" {
					err := fmt.Errorf("cannot find master branch for destination `%s`", rule.DestinationRepository)
					panic(err)
				}
			}
			// and skip the validation of publishing rules for it
			continue
		}

		// check 2, source directory checks
		for _, branch := range rule.Branches {
			if branch.Source.Dir != "" {
				err := fmt.Errorf("use of deprecated `dir` field in rules for `%s`", rule.DestinationRepository)
				panic(err)
			}
			if len(branch.Source.Dirs) > 1 {
				err := fmt.Errorf("cannot have more than one directory (%s) per source branch `%s` of `%s`",
					branch.Source.Dirs,
					branch.Source.Branch,
					rule.DestinationRepository,
				)
				panic(err)
			}
			if !strings.HasSuffix(branch.Source.Dirs[0], rule.DestinationRepository) {
				err := fmt.Errorf("copy/paste error `%s` refers to `%s`", rule.DestinationRepository, branch.Source.Dirs[0])
				panic(err)
			}
		}

		// check 3
		if mainBranch.Name != "master" {
			err := fmt.Errorf("cannot find master branch for destination `%s`", rule.DestinationRepository)
			panic(err)
		}

		// check 4
		if mainBranch.Source.Branch != "master" {
			err := fmt.Errorf("cannot find master source branch for destination `%s`", rule.DestinationRepository)
			panic(err)
		}

		// check 5
		// we specify the go version for all master branches through `default-go-version`
		// so ensure we don't specify explicit go version for master branch in rules
		if mainBranch.GoVersion != "" {
			err := fmt.Errorf("go version must not be specified for master branch for destination `%s`", rule.DestinationRepository)
			panic(err)
		}

		fmt.Printf("processing : %s", rule.DestinationRepository)
		if _, ok := gomodDependencies[rule.DestinationRepository]; !ok {
			err := fmt.Errorf("missing go.mod for `%s`", rule.DestinationRepository)
			panic(err)
		}
		processedRepos = append(processedRepos, rule.DestinationRepository)
		var processedDeps []string
		for _, dep := range gomodDependencies[rule.DestinationRepository] {
			found := false
			if len(mainBranch.Dependencies) > 0 {
				for _, dep2 := range mainBranch.Dependencies {
					processedDeps = append(processedDeps, dep2.Repository)
					if dep2.Branch != "master" {
						err := fmt.Errorf("looking for master branch of %s and found : %s for destination", dep2.Repository, rule.DestinationRepository)
						panic(err)
					}
					found = dep2.Repository == dep
				}
			} else {
				err := fmt.Errorf("Please add %s as dependencies under destination %s", gomodDependencies[rule.DestinationRepository], rule.DestinationRepository)
				panic(err)
			}
			if !found {
				err := fmt.Errorf("Please add %s as a dependency under destination %s", dep, rule.DestinationRepository)
				panic(err)
			} else {
				fmt.Printf("dependency %s found\n", dep)
			}
		}
		// check if all deps are processed.
		extraDeps := diffSlice(processedDeps, gomodDependencies[rule.DestinationRepository])
		if len(extraDeps) > 0 {
			err := fmt.Errorf("extra dependencies in rules for %s: %s", rule.DestinationRepository, strings.Join(extraDeps, ","))
			panic(err)
		}
	}
	// check if all repos are processed.
	items := diffSlice(getKeys(gomodDependencies), processedRepos)
	if len(items) > 0 {
		err := fmt.Errorf("missing rules for %s", strings.Join(items, ","))
		panic(err)
	}
}
