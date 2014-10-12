// This file contains code to format text files. It does
// 1. Remove empty lines.
// 2. Replace unicode spaces with normal spaces.
// 3. Merge consecutive spaces into one.
// 4. Fix unpaired quote.
// 5. Merge broken paragraphs.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const startQuote = '“'
const endQuote = '”'

func Format(input, output string) {
	// Open input.
	inputFile, err := os.Open(input)
	if err != nil {
		log.Printf("Failed to open %s for read: %s.", input, err)
		return
	}
	defer inputFile.Close()
	inputScanner := bufio.NewScanner(inputFile)

	// Open output.
	outputFile, err := os.Create(output)
	if err != nil {
		log.Printf("Failed to open %s for write: %s.", output, err)
	}
	defer outputFile.Close()

	var buffer bytes.Buffer
	inQuote := false
	inTable := false
	isSpace := false
	couldEnd := false
	lineNumber := 0
	var title, author string
	for inputScanner.Scan() {
		line := strings.TrimSpace(inputScanner.Text())
		lineNumber++
		if len(line) == 0 {
			continue
		}
		if len(title) == 0 {
			fmt.Fprintln(outputFile, line)
			title = line
			continue
		}
		if len(author) == 0 {
			fmt.Fprintln(outputFile, line)
			author = line
			continue
		}
		if line == "---" {
			inTable = !inTable
			if buffer.Len() > 0 || inQuote {
				if buffer.Len() > 0 {
					fmt.Fprintln(outputFile, buffer.String())
				}
				log.Fatalf("Error @%d: %s\n", lineNumber, line)
			}
			fmt.Fprintln(outputFile, line)
			continue
		}
		if inTable || strings.HasPrefix(line, "+") {
			if buffer.Len() > 0 || inQuote {
				if buffer.Len() > 0 {
					fmt.Fprintln(outputFile, buffer.String())
				}
				log.Fatalf("Error @%d: %s\n", lineNumber, line)
			}
			fmt.Fprintln(outputFile, line)
			continue
		}
		for _, runeValue := range line {
			if unicode.IsSpace(runeValue) {
				if isSpace {
					couldEnd = false
					continue
				} else {
					runeValue = ' '
					isSpace = true
				}
			}
			if runeValue == startQuote || runeValue == endQuote {
				if inQuote {
					runeValue = endQuote
				} else {
					runeValue = startQuote
				}
				inQuote = !inQuote
			}
			encodedBytes := make([]byte, utf8.RuneLen(runeValue))
			utf8.EncodeRune(encodedBytes, runeValue)
			buffer.Write(encodedBytes)
			switch runeValue {
			case '。', '”', '？', '！', '）', '》':
				couldEnd = true
			default:
				couldEnd = false
			}

		}
		if couldEnd && !inQuote {
			fmt.Fprintln(outputFile, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 || inQuote {
		if buffer.Len() > 0 {
			fmt.Fprintln(outputFile, buffer.String())
		}
		log.Fatalln("Error at end of file.")
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Printf("Usage: %s 1.txt [2.txt .. ]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	for _, input := range flag.Args() {
		if !strings.HasSuffix(input, ".txt") {
			log.Printf("Don't know how convert %s. Ignore it.", input)
			continue
		}
		output := input[0:len(input)-4] + ".new.txt"
		log.Printf("Format %s to %s ...\n", input, output)
		Format(input, output)
	}
}
