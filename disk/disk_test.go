package disk_test

import (
	"testing"

	"github.com/tmrekk121/relly/disk"
)

func TestNew(t *testing.T) {
	dm, err := disk.Open("../testdata/testnew.txt")
	if err != nil {
		t.Fatalf("failed test %v", err)
	}
	buf := make([]byte, 1024)
	var str string
	for {
		n, err := dm.HeapFile.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
		str = string(buf[:n])
	}
	if str != "Test New\n" {
		t.Fatalf("missing text: %+v", str)
	}
}
