package wal

import (
	"fmt"
	"os"
	"strings"

	"github.com/MitrickX/simple-kv/internal/config"
	utilsOs "github.com/MitrickX/simple-kv/internal/utils/os"
	utilsTime "github.com/MitrickX/simple-kv/internal/utils/time"
)

const (
	initialBufSize            = 4096
	fileNameFromNowTimeLayout = "20060102-150405-000"
)

type WAL struct {
	config          config.ConfigWAL
	buf             []byte
	batchSize       int          // текущее кол-во запросов в батче, записи между батчами разделяются \n
	file            utilsOs.File // текущий файл wal-сегмента
	currentFileSize int          // текущий размер в байтах wal-сегмента
	os              utilsOs.OS
	t               utilsTime.Time
}

func NewWAL(config config.ConfigWAL, os utilsOs.OS, t utilsTime.Time) *WAL {
	return &WAL{
		buf:    make([]byte, 0, initialBufSize),
		config: config,
		os:     os,
		t:      t,
	}
}

func (w *WAL) Write(query string) error {
	if w.batchSize > 0 {
		w.buf = append(w.buf, '\n')
	}
	w.buf = append(w.buf, query...)
	w.batchSize++

	if w.batchSize >= w.config.FlushingBatchSize {
		return w.Flush()
	}

	return nil
}

func (w *WAL) Flush() error {
	err := w.openWalSegmentFile()
	if err != nil {
		return err
	}

	n, err := w.file.Write(w.buf)
	if err != nil {
		return fmt.Errorf("fail to write to wal segment file (%s): %w", w.file.Name(), err)
	}

	if n < len(w.buf) {
		return fmt.Errorf("fail to write all buffered data to file (%s): %w", w.file.Name(), err)
	}

	w.currentFileSize += n

	err = w.file.Sync()
	if err != nil {
		return fmt.Errorf("fail to sync wal segment file (%s): %w", w.file.Name(), err)
	}

	w.buf = w.buf[:0]
	w.batchSize = 0

	if w.currentFileSize >= int(w.config.MaxSegmentSize) {
		err = w.file.Close()
		if err != nil {
			return fmt.Errorf("fail to close wal segment file (%s): %w", w.file.Name(), err)
		}
		w.file = nil
		w.currentFileSize = 0
	}

	return nil
}

func (w *WAL) openWalSegmentFile() error {
	if w.file != nil {
		return nil
	}

	ts := w.t.Now().Format(fileNameFromNowTimeLayout)

	fileName := ts
	filePath := strings.TrimSuffix(w.config.DataDirectory, string(os.PathSeparator)) + string(os.PathSeparator) + fileName
	file, err := w.os.OpenFile(filePath, os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("fail to open new wal segemnt file (%s): %w", fileName, err)
	}
	w.file = file
	w.currentFileSize = 0

	return nil
}
