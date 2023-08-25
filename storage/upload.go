package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func Upload(ctx context.Context, url string, content []byte) error {
	// upload content to url as S3 object using PUT
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Upload failed: %s", res.Status)
	}
	return err
}
