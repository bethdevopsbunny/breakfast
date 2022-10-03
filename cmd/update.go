package cmd

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{

	Use:   "update",
	Short: "checks for and updates store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		d, _ := getLatestReleaseData("danielmiessler", "SecLists")

		timee := d.PublishedAt.UnixMicro()
		zipfilepath := fmt.Sprintf("store/zip/%d-words.zip", timee)

		if _, err := os.Stat(zipfilepath); err == nil {

			println("Already Uptodate")

		} else {

			downloadFile(zipfilepath, d.ZipballURL)

			txtfilepath := fmt.Sprintf("store/txt/%d-words", timee)
			zipRoot := unzip(zipfilepath, "store/txt")
			zipRootPath := fmt.Sprintf("store/txt/%s", zipRoot)
			os.Rename(zipRootPath, txtfilepath)
			println("completed")

		}

		sa := ReadEachLine("store/txt/1659433997000000-words/Passwords/darkweb2017-top10000.txt")

		for _, element := range sa {
			println(fmt.Sprintf("%s:%s", element, GetMD5Hash(element)))
		}

	},
}

func ReadEachLine(filepath string) (fileLines []string) {

	readFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//gunzip added to unzip the kali wordlist containing rockyou
//not needed as it's included in seclists
//	downloadFile("store/zip/wordlist.txt.gz", "https://gitlab.com/kalilinux/packages/wordlists/-/raw/kali/master/rockyou.txt.gz?inline=false")
//	gunzip("store/zip/wordlist.txt.gz", "store/txt")
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

//unzip the provided archive at the provided destination
//returns the zips first/root name to be used for renaming if needed.
func unzip(archiveFilePath string, destination string) (zipRoot string) {

	archive, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	//dynamically pulls the root file name to return for renaming
	zipRoot = archive.Reader.File[0].FileHeader.Name

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
	return zipRoot
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

type ReleaseData struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string        `json:"node_id"`
	TagName         string        `json:"tag_name"`
	TargetCommitish string        `json:"target_commitish"`
	Name            string        `json:"name"`
	Draft           bool          `json:"draft"`
	Prerelease      bool          `json:"prerelease"`
	CreatedAt       time.Time     `json:"created_at"`
	PublishedAt     time.Time     `json:"published_at"`
	Assets          []interface{} `json:"assets"`
	TarballURL      string        `json:"tarball_url"`
	ZipballURL      string        `json:"zipball_url"`
	Body            string        `json:"body"`
	Reactions       struct {
		URL        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Num1       int    `json:"+1"`
		Num10      int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions"`
	MentionsCount int `json:"mentions_count"`
}

func getLatestReleaseData(owner string, repoName string) (ReleaseData, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ReleaseData{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ReleaseData{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReleaseData{}, err
	}

	var releaseData ReleaseData
	err = json.Unmarshal(body, &releaseData)
	if err != nil {
		return ReleaseData{}, err
	}

	return releaseData, nil
}

type ReleaseAssets []struct {
	URL                string    `json:"url"`
	BrowserDownloadURL string    `json:"browser_download_url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	State              string    `json:"state"`
	ContentType        string    `json:"content_type"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"uploader"`
}

func getReleaseAssets(owner string, repoName string, releaseID int) (ReleaseAssets, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d/assets", owner, repoName, releaseID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ReleaseAssets{}, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ReleaseAssets{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReleaseAssets{}, err
	}
	var releaseAssets ReleaseAssets

	err = json.Unmarshal(body, &releaseAssets)
	if err != nil {
		return ReleaseAssets{}, err
	}

	return releaseAssets, nil
}

type ReleaseAsset struct {
	URL                string    `json:"url"`
	BrowserDownloadURL string    `json:"browser_download_url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	State              string    `json:"state"`
	ContentType        string    `json:"content_type"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"uploader"`
}

func getLatestReleaseAsset(owner string, repoName string, assetID int) (ReleaseAsset, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/assets/%d", owner, repoName, assetID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ReleaseAsset{}, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ReleaseAsset{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReleaseAsset{}, err
	}
	var releaseAsset ReleaseAsset

	err = json.Unmarshal(body, &releaseAsset)
	if err != nil {
		return ReleaseAsset{}, err
	}

	return releaseAsset, nil
}
