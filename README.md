# Migration

### Go get
```shell
go get https://github.com/peeperklip/migration@0.0.1
# add the -u flag for updating
```

### (proposal on how to) Intergrate in your own project
1. Get according go get paragraph
2. Create a directory `cli`
   * as per https://github.com/golang-standards/project-layout/tree/master/cmd
3. create a file + function that you can exectute command line
```
//Exexute in: $PROJECT_ROOT
//File: $PROJECT_ROOT/cli/migration.go
//Execute as: 
func main() {
	mig := migrations.NewMigration(
		utils.CreateConnection(),
		"postgress",
		".")
	migrations.Init(*mig)
}
```


### Architecture:
<b>app.go</b> Is the entry point for the CLI<br>
<b>dialect.go</b> Will eventually be used to support multiple SQL dialects<br>
<b>migration.go</b> Holds all the logic for managing migrations<br>

### Codestyle:
structs go first, interfaces second, then the methods, then general functions. Besides that it's pretty much just `gofmt .`

### To be improved in 0.0.2 and 0.0.3:
* The unsustainable swtich case in dialect.go
* The typejuggeling througout migration.go
* The last few methods in migration.go should be moved to its own file and struct
* Inject a logger
* unify output and make more of the underlying code swapable
* general improvements throughout as well as a better distinction between exported and unexported methods