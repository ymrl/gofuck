package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Context struct {
	pointer int
	index   int
	memory  []int
}

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

	c := Context{}
	for length := len(code); c.index < length; {
		c = step(code, c)
	}
}

func step(code []string, c Context) Context {
	if c.pointer >= len(c.memory) {
		c.memory = append(c.memory, make([]int, c.pointer-len(c.memory)+1)...)
	}
	length := len(code)
	op := code[c.index]
	switch op {
	case ">":
		c.pointer++
	case "<":
		c.pointer--
		if c.pointer < 0 {
			fmt.Fprint(os.Stderr, "Error: Pointer goes negative\n")
			os.Exit(2)
		}
	case "+":
		c.memory[c.pointer] = (c.memory[c.pointer] + 1) & 0xff
	case "-":
		c.memory[c.pointer] = (c.memory[c.pointer] - 1) & 0xff
	case "[":
		if c.memory[c.pointer] == 0 {
			count := 0
			c.index++
			for ; ; c.index++ {
				if c.index >= length {
					fmt.Fprint(os.Stderr, "Error: ] Not found\n")
					os.Exit(2)
				} else if code[c.index] == "]" {
					if count == 0 {
						break
					} else {
						count--
					}
				} else if code[c.index] == "[" {
					count++
				}
			}
		}
	case "]":
		if c.memory[c.pointer] != 0 {
			count := 0
			c.index--
			for ; ; c.index-- {
				if c.index < 0 {
					fmt.Fprint(os.Stderr, "Error: [ Not found\n")
					os.Exit(2)
				} else if code[c.index] == "[" {
					if count == 0 {
						break
					} else {
						count--
					}
				} else if code[c.index] == "]" {
					count++
				}
			}
		}
	case ".":
		fmt.Print(string(c.memory[c.pointer]))
	case ",":
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		c.memory[c.pointer] = int(input[0]) & 0xff
	}
	c.index++
	return c
}
