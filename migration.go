package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Keeping an eye on this struct. It and its logic seem to be growing and being more all over the place...
type migration struct {
	Sql     *sql.DB
	dialect string
	baseDir string
}

func NewMigration(sql *sql.DB, dialect string, baseDir string) *migration {
	//Append a '/' if the string is not empty and doesn't already end with a '/'.
	//This is to avoid files/dirs are created/read in and from unexpected places
	if baseDir != "" && baseDir[len(baseDir)-1:] != "/" {
		baseDir += "/"
	}
	return &migration{Sql: sql, dialect: dialect, baseDir: baseDir}
}

func (mig migration) getRanMigrations() []int64 {

	var ranMigrations []int64

	result, err := mig.Sql.Query(QueryForRanMigrations(mig.dialect))

	if err != nil {
		return ranMigrations
	}

	defer func(result *sql.Rows) {
		err = result.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(result)

	for result.Next() {
		var currentMigration int64
		err = result.Scan(&currentMigration)
		{
			ranMigrations = append(ranMigrations, currentMigration)
		}
		if err != nil {
			fmt.Println(err)
		}

	}

	return ranMigrations

}

func (mig migration) prepareMigrationsTable() {
	_, err := mig.Sql.Exec(GetCreateTableByDialect(mig.dialect))

	if err != nil {
		fmt.Println(err)
	}
}
func (mig migration) HasMigrationRan(migrationToCheck string) bool {
	for _, item := range mig.getRanMigrations() {
		if strconv.FormatInt(item, 10) == migrationToCheck {
			return true
		}
	}

	return false

}
func (mig migration) GenerateMigration() {

	currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	migrationDirName := fmt.Sprintf("migrations/%s", currentTimestamp)

	mig.ensureDirExists("")
	mig.ensureDirExists("migrations")
	mig.ensureDirExists(migrationDirName)

	upMigFilename := fmt.Sprintf("%s/up.sql", migrationDirName)
	downMigFilename := fmt.Sprintf("%s/down.sql", migrationDirName)

	mig.createEmptyFile(upMigFilename)
	mig.createEmptyFile(downMigFilename)

}
func (mig migration) RunMigrations() {
	mig.ensureDirExists("migrations")
	dir, err := mig.readDir("migrations")
	if err != nil {
		return
	}

	regexFofMigrationFile := regexp.MustCompile("\\d+")

	for _, value := range dir {

		if !value.IsDir() {
			continue
		}

		if !regexFofMigrationFile.MatchString(value.Name()) {
			continue
		}

		if !mig.HasMigrationRan(value.Name()) {

			migrationFile := mig.readFile(value.Name())

			fmt.Println("Currently executing: " + value.Name())
			_, err := mig.Sql.Exec(string(migrationFile))
			if err != nil {
				fmt.Println("Error when running migration: ")
				fmt.Println(err)
			}

			_, err = mig.Sql.Exec(InsertNewEntry(mig.dialect), value.Name())

			if err != nil {
				fmt.Println("when marking the migration as ran: ")
				fmt.Println(err)
			}
		}
	}
}

// Consider moving to its own context (file and receiver)
func (mig migration) ensureDirExists(dir string) {
	dir = fmt.Sprintf("%s%s", mig.baseDir, dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err != nil {
			fmt.Println(err)
		}
		err = os.Mkdir(dir, 0771)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func (mig migration) createEmptyFile(filePath string) {
	filePath = fmt.Sprintf("%s%s", mig.baseDir, filePath)
	_, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)
	}
}

func (mig migration) readDir(dir string) ([]os.DirEntry, error) {
	dir = fmt.Sprintf("%s%s", mig.baseDir, dir)
	return os.ReadDir(dir)
}

func (mig migration) readFile(migrationFile string) []byte {
	migrationFileContents, err := os.ReadFile(fmt.Sprintf("%s%s/%s/up.sql", mig.baseDir, "migrations", migrationFile))
	if err != nil {
		fmt.Println(err)
	}
	return migrationFileContents
}
