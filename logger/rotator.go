package logger

import (
	"os"
	"sync"
)

// rotateWriter the new writer for rotating the files
type rotateWriter struct {
	sync.Mutex
	filename string // should be set to the actual filename
	fp       *os.File
}

// Write will call the write method fo fp i.e. the actual file
func (w *rotateWriter) Write(output []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	return w.fp.Write(output)
}

// Rotate will check if the writer has a open file if it has it will close it and will create a new file with the filename
func (w *rotateWriter) Rotate() (err error) {
	w.Lock()
	defer w.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return err
		}
	}

	w.fp, err = os.Create(w.filename)
	return err
}
