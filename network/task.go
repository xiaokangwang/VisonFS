package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type NetworkTaskQueue struct {
	srv *drive.Service
}

type NetworkUploadTask struct {
	Filename string
	Content  []byte
}
type NetworkUploadTaskResult struct {
}
type NetworkDownloadTask struct {
	Filename string
}
type NetworkDownloadTaskResult struct {
	Content []byte
}
type NetworkListTask struct {
	Dir string
}
type NetworkListTaskResult struct {
	Files []string
}

func (ntq *NetworkTaskQueue) EnqueueUploadTask(task NetworkUploadTask) {
	var err error
	var file drive.File
	file.Name = task.Filename
	_, err = ntq.srv.Files.Create(&file).Media(bytes.NewReader(task.Content)).Do()
	if err != nil {
		panic(err)
	}
}
func (ntq *NetworkTaskQueue) EnqueueDownloadTask(task NetworkDownloadTask) NetworkDownloadTaskResult {
	fn := task.Filename
	r, err := ntq.srv.Files.List().Q("name = '" + fn + "'").PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	did := r.Files[0].Id
	resp, err := ntq.srv.Files.Get(did).AcknowledgeAbuse(true).Download()
	if err != nil {
		panic(err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	var nt NetworkDownloadTaskResult
	nt.Content = c
	return nt
}
func (ntq *NetworkTaskQueue) EnqueueListTask(task NetworkListTask) NetworkListTaskResult {
	panic(nil)
}

func NewNetworkTaskQueue() *NetworkTaskQueue {
	return &NetworkTaskQueue{}
}
func (ntq *NetworkTaskQueue) ensureToken() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	srv, err := drive.New(getClient(config))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	ntq.srv = srv
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}
