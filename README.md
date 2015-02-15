This project contains ancient Chinese books in plain text. We provide converters written in Go programming language to convert plain texts to other format like pdf.


MANUAL: 
1. Install golang, git: 
  $ sudo apt-get install golang-go git 
2. Install texlive: https://www.tug.org/texlive/acquire-netinstall.html 
  After installation (takes hours!), add the bin folder to path, like 
  $ export PATH=/usr/local/texlive/2014/bin/x86_64-linux:$PATH 
  NOTE: Your installation path might be different. Use "ls /usr/local/texlive" to find yours. You can put the above export command in ~/.bashrc to run it automatically. 
3. Download source code of project ancient-chinese: 
  $ git clone https://code.google.com/p/ancient-chinese 
  $ cd ancient-chinese 
4. Install fonts: 
  $ sudo mkdir /usr/share/fonts/truetype/chinese/ 
  $ sudo cp fonts/* /usr/share/fonts/truetype/chinese/ 
5. Compile ancient-chinese: 
  $ cd go 
  $ go install tex 
6. Convert txt to tex format: 
  $ cd ../txt 
  $ ../go/bin/tex shiji-simplified.txt 
7. Convert tex to pdf format: 
  $ for i in 1 2 3 ; do xelatex shiji-simplified.txt ; done 
  NOTE: We need to xelatex three times to correctly generate TOC (table of content): 
  1st run: generate all pages w/o TOC. 
  2nd run: generate TOC and all pages w/o correct page numbers. 
  3nd run: generate TOC and all pagew w/ correct page numbers. 
  You don't need to worry about these details, Just run xelatex three times. \


Text file rules: 
1. All text files go under txt sub folder and Go code under go sub folder. 
2. Only a subset of ASCII characters are allowed in file names, including lowcase letters, numbers, - (dash) and . (dot). 
3. Use pinyin to replace Chinese character in file names. For example, "shiji" for "史记". Use suffix like "-simplified" or "-traditional" to indicate that the text is in simplified or traditional Chinese.  Prefer ".txt" extension. 
  e.g.  shiji-simplified.txt shiji-traditional.txt
4. All files are encoded with UTF8, W/O BOM byte. 
5. Rare characters are represented by multiple characters, enclosed by half-width parentheses. 
   e.g. 
     (土慮)   -- left & right composition. 
     (/窮)   -- / means up & down composition. 
     (𠂆*圭)  -- * means outside / inside composition. 
   NOTE: Rare characters are defined by that they are not included in 
   HanaMin(花園明朝) font, see http://www.zdic.net/appendix/f18.htm 
6. Comments are put inside （）. 

Go code rules: 
1. All Go code need to formated with gofmt. 
2. Every Go source file should have minimal comment explaining how to use the code. 


Format of text files: 
TITLE 
AUTHOR 
++CHAPTER1       // Num of '+' determines the type of structure, the more the smaller. Max 7. 
CONTENT          // One line for each paragraph. 
.... 
++CHAPTER2 
....

Table format: 
--- 
Column1|Column2|Column3 
Column1|Column2|Column3 
.... 
--- 
Note: 
1. Tables start and end with "---". 
2. Every row must have the same number of "|". 


FAQs 
Q: Why use text files? 
A: Text files save disk space. Most importantly, it's easy to edit text files and we can use source control system to record the change history of files. We require minimal formating in the text files and use Tex to format the books for different devices. 
Q: Why use go? 
A: No special reason. The author would like to try this relatively new language, which is cool and very concise. 
