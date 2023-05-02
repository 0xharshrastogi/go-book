package local

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path"

	"github.com/harshrastogiexe/bookmgmt/common/models"
	"github.com/harshrastogiexe/bookmgmt/common/repo"
)

const BASE_FILE_PATH = "temp/collections"

type BookLocalRepo struct {
	name string
}

func (b *BookLocalRepo) Create(book *models.BookInfo) error {
	books, err := b.Books()
	if err != nil {
		return err
	}
	books = append(books, *book)
	file, err := os.OpenFile(b.path(), os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer closeFile(file)
	if err := json.NewEncoder(file).Encode(&books); err != nil {
		return err
	}
	return nil
}

func (b *BookLocalRepo) Books() ([]models.BookInfo, error) {
	file, err := os.Open(b.path())
	if err != nil {
		return nil, err
	}
	defer closeFile(file)
	books := make([]models.BookInfo, 0, 10)
	if err := json.NewDecoder(file).Decode(&books); err != nil {
		if !errors.Is(err, io.EOF) {
			return books, err
		}
	}
	return books, nil
}

func (b *BookLocalRepo) CollectionName() string {
	return b.name
}

func (b *BookLocalRepo) path() string {
	return path.Join(BASE_FILE_PATH, b.CollectionName()+".json")
}

// create a new collection inside 'temp/data/' folder, may
// cause panic when not able to create file/folder
func NewCollection(name string) repo.BookRepository {
	repo := &BookLocalRepo{name: name}

	if isExist(repo.path()) {
		return repo
	}

	if err := os.MkdirAll(BASE_FILE_PATH, os.FileMode(0777)); err != nil {
		panic(err)
	}
	file, err := os.Create(repo.path())
	if err != nil {
		panic(err)
	}
	defer closeFile(file)
	return repo
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println(err)
	}
}
