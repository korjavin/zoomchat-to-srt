package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	timeFormat = "15:04:05"
	timeOutput = "15:04:05,000"
	maxTime    = time.Second * 10
	minTime    = time.Second
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("[ERROR] Usage `binary zoomfile")
	}
	inFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Can't open file %s: %v", os.Args[1], err)
	}
	defer inFile.Close()
	ext := path.Ext(os.Args[1])
	outFileName := os.Args[1][0:len(os.Args[1])-len(ext)] + ".srt"
	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("Can't create file %s: %v", outFileName, err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	shiftDuration := time.Second * 0
	if len(os.Args) > 2 { // shift time to args[2]
		shiftTime, err := time.Parse(timeFormat, os.Args[2])
		if err != nil {
			log.Fatalf("[ERROR] Can't parse shift time:%s %v", os.Args[2], err)
		}
		begin, _ := time.Parse(timeFormat, "0:0:0")
		shiftDuration = shiftTime.Sub(begin)
	}

	scanner := bufio.NewScanner(inFile)
	lineNumber := 0
	var timestamp, text string
	for scanner.Scan() {
		line := scanner.Text()
		// 20:02:25 From Teacher Lee Wright To Everyone:^M >-text
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			log.Printf("[Warn] Skipped wrong line %s", line)
			continue
		}
		if lineNumber > 0 {
			// We should have previous values, so we can calculate
			//	duration for previous message and write it
			tsPrev, err := time.Parse(timeFormat, timestamp)
			if err != nil {
				log.Fatalf("[ERROR] Can't parse time:%s %v", timestamp, err)
			}
			tsNow, err := time.Parse(timeFormat, parts[0])
			if err != nil {
				log.Fatalf("[ERROR] Can't parse time:%s %v", parts[0], err)
			}
			duration := tsNow.Sub(tsPrev)
			if duration > maxTime { // Close after 30 sec
				tsNow = tsPrev.Add(maxTime)
			}
			if duration < minTime { // Combine text with previos one
				text = fmt.Sprintf("%s\n%s", text, parts[1])
				continue
			}
			cleanedLine := fmt.Sprintf("%d\n%s --> %s\n%s\n",
				lineNumber,
				tsPrev.Add(-shiftDuration).Format(timeOutput),
				tsNow.Add(-shiftDuration).Format(timeOutput),
				fmt.Sprintf("%s\n", text),
			)

			_, err = writer.WriteString(cleanedLine)
			if err != nil {
				log.Printf("[ERROR] Can't write to outfile: %v", err)
			}
		}

		lineNumber++
		timestamp, text = parts[0], parts[1]
	}
	if scanner.Err() != nil {
		log.Fatalf("Can't finish file reading: %v", err)
	}
}
