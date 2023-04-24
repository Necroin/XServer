package builders

import (
	"fmt"
	"os/exec"
)

func Tool(tool string, filePath string, outputPath string, flags ...string) error {
	cmdArguments := append(flags, []string{"-o", outputPath, filePath}...)
	fmt.Println(cmdArguments)
	cmd := exec.Command(tool, cmdArguments...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Go(filePath string, outputPath string, flags ...string) error {
	cmdArguments := append([]string{"build"}, flags...)
	return Tool("go", filePath, outputPath, cmdArguments...)
}

func Cpp(filePath string, outputPath string, flags ...string) error {
	return Tool("go", filePath, outputPath, flags...)
}
