// This file contains code to convert text files to TeX.
// See README for the format of the text files.
//
// Compile:
//   cd ancient-chinese/go
//   go install tex
// Run
//   bin/tex txt/shiji-simplified.txt
// It generates a new file txt/shiji-simplified.tex
//
// You probably need to change --font-name if you don't have SimSun in your system.
// Run fc-list to find installed fonts. Refer to fc-cache about how to install new font.
//
// Use xelatex to convert TeX to PDF.
// You need at least the following packages to run xelatex:
//     sudo apt-get install texlive-xetex texlive-lang-cjk cjk-latex
// Suggest to use https://www.tug.org/texlive/acquire-netinstall.html.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var fontName = flag.String("font-name", "SimSun", "The font name.")
var fontSize = flag.Int("font-size", 16, "The font size. This default setting is for 9inch kindle.")
var titleFontName = flag.String("title-font-name", "KaiTi", "The title font name.")

func GetTitlePage(title, author string) string {
	return fmt.Sprintf(
		`\begin{titlepage}
\begin{center}
\vspace*{\fill}
\CJKfamily{title}
\textsc{\textbf{\huge %s}}\\[0.5cm]
\CJKfamily{}
\textsc{\large %s}\\[1.5cm]
\textsc{\small \url{https://code.google.com/p/ancient-chinese}}\\
{\today}
\vspace*{\fill}
\end{center}
\end{titlepage}`, title, author)
}

func GetChapterStart(title string) string {
	return fmt.Sprintf(
		`\cleardoublepage
\phantomsection
\chapter{%s}`, title, title)
}

func ConvertToTex(input, output string) {
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

	// Ouput headers.
	fmt.Fprintf(outputFile, "\\documentclass[fontsize=%dpt]{scrbook}\n", *fontSize)
	fmt.Fprintln(outputFile, `\usepackage{hyperref}`)
	fmt.Fprintln(outputFile, `\usepackage{indentfirst}`)
	fmt.Fprintln(outputFile, `\usepackage{xeCJK}`)
	fmt.Fprintln(outputFile, `\CJKspace`)
	fmt.Fprintf(outputFile, "\\setCJKmainfont{%s}\n", *fontName)
	fmt.Fprintf(outputFile, "\\setCJKfamilyfont{title}{%s}\n", *titleFontName)
	fmt.Fprintln(outputFile, `\XeTeXlinebreaklocale "zh"`)
	fmt.Fprintln(outputFile, `\XeTeXlinebreakskip 0pt plus 1pt`)
	fmt.Fprintln(outputFile, `\setcounter{secnumdepth}{-1}`)
	fmt.Fprintln(outputFile, `\setcounter{tocdepth}{0}`)
	fmt.Fprintln(outputFile, `\linespread{1.2}`)
	fmt.Fprintln(outputFile, `\setlength{\parindent}{3em}`)
	fmt.Fprintln(outputFile, `\sloppy`)
	fmt.Fprintln(outputFile, `\begin{document}`)

	var title string
	var author string
	var chapters []string
	currentChapter := 0
	for inputScanner.Scan() {
		line := strings.TrimSpace(inputScanner.Text())
		if len(line) == 0 {
			continue
		} else if len(title) == 0 {
			title = line
			log.Printf("Title: %s\n", title)
		} else if len(author) == 0 {
			author = line
			log.Printf("Author: %s\n", author)
			fmt.Fprintln(outputFile, GetTitlePage(title, author))
		} else if len(chapters) == 0 {
			for {
				chapters = append(chapters, line)
				if !inputScanner.Scan() {
					break
				}
				line = strings.TrimSpace(inputScanner.Text())
				if len(line) == 0 {
					break
				}
			}
			log.Printf("%d Chapters.\n", len(chapters))
			fmt.Fprintln(outputFile, `\tableofcontents{}`)
		} else {
			if currentChapter < len(chapters) && strings.HasSuffix(line, chapters[currentChapter]) {
				log.Printf("Chapter %d: %s\n", currentChapter, chapters[currentChapter])
				if len(line) > len(chapters[currentChapter]) {
					fmt.Fprintln(outputFile, "\\par\n"+line[0:len(line)-len(chapters[currentChapter])])
				}
				fmt.Fprintln(outputFile, GetChapterStart(chapters[currentChapter]))
				currentChapter++
			} else {
				fmt.Fprintln(outputFile, "\\par\n"+line)
			}
		}

	}
	fmt.Fprintln(outputFile, `"\end{document}`)
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
		output := input[0:len(input)-4] + ".tex"
		log.Printf("Converting %s to %s ...\n", input, output)
		ConvertToTex(input, output)
	}
}
