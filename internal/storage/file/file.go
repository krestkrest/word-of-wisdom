package file

import (
	"bufio"
	"math/rand"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/krestkrest/word-of-wisdom/internal/storage"
)

type Storage struct {
	fileName string

	quotes []string
}

var _ storage.Storage = (*Storage)(nil)

func NewStorage(fileName string) *Storage {
	return &Storage{fileName: fileName}
}

func (s *Storage) Start() error {
	f, err := os.Open(s.fileName)
	if err != nil {
		return errors.Wrap(err, "failed to open file with quotes")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		quote := strings.TrimSpace(scanner.Text())
		if len(quote) != 0 {
			s.quotes = append(s.quotes, quote)
		}
	}
	if err = scanner.Err(); err != nil {
		return errors.Wrap(err, "failed to scan file with quotes")
	}
	if len(s.quotes) == 0 {
		return errors.New("quotes list is empty!")
	}
	return nil
}

func (s *Storage) GetQuote() string {
	return s.GetQuoteByIndex(rand.Int())
}

func (s *Storage) GetQuoteByIndex(index int) string {
	return s.quotes[index%len(s.quotes)]
}
