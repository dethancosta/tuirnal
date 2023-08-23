package helpers

import (
	"fmt"
	"log"
	"os"
	//"os/user"
	//"path/filepath"
	"time"

	"github.com/dethancosta/tuirnal/internal/models"
	gap "github.com/muesli/go-app-paths"
	"go.etcd.io/bbolt"
)

type Application struct {
	Db           *bbolt.DB
	AuthorModel  models.AuthorModel
	JournalModel models.JournalModel
	EntryModel   models.EntryModel
}

func timeString(t time.Time) string {
	return fmt.Sprintf("%s %d, %d", t.Month().String(), t.Day(), t.Year())
}

// TODO refactor to return error rather than panic with log.Fatal
func InitApp(dbFilename string) *Application {
	scope := gap.NewScope(gap.User, "tuirnal")
	dirPath, err := scope.DataPath("")
	_, err = os.ReadDir(dirPath)

	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0777)
		if err != nil {
			log.Fatalf("Could not establish tuirnal files directory: %v", err)
		}
	} else if err != nil {
		log.Fatal(err.Error())
	}

	dbUrl, err := scope.DataPath(dbFilename)
	if err != nil {
		log.Fatal(err.Error())
	}

	jDb, err := bbolt.Open(dbUrl, 0666, nil)
	if err != nil {
		log.Fatalf("Could not open "+dbFilename+" file: %v", err)
	}

	err = jDb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("journals"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("entries"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("authors"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	app := Application{
		Db:           jDb,
		AuthorModel:  models.AuthorModel{DB: jDb},
		JournalModel: models.JournalModel{DB: jDb},
		EntryModel:   models.EntryModel{DB: jDb},
	}
	return &app
}
