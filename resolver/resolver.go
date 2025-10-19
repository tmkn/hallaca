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
	Resolve(name string, version string, options Options) (*pkg.Pkg, error)
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

func (item *ResolverQueueItem) String() string {
	return item.Name + "@" + item.Version
}

type StandardResolver struct{}

func (r *StandardResolver) Resolve(name string, version string, options Options) (*pkg.Pkg, error) {
	if version == "" {
		return nil, &ResolverError{Msg: "latest version feature is not yet supported"}
	}

	root := &pkg.Pkg{}
	queue := []ResolverQueueItem{{
		Name:            name,
		Version:         version,
		VisitedPackages: make(map[string]interface{}),
		Pkg:             root,
	}}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if pkg.GetDepth(item.Pkg) > uint(options.Depth) {
			continue
		}

		log.Infof("Evaluating %s", &item)

		versions, err := options.Provider.GetVersions(item.Name)
		if err != nil {
			return nil, &ResolverError{Msg: "couldn't get versions for %s", Pkg: &item}
		}

		resolvedVersion, err := ResolveVersion(item.Version, versions)
		if err != nil {
			return nil, err
		}
		item.Version = resolvedVersion
		item.Pkg.Version = item.Version
		item.Pkg.Name = item.Name

		itemID := item.String()
		item.VisitedPackages[itemID] = struct{}{}

		metadata, err := options.Provider.GetPackageMetadata(item.Name, item.Version)
		if err != nil {
			return nil, &ResolverError{Msg: "couldn't get metadata for %s", Pkg: &item}
		}

		item.Pkg.Metadata = metadata

		log.Infof("Found %d dependencies for %s", len(metadata.Dependencies), item.Name)

		for key, value := range metadata.Dependencies {
			var isAliased bool = false
			var nameToResolve = key
			versionToResolve := value

			if strings.HasPrefix(versionToResolve, "npm:") {
				packageSpec := strings.TrimPrefix(versionToResolve, "npm:")

				if _name, _version, err := pkg.ParsePackageSpec(packageSpec); err != nil {
					return nil, &ResolverError{Msg: "Couldn't parse package spec %q from %q@%q", Pkg: &item}
				} else {
					isAliased = true
					nameToResolve = _name
					versionToResolve = _version
				}
			}

			depVersions, err := options.Provider.GetVersions(nameToResolve)
			if err != nil {
				return nil, &ResolverError{Msg: "couldn't get versions for dependency %s: %v", Pkg: &item}
			}

			resolvedVersion, err := ResolveVersion(versionToResolve, depVersions)
			if err != nil {
				return nil, err
			}
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

	return root, nil
}

func ResolveVersion(toResolve string, versions []string) (string, error) {
	constraint, err := semver.NewConstraint(toResolve)
	if err != nil {
		return "", &ResolverError{Msg: "couldn't create version constraint: %s", Pkg: nil}
	}

	var matchingVersions []*semver.Version
	for _, versionString := range versions {
		version, err := semver.NewVersion(versionString)
		if err != nil {
			return "", &ResolverError{Msg: "invalid available version: %v", Pkg: nil}
		}

		if constraint.Check(version) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))

	if len(matchingVersions) > 0 {
		return matchingVersions[0].String(), nil
	}

	return "", &ResolverError{Msg: "Couldn't resolve version %s", Pkg: nil}
}
