package scribedb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func BackupDatabase(backupPath string, dbFile string) (int64, error) {

	ts := time.Now().String()[:22]
	dbBakFile := filepath.Join(backupPath, fmt.Sprintf("scribeNB.db-%s", ts))

	src, err := os.Open(dbFile)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(dbBakFile)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	// Copy source to destination
	bytesCopied, err := io.Copy(dst, src)
	if err != nil {
		return 0, err
	}
	return bytesCopied, nil
}
