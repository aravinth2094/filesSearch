/*
	Author	:	Aravinth Sundaram
	Created	:	29 July 2020
	License	:	MIT
*/

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// VERSION of this software
const VERSION = "1.1.0"

type appConfigProperties map[string]string

var cache map[string]int = make(map[string]int)

const layoutISO = "2006-01-02"

type searchedFile struct {
	Path string
	Info os.FileInfo
}

var wg sync.WaitGroup

func readPropertiesFile(filename string) (appConfigProperties, error) {
	config := appConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}

func fileCount(path string) (int, error) {
	if cache[path] != 0 {
		return cache[path], nil
	}
	i := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	cache[path] = i
	return i, nil
}

func write(writer *csv.Writer, dataChannel <-chan searchedFile) {
	for data := range dataChannel {
		count, _ := fileCount(filepath.Dir(data.Path))
		writer.Write([]string{data.Path, data.Info.ModTime().String(), fmt.Sprintf("%d", data.Info.Size()), fmt.Sprintf("%d", count)})
	}
	wg.Done()
}

func main() {
	log.Println(os.Args[0], "Version", VERSION)
	log.Println("Loading configuration from search.properties")
	properties, err := readPropertiesFile("search.properties")
	if err != nil {
		log.Fatalln(err)
	}
	outputFile := "output_" + time.Now().Format(layoutISO) + "_" + time.Now().Format("15-04-05") + ".csv"
	log.Println("Prepping output file:", outputFile)
	file, ferr := os.Create(outputFile)
	if ferr != nil {
		log.Fatalln(ferr)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"FILE_NAME", "FILE_MODIFIED_TIME", "FILE_SIZE", "DIRECTORY_FILES_COUNT"})
	dataChannel := make(chan searchedFile, 1000)
	wg.Add(1)
	go write(writer, dataChannel)
	log.Println("Searching...")
	filesFound := 0
	filesSearched := 0
	start := time.Now()
	for directory, dateStr := range properties {
		if "today" == dateStr {
			dateStr = time.Now().Format(layoutISO)
		}
		date, parseError := time.ParseInLocation(layoutISO, dateStr, time.Local)
		if parseError != nil {
			log.Fatalln(parseError)
		}
		log.Println("Directory:", directory, "=>", "Date:", date.Format(layoutISO))
		_, statErr := os.Stat(directory)
		if os.IsNotExist(statErr) {
			log.Fatalln(directory, "does not exist")
		}
		err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			filesSearched = filesSearched + 1
			if err != nil {
				log.Printf("Unable to access %q: %v\n", path, err)
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if info.ModTime().Format(layoutISO) != date.Format(layoutISO) {
				return nil
			}
			dataChannel <- searchedFile{Path: path, Info: info}
			filesFound = filesFound + 1
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the path %q: %v\n", directory, err)
		}
	}
	close(dataChannel)
	wg.Wait()
	log.Println(filesFound, "/", filesSearched, "found in", len(cache), "subdirectories in", (time.Now().Sub(start)))
	log.Println("Program Completed")
}
