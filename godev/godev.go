package godev

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/heyvito/goup/models"
	"github.com/levigross/grequests"
	"io"
	"os"
)

func ListVersions(all bool) ([]models.RemoteVersion, error) {
	url := "https://go.dev/dl/?mode=json"
	if all {
		url += "&include=all"
	}
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed obtaining list from go.dev: %w", err)
	}

	if !resp.Ok {
		return nil, fmt.Errorf("failed obtaining list from go.dev: HTTP %d. Please try again later.", resp.StatusCode)
	}

	var list []models.RemoteVersion
	if err = json.NewDecoder(resp).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed decoding data from go.dev: %s", err)
	}

	return list, nil
}

func DownloadVersion(url string) (string, error) {
	target, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	if err = target.Close(); err != nil {
		return "", err
	}

	resp, err := grequests.Get(url, &grequests.RequestOptions{})
	if err != nil {
		_ = os.Remove(target.Name())
		return "", fmt.Errorf("failed downloading %s: %w", url, err)
	}
	if !resp.Ok {
		_ = os.Remove(target.Name())
		return "", fmt.Errorf("failed downloading %s: HTTP %d", url, resp.StatusCode)
	}

	if err = resp.DownloadToFile(target.Name()); err != nil {
		_ = os.Remove(target.Name())
		return "", fmt.Errorf("failed downloading %s: %w", url, err)
	}

	return target.Name(), nil
}

func CheckShasum(path, sum string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	shaBuf := sha256.New()
	if _, err = io.Copy(shaBuf, f); err != nil {
		return false, err
	}

	return sum == hex.EncodeToString(shaBuf.Sum(nil)), nil
}
