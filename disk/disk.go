package disk

import "os"

const PAGE_SIZE = 4096

type DiskManager struct {
	HeapFile   *os.File
	NextPageId uint64
}

type PageId int

func New(heapFile *os.File) (*DiskManager, error) {
	// ファイルサイズの取得
	heapFileInfo, err := heapFile.Stat()

	if err != nil {
		return nil, err
	}

	heapFileSize := heapFileInfo.Size()
	nextPageId := heapFileSize / PAGE_SIZE

	return &DiskManager{heapFile, uint64(nextPageId)}, nil
}

func Open(fileName string) (*DiskManager, error) {
	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	return New(fp)
}

func (diskManager *DiskManager) AllocatePage() PageId {
	pageId := diskManager.NextPageId
	diskManager.NextPageId += 1

	return PageId(pageId)
}

func (diskManager *DiskManager) WritePageData(pageId PageId, data []byte) error {
	offset := PAGE_SIZE * pageId
	_, err := diskManager.HeapFile.WriteAt(data, int64(offset))
	if err != nil {
		return err
	}

	return nil
}

func (diskManager *DiskManager) ReadPageData(pageId PageId, data []byte) error {
	offset := PAGE_SIZE * pageId
	_, err := diskManager.HeapFile.ReadAt(data, int64(offset))
	if err != nil {
		return err
	}

	return nil
}

// impl DiskManager {
//    pub fn new(heap_file: File) -> io::Result<Self> {
//        let heap_file_size = heap_file.metadata()?.len();
//        let next_page_id = heap_file_size / PAGE_SIZE as u64;
//        Ok(Self {
//            heap_file,
//            next_page_id,
//        })
//    }
//
//     pub fn open(heap_file_path: impl AsRef<Path>) -> io::Result<Self> {
//         let heap_file = OpenOptions::new()
//             .read(true)
//             .write(true)
//             .create(true)
//             .open(heap_file_path)?;
//         Self::new(heap_file)
//     }
//
//     pub fn read_page_data(&mut self, page_id: PageId, data: &mut [u8]) -> io::Result<()> {
//         let offset = PAGE_SIZE as u64 * page_id.to_u64();
//         self.heap_file.seek(SeekFrom::Start(offset))?;
//         self.heap_file.read_exact(data)
//     }
//
//     pub fn write_page_data(&mut self, page_id: PageId, data: &[u8]) -> io::Result<()> {
//         let offset = PAGE_SIZE as u64 * page_id.to_u64();
//         self.heap_file.seek(SeekFrom::Start(offset))?;
//         self.heap_file.write_all(data)
//     }
//
