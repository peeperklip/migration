# Migration
## Ignore for now as it's not production ready

Initially started within a separate project of mine because of a desire for a migration tool
[See commit and message 06048d639eee4c994544091f0e8b4a255c021554](https://github.com/peeperklip/migration/commit/06048d639eee4c994544091f0e8b4a255c021554) From there on its restructuring, testing and developing

Made public for my own testing/developing purposes as it's slightly less of a hassle to `go get` a public repository. Canges and features are made if i require them in my own projects

### Go get
```shell
go get https://github.com/peeperklip/migration@{COMMIT_HASH}
# add the -u flag for updating
```

### (proposal on how to) Intergrate in your own project
1. Get according go get paragraph
2. Create a directory `cli`
   * as per https://github.com/golang-standards/project-layout/tree/master/cmd
3. create a file + function that you can exectute command line
4. 


### Architecture:
<b>main.go</b> Is the entry point for the CLI<br>
<b>dialect.go</b> Will eventually be used to support multiple SQL dialects<br>
<b>migration.go</b> Holds all the logic for managing migrations<br>
<b>dbUtils.go</b> In there to do more supporting tasks<br>

### Codestyle:
structs go first, interfaces second, then the methods, then general functions. Besides that it's pretty much just `gofmt .`

### To be improved
* The unsustainable swtich case in dialect.go
* The typejuggeling througout migration.go
* The last few methods in migration.go should be moved to its own file and struct
* Inject a logger
* unify output and make it swapable