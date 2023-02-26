package migrations

import (
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"text/tabwriter"
)

func Init(migrate migConfig) {
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
		migrate.Down()
		break

	case "revert":
		if len(args) == 2 {
			outputHelpText()
			panic("No migConfig given to revert")
		}

		requestedMigration := args[2]

		migrate.DownTo(requestedMigration)
		break

	case "status":
		migrate.Status()
		break
	default:
		outputHelpText()
		panic("This command is not supported")
	}
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
