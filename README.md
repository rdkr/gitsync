# gitsync
```
gitsync is a tool to keep local git repos in sync with remote git servers.

It supports individual repos and git service provider groups accessed over
HTTPS and authenticated either anonymously or with a token. Groups are
recursed to find projects and projects are concurrently cloned, pulled, or
fetched as appropriate.

Supported git service providers:
 - GitLab groups and repos over HTTPS (GITLAB_TOKEN env var should be set)
 - Anonymous repos over HTTPS

A .yaml config file is expected, and will be found from:
 - $HOME/.gitsync.yaml
 - $PWD/.gitsync.yaml
 - as specified using the --config/-c flag

The format of the config file is as follows:

gitlab:
  groups:
  - group: <group-id>
    location: <local path to sync to>
  projects:
  - url: <https clone url>
    location: <local path to sync to>
anon:
  projects:
  - url: <https clone url>
	location: <local path to sync to>

Usage:
  gitsync [flags]

Flags:
  -c, --config string   config file (default is $HOME/.gitsync.yaml)
  -h, --help            help for gitsync
  -v, --verbose         verbose Output
```
