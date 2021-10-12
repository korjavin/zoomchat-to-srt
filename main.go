package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

var re = regexp.MustCompile(`^(.+):\s`)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("[ERROR] Usage `chat-formatter infile outfile")
	}
	inFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Can't open file %s: %v", os.Args[1], err)
	}
	defer inFile.Close()
	outFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("Can't create file %s: %v", os.Args[2], err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		// 20:02:25 From Teacher Lee Wright To Everyone:^M >-text
		cleanedLine := re.ReplaceAllString(line, "")
		_, err := writer.WriteString(cleanedLine + "\n")
		if err != nil {
			log.Printf("[ERROR] Can't write to outfile: %v", err)
		}
	}
	if scanner.Err() != nil {
		log.Fatalf("Can't finish file reading: %v", err)
	}
}
