package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type Books struct {
	id       int
	bookName string
	Authors
	stockCode     string
	isbn          string
	pageNumber    int
	price         float64
	stockQuantity int
	isDeleted     bool
}

type Authors struct {
	authorName string
}

var books []Books

var bookNames = []string{
	"Don Quixote",
	"One Hundred Years of Solitude",
	"The Great Gatsby",
	"Moby Dick",
	"War and Peace",
	"Hamlet",
}

var authorNames = []string{
	"Miguel de Cervantes",
	"Gabriel García Márquez",
	"F. Scott Fitzgerald",
	"Herman Melville",
	"Leo Tolstoy",
	"William Shakespeare",
}

// Kitap nesnelerimizi oluşturup books sliceına ekliyor
func init() {
	for index, book := range bookNames {
		bookObject := Books{
			Authors: Authors{
				authorName: authorNames[index],
			},
			id:            index + 1,
			bookName:      book,
			pageNumber:    createRandomNumber(500) + 100,
			price:         float64(createRandomNumber(60) + 10),
			stockQuantity: createRandomNumber(20),
			isbn:          createIsbnCode(),
			isDeleted:     false,
			stockCode:     createIsbnCode()[:8],
		}

		books = append(books, bookObject)
	}
}

func main() {

	input := os.Args

	if len(input) >= 3 {
		if input[1] == "search" {
			searchKey := strings.Join(input[2:], " ")
			founded := Search(strings.ToLower(searchKey))

			if len(founded) != 0 {
				for _, i := range founded {
					fmt.Printf("Kitap adi: %v\nYazar: %v\nSayfa Sayısı: %d\nStok adedi: %d\nISBN: %v\nStok Kodu: %v\n\n------\n\n",
						i.bookName,
						i.authorName,
						i.pageNumber,
						i.stockQuantity,
						i.isbn,
						i.stockCode)
				}
			} else {
				fmt.Println("Kitap bulunamadı.")
			}
		} else if input[1] == "delete" {

			id, err := strconv.Atoi(input[2])

			if err == nil && len(input) == 3 {
				book, err := findBook(id)
				if err == nil {
					delErr := deleteInterface(&book)
					if delErr == nil {
						fmt.Println("Kitap silindi!")
					} else {
						fmt.Println(delErr)
					}
				} else {

					fmt.Println(err)
				}
			} else {
				fmt.Println("Yanlış input verildi 'delete id' şeklinde giriş yapınız.")
			}

		} else if input[1] == "buy" {
			id, errId := strconv.Atoi(input[2])
			if len(input) == 4 {
				count, errCount := strconv.Atoi(input[3])

				if errId == nil && errCount == nil {

					book, err := findBook(id)
					if !book.isDeleted {
						if err == nil {
							buyErr := book.Buy(count)

							if buyErr != nil {
								fmt.Println(buyErr)
							} else {
								fmt.Printf("%d adet kitap alındı\n", count)
							}
						} else {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Kitap silinmiş")
					}

				} else {
					fmt.Println("Yanlış input verildi.")
				}
			} else {
				fmt.Println("Yanlış input verildi.")
			}

		}

	} else {
		fmt.Println("Yanlış input verildi.")
	}
}

// Random bir numara üretip int şeklinde dönderiyor.
func createRandomNumber(b int64) int {
	var a = rand.Reader

	n, _ := rand.Int(a, big.NewInt(b))

	return int(n.Int64())

}

// Random ISBN numarası üretip string olarak döndürüyor
func createIsbnCode() string {
	var ısbn string

	for i := 0; i < 14; i++ {
		ısbn += strconv.Itoa(createRandomNumber(10))
	}

	return ısbn
}

// aldığı stringi kitap nesneleerinin adında yazar adında ve stok kodunda ve isbn kodunda arar eşleşen ve silinmemiş kitapları slice olarak döndürür
func Search(s string) []Books {
	var founded []Books

	for _, book := range books {
		if (!book.isDeleted) && (strings.Contains(strings.ToLower(book.bookName), s) || strings.Contains(strings.ToLower(book.authorName), s) || strings.Contains(strings.ToLower(book.stockCode), s) || strings.Contains(strings.ToLower(book.isbn), s)) {
			founded = append(founded, book)
		}
	}
	return founded
}

// daha önceden silinmemiş bir kitabı siler, öncede silimişse hata verir
func (b *Books) Delete() error {
	if !b.isDeleted {
		b.isDeleted = true
		return nil
	} else {
		delErr := errors.New("kitap önceden silinmiş")
		return delErr
	}
}

// silinmemiş bir kitabı satın alır kitap silindiyse uyarı verir, mevcut stokdan fazla sipariş verildiyse uyarı verilir
func (b *Books) Buy(n int) error {
	buyErr := errors.New("")
	if b.isDeleted {
		buyErr = fmt.Errorf("kitap silindiği için satın alınamıyor")
		return buyErr
	}

	if b.stockQuantity >= n {
		b.stockQuantity -= n
		return nil
	} else {
		buyErr = fmt.Errorf("%wLütfen %d adet veya daha az sipariş verin", buyErr, b.stockQuantity)
		return buyErr
	}
}

type Deletable interface {
	Delete() error
}

// delete Deletable interfaceini çağırır
func deleteInterface(d Deletable) error {
	err := d.Delete()
	return err
}

// girilen id numarına sahip kitabı döndürür, belirtilen id ye sahip kkitap yoksa hata verir
func findBook(n int) (Books, error) {
	for _, book := range books {
		if book.id == n {
			return book, nil

		}
	}

	buyErr := errors.New("kitap mevcut değil")

	return Books{}, buyErr
}
