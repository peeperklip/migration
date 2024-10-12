package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"

	_ "github.com/lib/pq"
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
	migconfig.generateMigration()
}

func (m migrate) run(migconfig migConfig, arguments []string) {
	migconfig.runMigrations()
}

func (s status) run(migconfig migConfig, arguments []string) {
	migconfig.status()
}

func (r revert) run(migconfig migConfig, arguments []string) {
	if len(arguments) == 2 {
		outputHelpText()
		panic("No migration given to revert")
	}

	requestedMigration := arguments[2]

	migconfig.downTo(requestedMigration)
}

func (d down) run(migconfig migConfig, arguments []string) {
	migconfig.down()
}

// Init is upposed to be called from the package this package is used in
func Init(migrate migConfig) {
	args := os.Args

	if len(args) == 0 {
		outputHelpText()
		return
	}

	for registeredCommand := range commands {
		if registeredCommand == args[1] {
			commands[registeredCommand].run(migrate, args)
			return
		}

	}
	panic("This command is not supported")
}

func Migrate(sql *sql.DB, dialect string, baseDir string) {
	m := NewMigration(sql, dialect, baseDir)
	commands["migrate"].run(*m, []string{})
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
