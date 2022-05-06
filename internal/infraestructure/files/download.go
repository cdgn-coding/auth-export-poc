package files

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader() Downloader {
	client := new(http.Client)
	return Downloader{
		client,
	}
}

func (d Downloader) Download(url, filepath string) error {
	out, errorCreatingFile := os.Create(filepath)
	defer out.Close()

	if errorCreatingFile != nil {
		return errorCreatingFile
	}

	request, errorCreatingRequest := http.NewRequest("GET", url, nil)
	if errorCreatingRequest != nil {
		return errorCreatingRequest
	}
	request.Header.Add("Accept-Encoding", "gzip")
	resp, errorRequestingFile := d.client.Do(request)
	defer resp.Body.Close()

	uncompressedBody, errorUncompressing := gzip.NewReader(resp.Body)
	defer uncompressedBody.Close()

	if errorUncompressing != nil {
		return errorUncompressing
	}

	if errorRequestingFile != nil {
		return errorRequestingFile
	}

	_, errorCopyingFile := io.Copy(out, uncompressedBody)

	if errorCopyingFile != nil {
		return errorCopyingFile
	}

	return nil
}
