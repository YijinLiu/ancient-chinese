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
// You probably need to download the free super-big font HanaMin(花園明朝)
// from http://www.zdic.net/appendix/f18.htm
// Please refer to fc-cache about how to install new font.
// Run fc-list to find installed fonts.
//
// Use xelatex to convert TeX to PDF.
// You need at least the following packages to run xelatex:
//     sudo apt-get install texlive-xetex texlive-lang-cjk cjk-latex
// Suggest to use https://www.tug.org/texlive/acquire-netinstall.html.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var fontName = flag.String("font-name", "HanaMinA", "The font name.")
var fallbackFontName = flag.String("fallback-font-name", "HanaMinB", "The fallback font name, for rare words not covered by the default font.")
var titleFontName = flag.String("title-font-name", "KaiTi", "The title font name.")
var fontSize = flag.Int("font-size", 16, "The font size. This default setting is for 9inch kindle.")

func GetLongtableDef() string {
	return `\usepackage{longtable,tabulary}
\makeatletter
\def\ltabulary{%
\def\endfirsthead{\\}%
\def\endhead{\\}%
\def\endfoot{\\}%
\def\endlastfoot{\\}%
\def\tabulary{%
  \def\TY@final{%
\def\endfirsthead{\LT@end@hd@ft\LT@firsthead}%
\def\endhead{\LT@end@hd@ft\LT@head}%
\def\endfoot{\LT@end@hd@ft\LT@foot}%
\def\endlastfoot{\LT@end@hd@ft\LT@lastfoot}%
\longtable}%
  \let\endTY@final\endlongtable
  \TY@tabular}%
\dimen@\columnwidth
\advance\dimen@-\LTleft
\advance\dimen@-\LTright
\tabulary\dimen@}
\def\endltabulary{\endtabulary}
\makeatother`
}

func GetTitlePage(title, author string) string {
	return fmt.Sprintf(
		`\begin{titlepage}
\begin{center}
\vspace*{\fill}
\emph{\textbf{\Huge{\kaiti %s}}}\\[0.5cm]
{\normalsize %s}\\[1.5cm]
{\small\url{https://code.google.com/p/ancient-chinese}}\\
{\small\today}
\vspace*{\fill}
\end{center}
\end{titlepage}`, title, author)
}

var kSectionNames = [...]string{
	"part",
	"chapter",
	"section",
	"subsection",
	"subsubsection",
	"subsubsubsection",
	"paragraph",
	"subparagraph",
}

// "title" is like
// 1) "+XXX" is part.
// 2) "++XXX" is chapter.
// 3) "+++XXX" is section.
// 4) "++++XXX" is subsection.
// 5) "+++++XXX" is subsubsection.
// 6) "++++++XXX" is subsubsubsection.
// 7) "+++++++XXX" is paragraph.
// 8) "++++++++XXX" is subparagraph.
// Returns
//   "start" the tex script to start the section.
//   "sectionType" the section type, 0..7.
//   "outTitle" the section title. (XXX)
func ParseTitleLine(title string) (sectionType int, start, outTitle string) {
	numOfPlus := 0
	for numOfPlus < len(title) && title[numOfPlus] == '+' {
		numOfPlus++
	}
	if numOfPlus < 1 || numOfPlus > len(kSectionNames) {
		log.Fatalf("Unknown title: %s.", title)
	}
	outTitle = title[numOfPlus:]
	if numOfPlus == 2 {
		start = `\cleardoublepage`
	} else {
		start = ""
	}
	sectionType = numOfPlus - 1
	start += `\phantomsection`
	return
}

var kCommentStart = "（"
var kCommentEnd = "）"

func ReplaceCommentWithScript(text string) string {
	var buffer bytes.Buffer
	for {
		start := strings.Index(text, kCommentStart)
		if start == -1 {
			break
		}
		end := strings.Index(text, kCommentEnd)
		if start > 0 {
			buffer.WriteString(text[:start])
		}
		start += len(kCommentStart)
		if start >= end {
			log.Fatalf("Invalid comment: %s\n", text)
		}
		buffer.WriteString(fmt.Sprintf(`{\scriptsize %s}`, text[start:end]))
		text = text[end+len(kCommentEnd):]
		if len(text) == 0 {
			break
		}
	}
	if len(text) > 0 {
		buffer.WriteString(text)
	}
	return buffer.String()
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
	fmt.Fprintln(outputFile, `\KOMAoptions{twoside=false}`)
	fmt.Fprintln(outputFile, `\usepackage{hyperref}`)
	fmt.Fprintln(outputFile, `\usepackage{indentfirst}`)
	fmt.Fprintln(outputFile, GetLongtableDef())
	fmt.Fprintln(outputFile, `\usepackage{xeCJK}`)
	fmt.Fprintln(outputFile, `\xeCJKsetup{AutoFallBack}`)
	fmt.Fprintln(outputFile, `\CJKspace`)
	fmt.Fprintf(outputFile, "\\setCJKmainfont[FallBack=%s]{%s}\n", *fallbackFontName, *fontName)
	fmt.Fprintf(outputFile, "\\newCJKfontfamily[kai]\\kaiti{%s}\n", *titleFontName)
	fmt.Fprintln(outputFile, `\XeTeXlinebreaklocale "zh"`)
	fmt.Fprintln(outputFile, `\XeTeXlinebreakskip 0pt plus 1pt`)
	fmt.Fprintln(outputFile, `\usepackage{fancyhdr}`)
	fmt.Fprintln(outputFile, `\pagestyle{fancy}`)
	fmt.Fprintln(outputFile, `\setcounter{secnumdepth}{-1}`)
	fmt.Fprintln(outputFile, `\setcounter{tocdepth}{2}`)
	fmt.Fprintln(outputFile, `\linespread{1.2}`)
	fmt.Fprintln(outputFile, `\setlength{\parindent}{3em}`)
	fmt.Fprintln(outputFile, `\sloppy`)
	fmt.Fprintln(outputFile, `\begin{document}`)

	var title string
	var author string
	const kTableMarker = "---"
	var sectionTypeToCount [len(kSectionNames)]int
	var sectionTypeToTitle [len(kSectionNames)]string
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
			fmt.Fprintln(outputFile, `\tableofcontents{}`)
			fmt.Fprintln(outputFile, `\newpage`)
		} else if line == kTableMarker {
			var tableRows []string
			for inputScanner.Scan() {
				line = strings.TrimSpace(inputScanner.Text())
				if line == kTableMarker {
					break
				}
				tableRows = append(tableRows, line)
			}
			if len(tableRows) == 0 {
				continue
			}
			fmt.Fprintln(outputFile, `\par`)
			fmt.Fprintln(outputFile, `\begin{scriptsize}`)
			// Allow line break after quote, colon.
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{CJK}{`：,`“}")
			// Allow line break before comma, period.
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{CJK}{`，,`、,`。,`”}")
			// Allow line break before numbers.
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{CJK}{`0,`1,`2,`3,`4,`5,`6,`7,`8,`9}")
			numColumns := -1
			for _, row := range tableRows {
				columns := strings.Split(row, "|")
				if numColumns == -1 {
					numColumns = len(columns)
					fmt.Fprintf(outputFile, "\\begin{ltabulary}{%s|}\n", strings.Repeat("|L", numColumns))
					fmt.Fprintln(outputFile, `\hline`)
				} else if numColumns != len(columns) {
					log.Fatalf("'%s' should have %d columns", row, numColumns)
				}
				fmt.Fprintf(outputFile, "%s \\\\ \\hline\n", strings.Join(columns, " & "))
			}
			fmt.Fprintln(outputFile, `\end{ltabulary}`)
			fmt.Fprintln(outputFile, `\xeCJKsetup{CheckSingle=true}`)
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{Default}{`0,`1,`2,`3,`4,`5,`6,`7,`8,`9}")
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{FullRight}{`，,`、,`。,`”}")
			fmt.Fprintln(outputFile, "\\xeCJKDeclareCharClass{FullLeft}{`：,`“}")
			fmt.Fprintln(outputFile, `\end{scriptsize}`)
			fmt.Fprintln(outputFile, `\par`)
		} else if line[0] == '+' {
			sectionType, start, title := ParseTitleLine(line)
			sectionTypeName := kSectionNames[sectionType]
			if sectionTypeToTitle[sectionType] == title {
				fmt.Printf("Ignoring %s: %s\n", sectionTypeName, title)
				continue
			}
			sectionTypeToTitle[sectionType] = title
			sectionTypeToCount[sectionType]++
			for i := sectionType + 1; i < len(kSectionNames); i++ {
				sectionTypeToTitle[i] = ""
				sectionTypeToCount[i] = 0
			}
			fmt.Printf("%s %d: %s\n", sectionTypeName, sectionTypeToCount[sectionType], title)
			fmt.Fprintf(outputFile, "%s\n\\%s{%s}\n", start, sectionTypeName, ReplaceCommentWithScript(title))
		} else {
			fmt.Fprintln(outputFile, "\\par\n"+ReplaceCommentWithScript(line))
		}

	}
	fmt.Fprintln(outputFile, `\end{document}`)
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
