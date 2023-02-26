package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"
)

// Keeping an eye on this struct. It and its logic seem to be growing and being more all over the place...
type migConfig struct {
	Sql     *sql.DB
	dialect string
	baseDir string
}

func NewMigration(sql *sql.DB, dialect string, baseDir string) *migConfig {
	//Append a '/' if the string is not empty and doesn't already end with a '/'.
	//This is to avoid files/dirs are created/read in and from unexpected places
	if baseDir != "" && baseDir[len(baseDir)-1:] != "/" {
		baseDir += "/"
	}
	return &migConfig{Sql: sql, dialect: dialect, baseDir: baseDir}
}

func (mig migConfig) Down() {
	var biggest int
	allMigrations := mig.GetAllMigrations()

	for i := 0; i < len(allMigrations); i++ {
		temp, _ := strconv.Atoi(allMigrations[i])
		if temp > biggest || biggest == 0 {
			biggest = temp
		}

		continue
	}

	mig.DownTo(string(rune(biggest)))
}

func (mig migConfig) DownTo(downto string) {
	mig.runSingleMigration(downto, "down")
}

func (mig migConfig) getRanMigrations() []string {

	var ranMigrations []string

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
		var currentMigration string
		err = result.Scan(&currentMigration)
		{
			ranMigrations = append(ranMigrations, currentMigration)
		}
		if err != nil {
			debug.PrintStack()
			fmt.Println(err)
		}

	}

	return ranMigrations

}

func (mig migConfig) HasMigrationRan(migrationToCheck string) bool {
	for _, item := range mig.getRanMigrations() {
		if item == migrationToCheck {
			return true
		}
	}

	return false

}
func (mig migConfig) GenerateMigration() {

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
func (mig migConfig) RunMigrations() {
	mig.ensureDirExists("migrations")
	for _, s := range mig.GetUnRanMigrations() {
		mig.runSingleMigration(s, "up")
	}
}

func (mig migConfig) runSingleMigration(s string, direction string) {
	migrationFile := mig.readFile(s, direction)

	fmt.Println("Currently executing: " + s)
	_, err := mig.Sql.Exec(string(migrationFile))
	if err != nil {
		fmt.Println("Error when running migConfig: ")
		fmt.Println(err)
		return
	}

	_, err = mig.Sql.Exec(GetCreateTableByDialect(mig.dialect))
	_, err = mig.Sql.Exec(InsertNewEntry(mig.dialect), s)

	if err != nil {
		debug.PrintStack()
		fmt.Println("when marking the migConfig as ran: ")
		fmt.Println(err)
	}
}

func (mig migConfig) GetAllMigrations() []string {
	var migrations = make([]string, 0)
	mig.ensureDirExists("migrations")
	dir, err := mig.readDir("migrations")
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	regexFofMigrationFile := regexp.MustCompile("\\d+")

	for _, value := range dir {

		if !value.IsDir() {
			continue
		}

		if !regexFofMigrationFile.MatchString(value.Name()) {
			continue
		}

		migrations = append(migrations, value.Name())
	}

	return migrations
}

func (mig migConfig) GetUnRanMigrations() []string {
	allMigs := mig.GetAllMigrations()
	ranMigs := mig.getRanMigrations()
	unranMigs := make([]string, 0)

	for i := 0; i < len(allMigs); i++ {
		appendItem := true
		for x := 0; x < len(ranMigs); x++ {
			if ranMigs[x] == allMigs[i] {
				appendItem = false
			}
		}

		if appendItem == true {
			unranMigs = append(unranMigs, allMigs[i])
		}
	}

	return unranMigs
}

func (mig migConfig) Status() {
	unranMigs := mig.GetUnRanMigrations()
	if len(unranMigs) == 0 {
		fmt.Println("all migrations have been ran")
		return
	}
	fmt.Println("these migrations were not ran")
	fmt.Println(unranMigs)
}

// Consider moving to its own context (file and receiver)
func (mig migConfig) ensureDirExists(dir string) {
	dir = fmt.Sprintf("%s%s", mig.baseDir, dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0771)
	}

}

func (mig migConfig) createEmptyFile(filePath string) {
	filePath = fmt.Sprintf("%s%s", mig.baseDir, filePath)
	_, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)
	}
}

func (mig migConfig) readDir(dir string) ([]os.DirEntry, error) {
	dir = fmt.Sprintf("%s%s", mig.baseDir, dir)
	return os.ReadDir(dir)
}

func (mig migConfig) readFile(migrationFile string, direction string) []byte {
	migrationFileContents, err := os.ReadFile(fmt.Sprintf("%s%s/%s/%s.sql", mig.baseDir, "migrations", migrationFile, direction))
	if err != nil {
		fmt.Println(err)
	}
	return migrationFileContents
}
