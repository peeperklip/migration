package migrations

import (
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"text/tabwriter"
)

type commandInterface interface {
	run(migconfig migConfig, arguments []string)
}

type generate struct{}
type migrate struct{}
type status struct{}
type revert struct{}
type down struct{}
type help struct{}

var commands = map[string]commandInterface{
	"generate": generate{},
	"migrate":  migrate{},
	"status":   status{},
	"revert":   revert{},
	"down":     down{},
	"":         help{},
	"help":     help{},
}

func (h help) run(migconfig migConfig, arguments []string) {
	outputHelpText()
}

func (c generate) run(migconfig migConfig, arguments []string) {
	migconfig.GenerateMigration()
}

func (m migrate) run(migconfig migConfig, arguments []string) {
	migconfig.RunMigrations()
}

func (s status) run(migconfig migConfig, arguments []string) {
	migconfig.Status()
}

func (r revert) run(migconfig migConfig, arguments []string) {
	if len(arguments) == 2 {
		outputHelpText()
		panic("No migConfig given to revert")
	}

	requestedMigration := arguments[2]

	migconfig.DownTo(requestedMigration)
}

func (d down) run(migconfig migConfig, arguments []string) {
	migconfig.Down()
}

func Init(migrate migConfig, command string) {
	args := os.Args

	for index, _ := range creationStatements {
		if index == command {
			break
		}

		outputHelpText()
		panic("This command is not supported")
	}

	commands[command].run(migrate, args)
}

func outputHelpText() {
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, _ = fmt.Println("Usage: go cli/migrations.go [COMMAND]")
	_, _ = fmt.Println()
	_, _ = fmt.Println("Available commands:")
	_, _ = fmt.Fprintln(w, "\tgenerate \t This generates a new migConfig")
	_, _ = fmt.Fprintln(w, "\tmigrate \t Run the migrations")
	_, _ = fmt.Fprintln(w, "\tdown \t the last ran migConfig according to down.sql")
	_, _ = fmt.Fprintln(w, "\trevert \t Revert to a specific migConfig")
	_, _ = fmt.Fprintln(w, "\tstatus \t Warn about any migConfig that hasn't run")

	_, _ = fmt.Fprintln(w)
	_ = w.Flush()
}
