package disk_test

import (
	"testing"
  "fmt"
	"github.com/tmrekk121/relly/disk"
)


func TestNew(t *testing.T) {
  dm, err := disk.Open("../testdata/testnew.txt")
  if err != nil {
    t.Fatalf("failed test %v", err)
  }
  buf := make([]byte, 1024)
  for {
      n, err := dm.HeapFile.Read(buf)
      if n == 0{
          break
      }
      if err != nil{
          break
      }
  }
  if string(buf) != "Test New" {
    fmt.Println(string(buf))
    t.Fatalf("missing text: %v", string(buf))
  }
}
