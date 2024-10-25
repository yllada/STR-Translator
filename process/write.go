package process

import (
	"fmt"
	"os"
)

func writeSRTFile(filename, content string) error {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
