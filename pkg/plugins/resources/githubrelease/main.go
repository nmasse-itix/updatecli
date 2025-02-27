package githubrelease

import (
	"github.com/mitchellh/mapstructure"
	"github.com/updatecli/updatecli/pkg/plugins/scms/github"
	"github.com/updatecli/updatecli/pkg/plugins/utils/version"
)

// Spec defines a specification for a "gittag" resource
// parsed from an updatecli manifest file
type Spec struct {
	Owner         string         // Owner specifies repository owner
	Repository    string         // Repository specifies the name of a repository for a specific owner
	Token         string         // Token specifies the credential used to authenticate with
	URL           string         // URL specifies the default github url in case of GitHub enterprise
	Username      string         // Username specifies the username used to authenticate with Github API
	VersionFilter version.Filter // VersionFilter provides parameters to specify version pattern and its type like regex, semver, or just latest.
}

// GitHubRelease defines a resource of kind "githubrelease"
type GitHubRelease struct {
	ghHandler     github.GithubHandler
	versionFilter version.Filter // Holds the "valid" version.filter, that might be different than the user-specified filter (Spec.VersionFilter)
	foundVersion  version.Version
}

// New returns a new valid GitHubRelease object.
func New(spec interface{}) (*GitHubRelease, error) {
	newSpec := Spec{}

	err := mapstructure.Decode(spec, &newSpec)
	if err != nil {
		return &GitHubRelease{}, err
	}

	newHandler, err := github.New(github.Spec{
		Owner:      newSpec.Owner,
		Repository: newSpec.Repository,
		Token:      newSpec.Token,
		URL:        newSpec.URL,
		Username:   newSpec.Username,
	})
	if err != nil {
		return &GitHubRelease{}, err
	}

	newFilter, err := newSpec.VersionFilter.Init()
	if err != nil {
		return &GitHubRelease{}, err
	}

	return &GitHubRelease{
		ghHandler:     newHandler,
		versionFilter: newFilter,
	}, nil
}
