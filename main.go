package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type ScannerInfo struct {
	paths []string
}

func (s *ScannerInfo) ScanAdminPanels(urlStr string, writeMethod func(string)) {

	if !IsValidURL(urlStr) {
		writeMethod("url is invalid!")
	}

	for _, path := range s.paths {
		fullURL := urlStr + path

		response, err := http.Get(fullURL)
		if err != nil {
			writeMethod("[-] " + fullURL + "\n")
			continue
		}
		defer response.Body.Close()

		statusCode := response.StatusCode
		if statusCode == 200 {
			writeMethod("[+] Admin panel found: " + fullURL + "\n")
		} else if statusCode == 404 {
			writeMethod("[-] " + fullURL + "\n")
		} else if statusCode == 302 {
			writeMethod("[+] Potential EAR vulnerability found: " + fullURL + "\n")
		} else {
			writeMethod("[-] " + fullURL + "\n")
		}
	}

}

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func CreateScanner(filePath string) (*ScannerInfo, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("file is not exists!")
	}
	defer file.Close()

	var paths []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		paths = append(paths, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("error occured when file reading!")
	}

	sInstance := &ScannerInfo{
		paths: paths,
	}

	return sInstance, nil
}

func main() {

	filePath := "paths.txt" // test purpose

	scannerInstance, err := CreateScanner(filePath)

	if err != nil {
		fmt.Println("Error creating scanner:", err)
		os.Exit(1)
	}

	urlStr := "https://example.com" // test purpose

	scannerInstance.ScanAdminPanels(urlStr, func(outputStr string) {
		fmt.Print(outputStr)
	})
}
