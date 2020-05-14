package tincd

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCreate(t *testing.T) {
	const netName = "alfaomega"

	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
		return
	}
	path := filepath.Join(tmp, netName)

	ntw, err := Create(path, "10.10.0.0/16")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", ntw)
}
