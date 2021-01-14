package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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

func getFileAndCopy(videos *[]string, targetFolder string) {
	sourcePath := (*videos)[0]
	*videos = (*videos)[1:]
	source, _ := os.Open(sourcePath)
	defer source.Close()

	tmp := strings.Split(sourcePath, "/")
	filename := tmp[len(tmp)-1]
	destPath := targetFolder + "/" + filename
	dest, _ := os.Create(destPath)
	defer dest.Close()
	io.Copy(dest, source)
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

func main() {
	// videos := []string{"/Users/liujiayu/Movies/NEON/0001-0120.mp4"}
	// targetFolder := "/Users/liujiayu/Movies/Videos"
	// getFileAndCopy(&videos, targetFolder)
	// fmt.Println(videos)

	path := "/Users/liujiayu/Movies/NEON/0001-0120.mp4"
	fmt.Println(getFileCreateDate(path))
}
