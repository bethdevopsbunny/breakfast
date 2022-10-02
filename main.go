package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	downloadFile("wordLists/zip/secLists.zip", "https://github.com/danielmiessler/SecLists/archive/master.zip")

	unzip("wordLists/zip/secLists.zip", "wordLists/txt")

}

func checkWordListVersion() {

}

//gunzip added to unzip the kali wordlist containing rockyou
//not needed as it's included in seclists
//	downloadFile("wordLists/zip/wordlist.txt.gz", "https://gitlab.com/kalilinux/packages/wordlists/-/raw/kali/master/rockyou.txt.gz?inline=false")
//	gunzip("wordLists/zip/wordlist.txt.gz", "wordLists/txt")
func gunzip(archiveFilePath string, destination string) {

	gzipfile, err := os.Open(archiveFilePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer reader.Close()

	newfilename := strings.TrimSuffix(archiveFilePath, ".gz")
	newfilenamesplit := strings.Split(newfilename, "/")
	newfilename = newfilenamesplit[len(newfilenamesplit)-1]
	writer, err := os.Create(destination + "/" + newfilename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func unzip(archiveFilePath string, destination string) {

	archive, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		archiveFilePath := filepath.Join(destination, f.Name)
		fmt.Println("unzipping file ", archiveFilePath)

		if !strings.HasPrefix(archiveFilePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(archiveFilePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(archiveFilePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(archiveFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
