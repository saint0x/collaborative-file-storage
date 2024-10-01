package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type B2Service struct {
	accountID      string
	applicationKey string
	apiURL         string
	authToken      string
	bucketID       string
	client         *http.Client
}

func NewB2Service(accountID, applicationKey, bucketID string) (*B2Service, error) {
	s := &B2Service{
		accountID:      accountID,
		applicationKey: applicationKey,
		bucketID:       bucketID,
		client:         &http.Client{Timeout: 30 * time.Second},
	}

	if err := s.authorize(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *B2Service) authorize() error {
	req, err := http.NewRequest("GET", "https://api.backblazeb2.com/b2api/v2/b2_authorize_account", nil)
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.accountID + ":" + s.applicationKey))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var authResponse struct {
		ApiUrl             string `json:"apiUrl"`
		AuthorizationToken string `json:"authorizationToken"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return err
	}

	s.apiURL = authResponse.ApiUrl
	s.authToken = authResponse.AuthorizationToken

	return nil
}

func (s *B2Service) UploadFile(ctx context.Context, key string, body io.Reader) error {
	// First, get upload URL and authorization token
	uploadURL, uploadAuthToken, err := s.getUploadURL()
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", uploadAuthToken)
	req.Header.Set("X-Bz-File-Name", key)
	req.Header.Set("Content-Type", "b2/x-auto")
	req.Header.Set("X-Bz-Content-Sha1", "do_not_verify")

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func (s *B2Service) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/b2api/v2/b2_download_file_by_name?bucketName=%s&fileName=%s", s.apiURL, s.bucketID, key)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	return resp.Body, nil
}

func (s *B2Service) DeleteFile(ctx context.Context, key string) error {
	// First, get file info to retrieve fileId
	fileInfo, err := s.getFileInfo(ctx, key)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/b2api/v2/b2_delete_file_version", s.apiURL)
	data := map[string]string{
		"fileName": key,
		"fileId":   fileInfo.FileID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete file: %s", resp.Status)
	}

	return nil
}

func (s *B2Service) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	url := fmt.Sprintf("%s/b2api/v2/b2_list_file_names", s.apiURL)
	data := map[string]string{
		"bucketId": s.bucketID,
		"prefix":   prefix,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list files: %s", resp.Status)
	}

	var result struct {
		Files []struct {
			FileName string `json:"fileName"`
		} `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var files []string
	for _, file := range result.Files {
		files = append(files, file.FileName)
	}

	return files, nil
}

func (s *B2Service) GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	url := fmt.Sprintf("%s/b2api/v2/b2_get_download_authorization", s.apiURL)
	data := map[string]interface{}{
		"bucketId":               s.bucketID,
		"fileNamePrefix":         key,
		"validDurationInSeconds": int(expiration.Seconds()),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get download authorization: %s", resp.Status)
	}

	var result struct {
		AuthorizationToken string `json:"authorizationToken"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	downloadURL := fmt.Sprintf("%s/file/%s/%s?Authorization=%s", s.apiURL, s.bucketID, key, result.AuthorizationToken)
	return downloadURL, nil
}

func (s *B2Service) getUploadURL() (string, string, error) {
	url := fmt.Sprintf("%s/b2api/v2/b2_get_upload_url", s.apiURL)
	data := map[string]string{
		"bucketId": s.bucketID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get upload URL: %s", resp.Status)
	}

	var result struct {
		UploadUrl          string `json:"uploadUrl"`
		AuthorizationToken string `json:"authorizationToken"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result.UploadUrl, result.AuthorizationToken, nil
}

func (s *B2Service) getFileInfo(ctx context.Context, key string) (*struct {
	FileID string `json:"fileId"`
}, error) {
	url := fmt.Sprintf("%s/b2api/v2/b2_get_file_info", s.apiURL)
	data := map[string]string{
		"fileName": key,
		"bucketId": s.bucketID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file info: %s", resp.Status)
	}

	var result struct {
		FileID string `json:"fileId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Add the Close method
func (s *B2Service) Close() error {
	// Perform any necessary cleanup
	s.client.CloseIdleConnections()
	// Reset auth token
	s.authToken = ""
	return nil
}
