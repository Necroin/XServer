package builders

import (
	"os/exec"
)

var (
	Go = build_go
)

func build_go(path string, outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
