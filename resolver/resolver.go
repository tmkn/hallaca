package resolver

import (
	"sort"
	"strings"

	"maps"

	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/log"
	"github.com/tmkn/hallaca/pkg"
	"github.com/tmkn/hallaca/provider"
)

type Resolver interface {
	Resolve(name string, version string, options Options) *pkg.Pkg
}

type Options struct {
	Provider provider.Provider
	Depth    uint32
}

type ResolverQueueItem struct {
	Name            string
	Version         string
	VisitedPackages map[string]interface{}
	Pkg             *pkg.Pkg
}

func (item ResolverQueueItem) String() string {
	return item.Name + "@" + item.Version
}

type StandardResolver struct{}

func (r *StandardResolver) Resolve(name string, version string, options Options) *pkg.Pkg {
	if version == "" {
		log.Fatalf("latest version feature is not yet supported")
	}

	root := &pkg.Pkg{}
	queue := []ResolverQueueItem{{
		Name:            name,
		Version:         version,
		VisitedPackages: make(map[string]interface{}),
		Pkg:             root,
	}}
	dependencyKey := "dependencies"

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		log.Infof("Evaluating %s", item)

		versions, err := options.Provider.GetVersions(item.Name)
		if err != nil {
			log.Fatalf("couldn't get versions for %s", item)
		}

		item.Version = ResolveVersion(item.Version, versions)
		item.Pkg.Version = item.Version
		item.Pkg.Name = item.Name

		itemID := item.String()
		item.VisitedPackages[itemID] = struct{}{}

		metadata, err := options.Provider.GetPackageMetadata(item.Name, item.Version)
		if err != nil {
			log.Fatalf("couldn't get metadata for %s", item)
		}

		dependencies, ok := metadata[dependencyKey].(map[string]interface{})

		log.Infof("Found %d dependencies for %s", len(dependencies), item.Name)

		if !ok {
			continue
		}

		for key, value := range dependencies {
			var isAliased bool = false
			var nameToResolve = key
			versionToResolve, ok := value.(string)
			if !ok {
				log.Fatalf("Value for key %s is not a string\n", key)
			}

			if strings.HasPrefix(versionToResolve, "npm:") {
				packageSpec := strings.TrimPrefix(versionToResolve, "npm:")

				if _name, _version, error := pkg.ParsePackageSpec(packageSpec); error != nil {
					log.Fatalf("Couldn't parse package spec %q from %q@%q", packageSpec, key, versionToResolve)

				} else {
					isAliased = true
					nameToResolve = _name
					versionToResolve = _version
				}
			}

			depVersions, err := options.Provider.GetVersions(nameToResolve)
			if err != nil {
				log.Fatalf("couldn't get versions for dependency %s: %v", nameToResolve, err)
			}

			resolvedVersion := ResolveVersion(versionToResolve, depVersions)
			dependencyID := nameToResolve + "@" + resolvedVersion

			depPkg := &pkg.Pkg{
				Name:    nameToResolve,
				Version: resolvedVersion,
				Parent:  item.Pkg,
			}

			if isAliased {
				depPkg.Name = key
				depPkg.AliasedFrom = nameToResolve
			}

			item.Pkg.Dependencies = append(item.Pkg.Dependencies, depPkg)

			log.Infof("Evaluated %s", depPkg)

			if _, exists := item.VisitedPackages[dependencyID]; exists {
				log.Warnf("Detected loop for %s", dependencyID)
				depPkg.IsLoop = true

				continue
			}

			newVisited := make(map[string]interface{})
			maps.Copy(newVisited, item.VisitedPackages)
			newVisited[dependencyID] = struct{}{}

			dependency := ResolverQueueItem{
				Name:            nameToResolve,
				Version:         resolvedVersion,
				VisitedPackages: newVisited,
				Pkg:             depPkg,
			}
			queue = append(queue, dependency)
		}
	}

	return root
}

func ResolveVersion(toResolve string, versions []string) string {
	constraint, err := semver.NewConstraint(toResolve)
	if err != nil {
		log.Fatalf("couldn't create version constraint: %s", toResolve)
	}

	var matchingVersions []*semver.Version
	for _, versionString := range versions {
		version, err := semver.NewVersion(versionString)
		if err != nil {
			log.Fatalf("invalid available version: %v", err)
		}

		if constraint.Check(version) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))

	if len(matchingVersions) > 0 {
		return matchingVersions[0].String()
	}

	log.Fatalf("Couldn't resolve version %s", toResolve)
	return ""
}
