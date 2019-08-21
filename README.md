# gitsync [![Build Status](https://travis-ci.org/rdkr/gitsync.svg)](https://travis-ci.org/rdkr/gitsync) [![codecov.io](https://codecov.io/github/rdkr/gitsync/coverage.svg)](https://codecov.io/github/rdkr/gitsync)
```
gitsync is a tool to keep local git repos in sync with remote git servers.

It supports individual repos and git service provider groups accessed over
HTTPS and authenticated either anonymously or with a token.

Groups are recursed to find projects and projects are concurrently cloned,
pulled, or fetched as appropriate.

Supported git service providers:
 - GitLab groups and repos over HTTPS
 - Anonymous repos over HTTPS

A .yaml config file is expected, The format of the config file is:

gitlab:         # optional: defines GitLab resources
  token:        # required: a GitLab API token
  groups:       # optional: defines GitLab groups
  - group:      # required: group ID number
    location:   # required: local path to sync to
  projects:     # optional: defines GitLab projects
  - url:        # required: https clone url
    location:   # required: local path to sync to
anon:           # optional: defines any other resources
  projects:     # optional: defines any HTTPS projects
  - url:        # required: https clone url
    location:   # required: local path to sync to
    token:      # optional: HTTPS token to use

The config file will will be found, by order of precedence, from:
 - $HOME/.gitsync.yaml
 - $PWD/.gitsync.yaml
 - as specified using the --config/-c flag

Treat this file with care, as it may contain secrets.

Usage:
  gitsync [flags]

Flags:
  -c, --config string   config file location
  -h, --help            help for gitsync
  -v, --verbose         verbose / script friendly output
```
