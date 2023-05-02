package main

import (
	"github.com/harshrastogiexe/bookmgmt/common/models"
	"github.com/harshrastogiexe/bookmgmt/db/bookrepo/local"
)

func main() {
	store := local.NewCollection("Javascript")
	store.Create(&models.BookInfo{})

}
