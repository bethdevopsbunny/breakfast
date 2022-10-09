package source

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Gunzip added to unzip the kali wordlist containing rockyou
//not needed as it's included in seclists
//	downloadFile("store/zip/wordlist.txt.gz", "https://gitlab.com/kalilinux/packages/wordlists/-/raw/kali/master/rockyou.txt.gz?inline=false")
//	gunzip("store/zip/wordlist.txt.gz", "store/txt")
func Gunzip(archiveFilePath string, destination string) {

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

//Unzip the provided archive at the provided destination
//returns the zips first/root name to be used for renaming if needed.
func Unzip(archiveFilePath string, destination string) (zipRoot string) {

	archive, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	//dynamically pulls the root file name to return for renaming
	zipRoot = archive.Reader.File[0].FileHeader.Name

	for i, f := range archive.File {
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
		fmt.Sprintf("unzipped %d of %d", i, len(archive.File))
	}
	return zipRoot
}
