# gitsync [![Build Status](https://travis-ci.org/rdkr/gitsync.svg)](https://travis-ci.org/rdkr/gitsync) [![codecov.io](https://codecov.io/github/rdkr/gitsync/coverage.svg)](https://codecov.io/github/rdkr/gitsync) [![Maintainability](https://api.codeclimate.com/v1/badges/c86f6cac36c28c9ea55f/maintainability)](https://codeclimate.com/github/rdkr/gitsync/maintainability) [![License](https://img.shields.io/github/license/rdkr/gitsync)](https://github.com/rdkr/gitsync/blob/master/LICENSE) [![Releases](https://img.shields.io/github/v/release/rdkr/gitsync)](https://github.com/rdkr/gitsync/releases)

gitsync is a tool to keep many local repos in sync with their remote hosts.

It supports syncing GitHub / GitLab users' repos, recursively syncing GitHub / GitLab
groups, and syncing individual Git repos, all over HTTPS and optionally using auth tokens.

## Install

```
go install github.com/rdkr/gitsync
```

### MacOS
```
brew tap rdkr/taps
brew install gitsync
```

## Help text
```
gitsync is a tool to keep many local repos in sync with their remote hosts.

It supports syncing GitHub / GitLab users' repos, recursively syncing GitHub / GitLab
groups, and syncing generic Git repos, all over HTTPS and optionally using auth tokens.

                Users'   Groups'   Individual
    GitHub        x
    GitLab                  x          x
    Generic                            x

Groups are recursed to find projects. All projects are concurrently synced, i.e:
 - cloned, if the local repo doesn't exist
 - pulled, if the local repo exists and is on master
 - fetched, if neither of the above

A .yaml config file is expected, The format of the config file is:

github:       # optional: defines GitHub resources
  token:        # required: a GitHub API token
  users:        # optional: defines GitHub users
  - name:         # required: GitHub username
    location:     # required: local path to sync to
gitlab:       # optional: defines GitLab resources
  baseurl:      # optional: a custom GitLab API URL
  token:        # optional: a GitLab API token
  groups:       # optional: defines GitLab groups
  - group:        # required: group ID number
    location:     # required: local path to sync to
  projects:     # optional: defines GitLab projects
  - url:          # required: https clone url
    location:     # required: local path to sync to
anon:         # optional: defines any other resources
  projects:     # optional: defines any HTTPS projects
  - url:          # required: https clone url
    location:     # required: local path to sync to
    token:        # optional: HTTPS token to use

The config file will will be found, by order of precedence, from:
 - $HOME/.gitsync.yaml
 - $PWD/.gitsync.yaml
 - as specified using the --config/-c flag

Treat this file with care, as it may contain secrets.
```

## Contributing
Contributions are welcomed, especially to complete and add new providers! Feel free to submit an issue or PR :)
