package migrations

import (
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"text/tabwriter"
)

func Init(migrate migration) {
	args := os.Args

	if len(args) == 1 {
		outputHelpText()
		return
	}

	command := args[1]

	switch command {

	case "generate":

		migrate.GenerateMigration()
		break

	case "migrate":
		migrate.RunMigrations()
		break

	case "down":

		break

	case "revert_to":

		break

	case "status":

		break

	case "":
		outputHelpText()
		break
	default:
		outputHelpText()
		panic("This command is not supported")
	}
}

func outputHelpText() {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, _ = fmt.Println("Usage: go cli/migrations.go [COMMAND]")
	_, _ = fmt.Println()
	_, _ = fmt.Println("Available commands:")
	_, _ = fmt.Fprintln(w, "\tgenerate \t This generates a new migration")
	_, _ = fmt.Fprintln(w, "\tmigrate \t Run the migrations")
	_, _ = fmt.Fprintln(w, "\tdown \t Undo a specific migration according to down.sql")
	_, _ = fmt.Fprintln(w, "\trevert_to \t Undo all migrations according to down.sql until, but >not< including the specified migration")
	_, _ = fmt.Fprintln(w, "\tstatus \t Warn about any migration that hasn't run")

	_, _ = fmt.Fprintln(w)
	_ = w.Flush()
}
