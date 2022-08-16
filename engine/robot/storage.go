package robot

import "os"

type JsonLocalStorage struct {
	FileName string
	file     *os.File
}

func (j *JsonLocalStorage) Read(p []byte) (n int, err error) {
	if j.file == nil {
		j.file, err = os.Open(j.FileName)
		if err != nil {
			return 0, err
		}
	}
	return j.file.Read(p)
}

func (j *JsonLocalStorage) Write(p []byte) (n int, err error) {
	j.file, err = os.Create(j.FileName)
	if err != nil {
		return 0, err
	}
	return j.file.Write(p)
}
