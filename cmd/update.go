package cmd

import (
	"archive/zip"
	"compress/gzip"
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
	Short: "checks for and updates wordlists",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		downloadFile("wordlists/zip/secLists.zip", "https://github.com/danielmiessler/SecLists/archive/master.zip")
		unzip("wordlists/zip/secLists.zip", "wordlists/txt")

	},
}

//gunzip added to unzip the kali wordlist containing rockyou
//not needed as it's included in seclists
//	downloadFile("wordlists/zip/wordlist.txt.gz", "https://gitlab.com/kalilinux/packages/wordlists/-/raw/kali/master/rockyou.txt.gz?inline=false")
//	gunzip("wordlists/zip/wordlist.txt.gz", "wordlists/txt")
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

	res, err := http.NewRequest("GET", url, nil)
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

func getLatestReleaseAsset(owner string, repoName string) (ReleaseData, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repoName)

	res, err := http.NewRequest("GET", url, nil)
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
