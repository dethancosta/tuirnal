package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/dethancosta/tuirnal/internal/models"
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

func getNewJournalName(jModel models.JournalModel, author string, name string) (string, error) {
	j, err := jModel.Get(author, name)
	var NoJournalError *models.NoJournalError
	if errors.As(err, &NoJournalError) {
		err = jModel.Insert(author, name)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	for j != nil {
		fmt.Println("A journal with that name already exists. Please give a different name")
		fmt.Scanln(&name)
		j, err = jModel.Get(author, name)
		if err != nil && !errors.Is(err, &models.NoJournalError{}) {
			return "", err
		}
	}

	return name, nil
}

// TODO refactor to return error rather than panic with log.Fatal
func InitApp() *Application {
	user, err := user.Current()
	if err != nil {
		log.Fatal("Could not establish username")
	}
	userName := user.Username

	dirPath := filepath.Join("/", "Users", userName, ".journi")
	_, err = os.ReadDir(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0777)
		if err != nil {
			log.Fatalf("Could not establish journi files directory: %v", err)
		}
	} else if err != nil {
		log.Fatal(err.Error())
	}
	dbUrl := filepath.Join(dirPath, "journi.db") // TODO update to be portable; this is specific to MacOS
	jDb, err := bbolt.Open(dbUrl, 0666, nil)
	if err != nil {
		log.Fatalf("Could not open journi.db file: %v", err)
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
