package buffer

import (
	"fmt"

	"github.com/getlantern/deepcopy"
	"github.com/tmrekk121/relly/disk"
)

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
	Buffers      []Frame
	NextVictimId BufferId
}

type BufferPoolManager struct {
	Disk      disk.DiskManager
	Pool      BufferPool
	PageTable map[disk.PageId]BufferId
}

func (bufferPool *BufferPool) Size() int {
	return len(bufferPool.Buffers)
}

func (bufferPool *BufferPool) incrementId(bufferId BufferId) BufferId {
	return BufferId(int(bufferId+1) % bufferPool.Size())
}

func (bufferPool *BufferPool) Evict() (BufferId, error) {
	poolSize := bufferPool.Size()
	consecutivePinned := 0

	for {
		nextVictimId := bufferPool.NextVictimId
		frame := bufferPool.Buffers[nextVictimId]

		if frame.UsageCount == 0 {
			return nextVictimId, nil
		}

		if frame.refCount == 0 {
			frame.UsageCount -= 1
			consecutivePinned = 0
		} else {
			consecutivePinned += 1
			if consecutivePinned >= poolSize {
				return -1, fmt.Errorf("no bufferId")
			}
		}

		bufferPool.NextVictimId = bufferPool.incrementId(bufferPool.NextVictimId)
	}
}

func (BufferPoolManager *BufferPoolManager) FetchPage(pageId disk.PageId) (Buffer, error) {
	bufferId, ok := BufferPoolManager.PageTable[pageId]

	// ページがバッファプールにある場合
	if ok {
		BufferPoolManager.Pool.Buffers[bufferId].refCount += 1
		frame := BufferPoolManager.Pool.Buffers[bufferId]

		buffer := Buffer{
			PageId:  frame.Buffer.PageId,
			Page:    frame.Buffer.Page,
			IsDirty: frame.Buffer.IsDirty,
		}

		frame.UsageCount += 1
		frame.refCount -= 1

		return buffer, nil
	}

	// ページがバッファプールにない場合
	bufferId, err := BufferPoolManager.Pool.Evict()
	if err != nil {
		return Buffer{}, err
	}
	frame := BufferPoolManager.Pool.Buffers[bufferId]
	evictPageId := frame.Buffer.PageId
	buffer := frame.Buffer
	if buffer.IsDirty {
		err = BufferPoolManager.Disk.WritePageData(evictPageId, buffer.Page[:])
		if err != nil {
			return Buffer{}, err
		}
	}
	buffer.PageId = pageId
	buffer.IsDirty = false
	err = BufferPoolManager.Disk.ReadPageData(buffer.PageId, buffer.Page[:])
	if err != nil {
		return Buffer{}, err
	}
	frame.UsageCount = 1

	page := Buffer{}
	deepcopy.Copy(page, buffer)
	delete(BufferPoolManager.PageTable, evictPageId)
	BufferPoolManager.PageTable[pageId] = bufferId
	return page, nil
}
