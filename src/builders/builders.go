package builders

import (
	"os/exec"
)

var (
	Go = build_go
)

func build_go(filePath string, outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, filePath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
