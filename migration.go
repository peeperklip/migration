package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/peeperklip/migration/internal"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"
)

type migConfig struct {
	Sql     *sql.DB
	dialect string
	baseDir string
}

var stateMap = [3]string{"RAN", "UNRAN", "REVERTED"}

type migration struct {
	id    string
	state string
}

func createNewMigration(migrationId string, stateString string) migration {
	var migration migration

	migration.id = migrationId
	migration.setState(stateString)
	return migration
}

func (migration *migration) setState(stateString string) {
	if migration.state == "UNRAN" && stateString == "REVERTED" {
		panic("This could never happen")
	}

	if migration.state == "REVERTED" && stateString == "UNRAN" {
		panic("This could never happen")
	}

	for _, val := range stateMap {
		if val == stateString {
			migration.state = stateString
		}
	}

	//Trigger error?
}

// table > migration
// 1676928405
// 1677431029
// 1677431498
// 1677431029

func loadMigrations(migrator migConfig) []migration {
	migrator.ensureDirExists("migrations")
	migrator.ensureTableExists()
	miglist := make([]migration, 0)
	ranMigs := migrator.getRanMigrations()
	allMigs := migrator.getAllMigrations()

mainLoop:
	for _, mig := range allMigs {
		for _, ranMig := range ranMigs {
			if ranMig.id == mig {
				continue mainLoop
			}
		}

		//The migrations that have been ran before already had their statusses and are bing tracked
		//The untracked ones are by definition unran
		miglist = append(miglist, createNewMigration(mig, "UNRAN"))

	}

	return miglist
}

//func NewMigration(sql *sql.DB, dialect string, baseDir string) *migConfig {
//	//Append a '/' if the string is not empty and doesn't already end with a '/'.
//	//This is to avoid files/dirs are created/read in and from unexpected places
//	if baseDir != "" && baseDir[len(baseDir)-1:] != "/" {
//		baseDir += "/"
//	}
//	return &migConfig{Sql: sql, dialect: dialect, baseDir: baseDir}
//}

func (mig migConfig) Down() {
	var biggest int
	allMigrations := mig.getAllMigrations()

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
	migs := loadMigrations(mig)
	for _, m := range migs {
		if m.id != downto {
			continue
		}

		if m.state == "RAN" || m.state == "REVERTED" {
			break
		}

		mig.runSingleMigration(&m, "down")
	}
}

func (mig migConfig) getRanMigrations() []migration {

	var ranMigrations []migration

	result, err := mig.Sql.Query(QueryForRanMigrations(mig.dialect))

	if err != nil {
		internal.AddError(err)
		return ranMigrations
	}

	defer func(result *sql.Rows) {
		err = result.Close()
		internal.AddError(errors.New("error when selecting migrations"))
		internal.AddError(err)
	}(result)

	for result.Next() {
		var currentMigration migration
		err = result.Scan(&currentMigration.id, &currentMigration.state)
		{
			ranMigrations = append(ranMigrations, currentMigration)
		}
		internal.AddError(err)
	}

	return ranMigrations
}

func (mig migConfig) HasMigrationRan(migrationToCheck string) bool {
	for _, item := range mig.getRanMigrations() {
		if item.id == migrationToCheck {
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
	migrations := loadMigrations(mig)

	for _, m := range migrations {
		if m.state == "REVERTED" {
			continue
		}

		mig.runSingleMigration(&m, "up")
	}

}

func (mig migConfig) runSingleMigration(s *migration, direction string) {
	migrationFile := mig.readFile(s.id, direction)

	log.Println("Currently executing: " + s.id)
	_, err := mig.Sql.Exec(string(migrationFile))
	if err != nil {
		internal.AddError(err)
		return
	}

	_, err = mig.Sql.Exec(InsertNewEntry(mig.dialect), s.state)
	internal.AddError(err)
}

// Gets a list of all migration files in dir
func (mig migConfig) getAllMigrations() []string {
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
	allMigs := mig.getAllMigrations()
	ranMigs := mig.getRanMigrations()
	unranMigs := make([]string, 0)

	for i := 0; i < len(allMigs); i++ {
		appendItem := true
		for x := 0; x < len(ranMigs); x++ {
			if ranMigs[x].id == allMigs[i] {
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
		log.Println("all migrations have been ran")
		return
	}
	log.Println("these migrations were not ran")
	log.Println(unranMigs)
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
		log.Println(err)
	}
}

func (mig migConfig) readDir(dir string) ([]os.DirEntry, error) {
	dir = fmt.Sprintf("%s%s", mig.baseDir, dir)
	return os.ReadDir(dir)
}

func (mig migConfig) readFile(migrationFile string, direction string) []byte {
	migrationFileContents, err := os.ReadFile(fmt.Sprintf("%s%s/%s/%s.sql", mig.baseDir, "migrations", migrationFile, direction))
	if err != nil {
		log.Println(err)
	}
	return migrationFileContents
}

func (mig migConfig) ensureTableExists() {
	_, err := mig.Sql.Exec(GetCreateTableByDialect(mig.dialect))
	if err != nil {
		internal.AddError(err)
	}
}
