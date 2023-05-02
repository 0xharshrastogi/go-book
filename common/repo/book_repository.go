package repo

import "github.com/harshrastogiexe/bookmgmt/common/models"

type BookRepository interface {
	// Create saves the book, if success returns error as nil else return error
	Create(b *models.BookInfo) error

	// Books returns all the books stored
	Books() ([]models.BookInfo, error)

	// CollectionName returns the collection name
	CollectionName() string
}
