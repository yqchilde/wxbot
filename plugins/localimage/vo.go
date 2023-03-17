package localimage

import "os"

type DirEntryVo struct {
	folder       os.DirEntry
	parentFolder os.DirEntry
}
