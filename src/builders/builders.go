package builders

import (
	"fmt"
	"os/exec"
)

var (
	Go     = build_go
	Custom = tool_build
)

func tool_build(tool string, filePath string, outputPath string, flags ...string) error {
	cmdArguments := append(flags, []string{"-o", outputPath, filePath}...)
	fmt.Println(cmdArguments)
	cmd := exec.Command(tool, cmdArguments...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func build_go(filePath string, outputPath string, flags ...string) error {
	cmdArguments := append([]string{"build"}, flags...)
	return tool_build("go", filePath, outputPath, cmdArguments...)
}
