package dbg

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func CreateLog(filename, content string) error {
	logFolder := filepath.Join("logs", fmt.Sprintf("%v", time.Now().Unix()))
	if err := os.MkdirAll(logFolder, 0777); err != nil {
		return err
	}

	err := os.WriteFile(filepath.Join(logFolder, fmt.Sprintf("F%vDT%v", filename, time.Now().Format(time.RFC3339))), []byte(content), 0777)
	if err != nil {
		return err
	}

	return nil
}
