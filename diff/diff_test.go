package diff

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseFiles(t *testing.T) {
	f, err := os.Open("test.diff")
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	files, err := ParseFiles(string(b))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(files), 15; want != got {
		t.Errorf("wrong number of files: got %d, want %d", got, want)
	}
}
