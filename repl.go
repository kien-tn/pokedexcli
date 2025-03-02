package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		outout := strings.Fields(strings.ToLower(text))
		fmt.Printf("Your command was: %s\n", outout[0])
	}
}
