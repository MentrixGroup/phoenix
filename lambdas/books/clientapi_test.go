package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikimedia/phoenix/common"
)

const testGoogleBook = `{
  "kind": "books#volumes",
  "totalItems": 1,
  "items": [
    {
      "kind": "books#volume",
      "id": "FHdp36BsN9AC",
      "etag": "OoDIpuFKuzY",
      "selfLink": "https://www.googleapis.com/books/v1/volumes/FHdp36BsN9AC",
      "volumeInfo": {
        "title": "Einstein:",
        "subtitle": "The Life and Times",
        "authors": [
          "Ronald W. Clark"
        ],
        "publisher": "Harper Collins",
        "publishedDate": "1984",
        "description": "A compelling biography of the great physicist focuses on the intellect, and the philosophical tensions, that made Einstein such great scientist, and an interesting man. Reissue. PW. NYT.",
        "industryIdentifiers": [
          {
            "type": "ISBN_13",
            "identifier": "9780380011599"
          },
          {
            "type": "ISBN_10",
            "identifier": "038001159X"
          }
        ],
        "readingModes": {
          "text": false,
          "image": false
        },
        "pageCount": 878,
        "printType": "BOOK",
        "categories": [
          "Biography & Autobiography"
        ],
        "averageRating": 4,
        "ratingsCount": 2,
        "maturityRating": "NOT_MATURE",
        "allowAnonLogging": false,
        "contentVersion": "0.1.1.0.preview.0",
        "panelizationSummary": {
          "containsEpubBubbles": false,
          "containsImageBubbles": false
        },
        "imageLinks": {
          "smallThumbnail": "http://books.google.com/books/content?id=FHdp36BsN9AC&printsec=frontcover&img=1&zoom=5&source=gbs_api",
          "thumbnail": "http://books.google.com/books/content?id=FHdp36BsN9AC&printsec=frontcover&img=1&zoom=1&source=gbs_api"
        },
        "language": "un",
        "previewLink": "http://books.google.nl/books?id=FHdp36BsN9AC&dq=isbn:9780380011599&hl=&cd=1&source=gbs_api",
        "infoLink": "http://books.google.nl/books?id=FHdp36BsN9AC&dq=isbn:9780380011599&hl=&source=gbs_api",
        "canonicalVolumeLink": "https://books.google.com/books/about/Einstein.html?hl=&id=FHdp36BsN9AC"
      },
      "saleInfo": {
        "country": "NL",
        "saleability": "NOT_FOR_SALE",
        "isEbook": false
      },
      "accessInfo": {
        "country": "NL",
        "viewability": "NO_PAGES",
        "embeddable": false,
        "publicDomain": false,
        "textToSpeechPermission": "ALLOWED",
        "epub": {
          "isAvailable": false
        },
        "pdf": {
          "isAvailable": false
        },
        "webReaderLink": "http://play.google.com/books/reader?id=FHdp36BsN9AC&hl=&printsec=frontcover&source=gbs_api",
        "accessViewStatus": "NONE",
        "quoteSharingAllowed": false
      },
      "searchInfo": {
        "textSnippet": "A compelling biography of the great physicist focuses on the intellect, and the philosophical tensions, that made Einstein such great scientist, and an interesting man. Reissue. PW. NYT."
      }
    }
  ]
}`

func getGoogleBookHandler(rw http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["q"]

	if !ok || len(keys[0]) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	q := keys[0]

	if q != fmt.Sprintf("isbn:%s", testIsbn) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	_, _ = rw.Write([]byte(testGoogleBook))
}

func badReqTestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

const (
	testIsbn   = "9780380011599"
	testGoogle = "/books/v1/volumes"
)

func createTestServer() http.Handler {
	srv := http.NewServeMux()

	srv.HandleFunc(testGoogle, getGoogleBookHandler)

	return srv
}

func TestClient(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createTestServer())
	defer srv.Close()

	client := NewClient(srv.URL)

	t.Run("get google book", func(t *testing.T) {
		testBook := &common.Book{Isbn: testIsbn}
		book, err := client.GetBook(testIsbn)

		assert.NoError(err)
		assert.Equal(testBook.Isbn, book.Isbn)
	})
}
