package disk_test

import (
	"fmt"
	"testing"

	"github.com/tmrekk121/relly/disk"
)

func TestOpen(t *testing.T) {
	dm, err := disk.Open("../testdata/testopen.txt")
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
	if str != "Test Open\n" {
		t.Fatalf("missing text: %+v", str)
	}

	nextPageId := dm.NextPageId

	if nextPageId != 0 {
		t.Fatalf("invalid NextPageId: %+v", dm)
	}
}

func TestAllocatePage(t *testing.T) {
	dm, err := disk.Open("../testdata/testopen.txt")
	if err != nil {
		t.Fatalf("failed test %v", err)
	}

	pageId := dm.AllocatePage()

	if pageId != 0 {
		t.Fatalf("invalid pageId: %+v", pageId)
	}

	if dm.NextPageId != 1 {
		t.Fatalf("invalid NextPageId: %+v", dm)
	}
}

// func TestWritePageData(t *testing.T) {
// 	dm, err := disk.Open("../testdata/write.txt")
// 	if err != nil {
// 		t.Fatalf("failed test %v", err)
// 	}

// 	dm.WritePageData(disk.PageId(dm.NextPageId), []byte("ハローワールド"))
// 	dm.AllocatePage()
// }

func TestReadPageData(t *testing.T) {
	dm, err := disk.Open("../testdata/write.txt")
	if err != nil {
		t.Fatalf("failed test %v", err)
	}

	data := make([]byte, 4096)
	dm.ReadPageData(disk.PageId(0), data)

	fmt.Println(string(data))
}
