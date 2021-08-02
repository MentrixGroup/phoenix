package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBook  = "{{cite book |editor-last=Heilbron |author1=Fujia Yang |author2=Joseph H. Hamilton |editor-first=John L. |title=The Oxford Companion to the History of Modern Science |url=https://books.google.com/books?id=abqjP-_KfzkC&pg=PA233 |date=2003 |publisher=Oxford University Press |isbn=978-0-19-974376-6 |page=233}}"
	testTitle = "The Oxford Companion to the History of Modern Science"
)

func TestSmatch(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful simple match", func(t *testing.T) {
		title := smatch(testBook, "title")

		assert.Equal(testTitle, title)
	})

	t.Run("no simple match", func(t *testing.T) {
		res := smatch(testBook, "nomatch")

		assert.Equal("", res)
	})
}

func TestSnmatch(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful number match", func(t *testing.T) {
		authors := snmatch(testBook, "author")

		assert.Equal(testAuthor[0], authors[0])
	})

	t.Run("no number match", func(t *testing.T) {
		authors := snmatch(testBook, "nomatch")

		assert.Equal(0, len(authors))
	})
}
