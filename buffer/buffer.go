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

func (bufferPool *BufferPool) Evict() *BufferId {
	poolSize := bufferPool.Size()
	consecutivePinned := 0

	for {
		nextVictimId := bufferPool.NextVictimId
		frame := bufferPool.Buffers[nextVictimId]

		if frame.UsageCount == 0 {
			return &nextVictimId
		}

		if frame.refCount == 0 {
			frame.UsageCount -= 1
			consecutivePinned = 0
		} else {
			consecutivePinned += 1
			if consecutivePinned >= poolSize {
				return nil
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
}
