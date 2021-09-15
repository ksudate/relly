package buffer

import "github.com/tmrekk121/relly/disk"

type Page [disk.PAGE_SIZE]byte
type BufferId int

type Buffer struct {
  PageId  disk.PageId
  Page    Page
  IsDirty bool
}

type Frame struct {
  UsageCount int
  Buffer     *Buffer
  refCount   int
}

type BufferPool struct {
  Buffers []Frame
  NextVictimId BufferId
}

type BufferPoolManager struct {
  Disk disk.DiskManager
  Pool BufferPool
  PageTable map[disk.PageId]BufferId
}
