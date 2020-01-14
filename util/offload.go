package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fsnotify/fsnotify"
)

func watcher(wg *sync.WaitGroup) {
	defer wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err = Log(err); err != nil {
		return
	}
	defer watcher.Close()

	if err = Log(watcher.Add(watch)); err != nil {
		return
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			if strings.Contains(event.Op.String(), "CREATE") || strings.Contains(event.Op.String(), "WRITE") || strings.Contains(event.Op.String(), "RENAME") || strings.Contains(event.Op.String(), "CHMOD") {
				fileQueue[event.Name] = false
			}
		case err := <-watcher.Errors:
			Log(err)
		}
	}

}

func offloader(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		for filename := range fileQueue {
			if _, err := offload(filename); err == nil {
				delete(fileQueue, filename)
			}
		}
		time.Sleep(queueSpeed * time.Millisecond)
	}

}

func offload(filename string) (url string, err error) {
	extension := filepath.Ext(filename)
	parts := strings.Split(filename, "/")
	key := strings.TrimSuffix(parts[len(parts)-1], extension)
	if url, err = saveToS3(filename, fmt.Sprintf("%s/%s", prefix, key)); err == nil {
		Log(fmt.Sprintf("Uploaded %s to %s", filename, url))
	}
	Log(err)

	return
}

func saveToS3(filename string, key string) (url string, err error) {
	file, err := os.Open(filename)
	if err = Log(err); err != nil {
		return "", Log(errors.New(fmt.Sprintf("Failed to open %s", filename)))
	}

	defer file.Close()

	fileData, err := ioutil.ReadFile(filename)
	if err = Log(err); err != nil {
		return "", errors.New("Could not read KYC file")
	}

	Log(fmt.Sprintf("Uploading %s to %s/%s/%s", filename, bucket, prefix, key))

	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ACL:         &acl,
		ContentType: aws.String(http.DetectContentType(fileData)),
	})

	if err = Log(err); err != nil {
		return "", errors.New("Could not upload file")
	}

	url = result.Location

	// Delete file from disk
	Log(os.Remove(filename))
	return
}
