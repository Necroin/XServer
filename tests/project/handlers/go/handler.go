package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("[Go Hanler] Started")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
