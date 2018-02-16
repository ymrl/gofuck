package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed reading %s\n", os.Args[1])
		os.Exit(1)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	var code []string
	for scanner.Scan() {
		code = append(code, strings.Split(scanner.Text(), "")...)
	}

	var pointer int = 0
	var memory []int
	for length, index := len(code), 0; index < length; index++ {
		if pointer >= len(memory) {
			memory = append(memory, make([]int, pointer-len(memory)+1)...)
		}
		op := code[index]
		switch op {
		case ">":
			pointer++
		case "<":
			pointer--
			if pointer < 0 {
				fmt.Fprint(os.Stderr, "Error: Pointer goes negative\n")
				os.Exit(2)
			}
		case "+":
			memory[pointer] = (memory[pointer] + 1) & 0xff
		case "-":
			memory[pointer] = (memory[pointer] - 1) & 0xff
		case "[":
			if memory[pointer] == 0 {
				count := 0
				index++
				for ; ; index++ {
					if index >= length {
						fmt.Fprint(os.Stderr, "Error: ] Not found\n")
						os.Exit(2)
					} else if code[index] == "]" {
						if count == 0 {
							break
						} else {
							count--
						}
					} else if code[index] == "[" {
						count++
					}
				}
			}
		case "]":
			if memory[pointer] != 0 {
				count := 0
				index--
				for ; ; index-- {
					if index < 0 {
						fmt.Fprint(os.Stderr, "Error: [ Not found\n")
						os.Exit(2)
					} else if code[index] == "[" {
						if count == 0 {
							break
						} else {
							count--
						}
					} else if code[index] == "]" {
						count++
					}
				}
			}
		case ".":
			fmt.Print(string(memory[pointer]))
		case ",":
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			memory[pointer] = int(input[0]) & 0xff
		}
	}
}
