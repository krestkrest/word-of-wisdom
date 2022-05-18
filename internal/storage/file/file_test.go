package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krestkrest/word-of-wisdom/internal/storage/file"
)

const (
	quotesPath = "../../../build/quotes"
)

func TestStorage(t *testing.T) {
	s := file.NewStorage(quotesPath)
	require.NoError(t, s.Start())

	assert.Equal(t, "Why so serious?", s.GetQuoteByIndex(8))
	assert.Equal(t, "Let's put a smile on your face", s.GetQuoteByIndex(20))
}
