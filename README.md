# Migration
## Ignore for now as it's not production ready

Initially started within a separate project of mine because of a desire for a migration tool
[See commit and message 06048d639eee4c994544091f0e8b4a255c021554](https://github.com/peeperklip/migration/commit/06048d639eee4c994544091f0e8b4a255c021554) From there on its restructuring, testing and developing

Made public for my own testing/developing purposes as it's slightly less of a hassle to `go get` a public repository

[All updates are being done in the develop branch](https://github.com/peeperklip/migration/tree/develop)

### Architecture:
main.go:main is the entry point for the CLI
dialect.go will eventually be used to support multiple SQL dialects
migration.go holds all the logic for managing migrations
dbUtils.go in there to do more supporting tasks

### Codestyle:
structs go first, interfaces second, then the methods, then general functions. Besides that it's pretty much just `gofmt .`

### Go get
```shell
go get https://github.com/peeperklip/migration@develop
# add the -u flag for updating
```