// This is the mac version

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var monthMap map[string]string = map[string]string{
	"01": "Jan",
	"02": "Feb",
	"03": "Mar",
	"04": "Apr",
	"05": "May",
	"06": "Jun",
	"07": "Jul",
	"08": "Aug",
	"09": "Sep",
	"10": "Oct",
	"11": "Nov",
	"12": "Dec"}

const concurrent int = 10

type sourceVideoFiles struct {
	mu   sync.Mutex
	list []string
}

var sourceFiles sourceVideoFiles

func isVideo(ext string) bool {
	switch ext {
	case
		".mp4",
		".MP4",
		".mov",
		".MOV":
		return true
	}
	return false
}

func getExt(name string) string {
	// filepath.Ext() has problems dealing with whitespaces
	sep := strings.Split(name, ".")
	return "." + sep[len(sep)-1]
}

func getFileModDate(path string) string {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	s := file.ModTime().String()
	date := strings.Split(s, " ")[0]
	temp := strings.Split(date, "-")
	year := temp[0]
	month := temp[1]
	day := temp[2]

	return monthMap[month] + " " + day + ", " + year
}

func locateFiles(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		ext := getExt(info.Name())
		if isVideo(ext) {
			sourceFiles.list = append(sourceFiles.list, path)
		}
	}
	return err
}

func getFileAndCopy(videos *[]string, targetFolder string, wg *sync.WaitGroup, index int) {
	defer wg.Done()
	for {
		fmt.Println("Thread ", index, "is retrieving source path... Locking...")
		sourceFiles.mu.Lock()
		if len(*videos) == 0 {
			fmt.Println("Thread ", index, "found out the work is done.")
			sourceFiles.mu.Unlock()
			return
		}
		sourcePath := (*videos)[0]
		*videos = (*videos)[1:]
		sourceFiles.mu.Unlock()
		fmt.Println("Thread ", index, "is done retrieving source path, unlocking...")

		source, _ := os.Open(sourcePath)
		defer source.Close()

		fmt.Println(sourcePath)
		newFolderName := getFileModDate(sourcePath)
		if _, err := os.Stat(targetFolder + "/" + newFolderName); os.IsNotExist(err) {
			os.Mkdir(targetFolder+"/"+newFolderName, os.ModePerm)
		}

		tmp := strings.Split(sourcePath, "/")
		filename := tmp[len(tmp)-1]
		destPath := targetFolder + "/" + newFolderName + "/" + filename
		dest, _ := os.Create(destPath)
		defer dest.Close()
		io.Copy(dest, source)
		fmt.Println("Thread ", index, "is done copying ", sourcePath)
	}
}

func startImporting(videos *[]string, targetFolder string) {
	var wg sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go getFileAndCopy(videos, targetFolder, &wg, i)
	}
	wg.Wait()
}

func main() {
	fmt.Println("Input the source folder: ")
	var sourceFolderPath string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		// when dragging in the folder to the terminal
		// there will be extra whitespace at the end
		// and "\" in between if your filenames have whitespaces
		sourceFolderPath = strings.Replace(strings.TrimSpace(scanner.Text()), "\\", "", -1)
	}

	filepath.Walk(sourceFolderPath, locateFiles)
	fmt.Println("Found ", len(sourceFiles.list), " videos.")

	fmt.Println("Input the target folder: ")
	var targetFolderPath string
	if scanner.Scan() {
		// when dragging in the folder to the terminal
		// there will be extra whitespace at the end
		// and "\" in between if your filenames have whitespaces
		targetFolderPath = strings.Replace(strings.TrimSpace(scanner.Text()), "\\", "", -1)
	}

	startImporting(&sourceFiles.list, targetFolderPath)
}
