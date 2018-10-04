package imageupload

import (
	"strings"
	"testing"
)

func init() {
	initExtMap()
	fs = dummyFS{}
}

func TestUnknownFormat(t *testing.T) {
	_, err := saveFile(strings.NewReader("nop"), "/", "testID", "unknown", 0)
	if err != ErrFileNotSupported {
		t.Errorf("unexpected error: %v", err)
	}
}

type dummyFile struct{}

func (dummyFile) Close() error { return nil }
func (dummyFile) Write(p []byte) (n int, err error) { return 0, nil}

type dummyFS struct{}

func (dummyFS) Create(name string) (file, error) { return dummyFile{}, nil }