package input

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func ClearCache() error {
	return os.RemoveAll(cachedir())
}

func daydata(n int) ([]byte, error) {
	if !IgnoreCache {
		b, err := ioutil.ReadFile(dayinputfn(n))
		if err == nil {
			return b, nil
		}
	}

	if err := downloaddaydata(n); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(dayinputfn(n))
}

func cachedir() string {
	d, err := os.UserCacheDir()
	if err == nil {
		return filepath.Join(d, cachesubdir())
	}
	return ".input"
}

func dayinputfn(day int) string {
	return filepath.Join(cachedir(), fmt.Sprint(day))
}

var dataclient struct {
	once sync.Once

	c   *http.Client
	err error
}

func createFile(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

////////////////////////////////////////////////////////////////////////////////

func DayBlocks(day int) *BlockScanner {
	return NewBlockScanner(Reader(day))
}

// BlockScanner returns multi-line strings in the input
// separated by an empty line.
type BlockScanner struct {
	scanner *bufio.Scanner
	token   []byte
}

func NewBlockScanner(r io.Reader) *BlockScanner {
	return &BlockScanner{scanner: bufio.NewScanner(r)}
}

func (b *BlockScanner) Scan() bool {
	b.token = b.token[:0]
	for b.scanner.Scan() {
		p := b.scanner.Bytes()
		if len(p) == 0 {
			return true
		}
		if len(b.token) != 0 {
			b.token = append(b.token, '\n')
		}
		b.token = append(b.token, p...)
	}
	return len(b.token) != 0
}

func (b *BlockScanner) Bytes() []byte { return b.token }
func (b *BlockScanner) Text() string  { return string(b.token) }
func (b *BlockScanner) Err() error    { return b.scanner.Err() }
