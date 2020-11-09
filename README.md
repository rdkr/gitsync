# gitsync [![Build Status](https://travis-ci.org/rdkr/gitsync.svg)](https://travis-ci.org/rdkr/gitsync) [![codecov.io](https://codecov.io/github/rdkr/gitsync/coverage.svg)](https://codecov.io/github/rdkr/gitsync) [![Maintainability](https://api.codeclimate.com/v1/badges/c86f6cac36c28c9ea55f/maintainability)](https://codeclimate.com/github/rdkr/gitsync/maintainability) [![License](https://img.shields.io/github/license/rdkr/gitsync)](https://github.com/rdkr/gitsync/blob/master/LICENSE) [![Releases](https://img.shields.io/github/v/release/rdkr/gitsync)](https://github.com/rdkr/gitsync/releases)

gitsync is a tool to keep many local repos in sync with their remote hosts.

It supports recursively syncing GitHub orgs, teams, and users; GitLab groups; and individual repos. Repos are synced over HTTPS, optionally / where required using auth tokens.

## ⚠️ Breaking change warning!
The config file structure changed recently. The GitHub and GitLab sections are now lists of the original maps: https://github.com/rdkr/gitsync/commit/7135e5a38e087efaa2dbf3eedd94ea849812172e#

## Install

```
go install github.com/rdkr/gitsync
```

## Help text
```
gitsync is a tool to keep many local repos in sync with their remote hosts.

It supports recursively syncing GitHub orgs, teams, and users; GitLab groups; and individual
repos. Repos are synced over HTTPS, optionally / where required using auth tokens.

              Orgs'    Groups' / Teams'    Users'    Repos'
    GitHub      x              x             x
    GitLab                     x                       x
    HTTPS                                              x

Orgs / groups / user profiles are enumerated / recursed to find projects. All projects
are then concurrently synced, i.e:
 - cloned, if the local repo doesn't exist
 - pulled, if the local repo exists and is on main
 - fetched, if neither of the above

A .yaml config file is expected, The format of the config file is:

github:       # optional: defines GitHub resources
- baseurl:      # optional: a custom GitHub API URL
  token:        # required: a GitHub API token
  users:        # optional: defines GitHub users
  - name:         # required: GitHub username
    location:     # required: local path to sync to
  orgs:         # optional: defines GitHub organisations
  - name:         # required: GitHub org name
    location:     # required: local path to sync to
  teams:        # optional: defines GitHub teams
  - org:          # required: GitHub org name
    name:         # required: GitHub team name
    location:     # required: local path to sync to
gitlab:       # optional: defines GitLab resources
- baseurl:      # optional: a custom GitLab API URL
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

Usage:
  gitsync [flags]

Flags:
  -c, --config string   config file location
  -d, --debug           debug output (implies verbose)
  -h, --help            help for gitsync
  -v, --verbose         verbose output instead of pretty output
```

## Contributing
Contributions are welcomed, especially to complete and add new providers! Feel free to submit an issue or PR :)
