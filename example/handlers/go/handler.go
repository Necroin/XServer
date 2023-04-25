package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("[Go Handler] Started")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
