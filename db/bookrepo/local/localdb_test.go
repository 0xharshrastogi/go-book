package local

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/harshrastogiexe/bookmgmt/common/models"
)

func toBookLocalRepo(b any) *BookLocalRepo {
	v, ok := b.(*BookLocalRepo)
	if !ok {
		panic(fmt.Errorf("conversion to *BoolLocalRepo failed"))
	}
	return v
}

func TestNewCollection(t *testing.T) {
	t.Cleanup(func() {
		fmt.Println("removing", filepath.Dir(BASE_FILE_PATH))
		if err := os.RemoveAll(filepath.Dir(BASE_FILE_PATH)); err != nil {
			log.Println("failed to remove test folder:", err.Error())
		}
	})

	filename := "textCollection"
	t.Run("when file not exist", func(t *testing.T) {
		name := filename
		NewCollection(name)
		p := path.Join(BASE_FILE_PATH, name+".json")
		if isExist(p) {
			return
		}
		t.Errorf("function failed to create %s", p)
	})
}

func TestCreate(t *testing.T) {
	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Dir(BASE_FILE_PATH)); err != nil {
			log.Println("failed to remove test folder:", err.Error())
		}
	})

	t.Run("if no records exist", func(t *testing.T) {
		name := "zlib"
		zlib := toBookLocalRepo(NewCollection(name))
		b := &models.BookInfo{
			Code:  "1",
			Name:  "foo",
			Pages: 100,
		}
		err := zlib.Create(b)
		if err != nil {
			t.Fatal(err)
		}
		books, err := zlib.Books()
		if err != nil {
			t.Fatal(err)
		}
		for _, book := range books {
			if book.Code == b.Code {
				return
			}
		}
		t.Errorf("record not found in records with ID=%v", b.Code)
	})
}

func TestBooks(t *testing.T) {
	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Dir(BASE_FILE_PATH)); err != nil {
			log.Println("failed to remove test folder:", err.Error())
		}
	})
	r := toBookLocalRepo(NewCollection(time.Now().String()))
	t.Run("if collection not exist, should return empty slice", func(t *testing.T) {
		books, err := r.Books()
		if err != nil {
			t.Fatal(err)
		}
		if l := len(books); l != 0 {
			t.Errorf("expected length to be zero, got %d", l)
		}
	})

	t.Run("if collection exist, should return same count in file", func(t *testing.T) {
		option := os.O_RDWR | os.O_TRUNC
		file, err := os.OpenFile(r.path(), option, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		defer closeFile(file)
		want := rand.Intn(100)
		b := make([]models.BookInfo, want)
		for i := 0; i < want; i++ {
			b[i] = models.BookInfo{
				Name:  fmt.Sprint("book", i),
				Pages: uint(i) * 7,
				Code:  fmt.Sprint(i),
			}
		}
		if err := json.NewEncoder(file).Encode(&b); err != nil {
			t.Fatal(err)
		}
		get, err := r.Books()
		if err != nil {
			t.Fatal(err)
		}
		if len(get) != want {
			t.Errorf("want length to be %d, got length %d", want, len(get))
		}
	})
}

func BenchmarkCreate(b *testing.B) {
	b.Cleanup(func() {
		if err := os.RemoveAll(filepath.Dir(BASE_FILE_PATH)); err != nil {
			log.Println("failed to remove test folder:", err.Error())
		}
	})
	r := toBookLocalRepo(NewCollection("zlib"))
	b.StartTimer()
	for i := 0; i < 10; i++ {
		err := r.Create(&models.BookInfo{
			Name:  fmt.Sprint("book", i),
			Pages: uint(i) * 7,
			Code:  fmt.Sprint(i),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.Elapsed()
	b.Logf("time elapsed %d", b.Elapsed().Microseconds())
}
