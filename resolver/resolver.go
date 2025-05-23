package resolver

import (
	"log"
	"sort"

	"github.com/Masterminds/semver/v3"
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
		log.Fatalln("latest version feature is not yet supported")
	}

	var root *pkg.Pkg = &pkg.Pkg{}
	var queue []ResolverQueueItem = []ResolverQueueItem{{Name: name, Version: version, VisitedPackages: make(map[string]interface{}), Pkg: root}}
	var dependencyKey = "dependencies"
	// var visitedPackages map[string]interface{} = make(map[string]interface{})
	// var currentDepth = 0

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		log.Println("Evaluating", item)

		versions, err := options.Provider.GetVersions(item.Name)

		if err != nil {
			log.Fatalln("couldn't get versions for", item)
		}

		item.Version = ResolveVersion(item.Version, versions)
		item.Pkg.Version = item.Version
		item.Pkg.Name = item.Name

		item.VisitedPackages[item.String()] = struct{}{}

		metadata, err := options.Provider.GetPackageMetadata(item.Name, item.Version)

		if err != nil {
			log.Fatalln("couldn't get metadata for", item)
		}

		dependencies, ok := metadata[dependencyKey].(map[string]interface{})

		if !ok {
			// log.Fatalln("couldn't cast dependencies for", item.Name, item.Version)
			// log.Println("no dependencies for", item.Name, item.Version)
		} else {
			log.Println("found", len(dependencies), "dependencies for", item.Name)

			for key, value := range dependencies {
				strValue, ok := value.(string)
				if !ok {
					log.Fatalf("Value for key %s is not a string\n", key)
				}

				depPkg := &pkg.Pkg{Parent: item.Pkg}
				item.Pkg.Dependencies = append(item.Pkg.Dependencies, depPkg)

				dependency := ResolverQueueItem{Name: key, Version: strValue, VisitedPackages: make(map[string]interface{}), Pkg: depPkg}
				// todo implement loop logic
				if _, exists := item.VisitedPackages[dependency.String()]; !exists {
					queue = append(queue, dependency)
				} else {
					log.Fatal("Found loop, todo set loop flag")
				}

			}
		}

	}

	return root
}

func ResolveVersion(toResolve string, versions []string) string {
	constraint, err := semver.NewConstraint(toResolve)

	if err != nil {
		log.Fatalln("couldn't create version constraint:", toResolve)
	}

	var matchingVersions []*semver.Version

	for _, versionString := range versions {
		version, err := semver.NewVersion(versionString)
		if err != nil {
			log.Fatalf("invalid available version: %w", err)
		}

		if constraint.Check(version) {
			matchingVersions = append(matchingVersions, version)
		}

	}

	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))
	// result := make([]string, len(matchingVersions))

	if len(matchingVersions) > 0 {
		return matchingVersions[0].String()
	}

	log.Fatalln("Couldn't resolve version", toResolve)

	return ""
}
