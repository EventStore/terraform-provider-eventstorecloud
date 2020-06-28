package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type tokenStore struct {
	path string
}

// Check if a token exists within the store
func (t *tokenStore) exists(audience string) bool {
	if _, err := os.Stat(t.filePath(audience)); os.IsNotExist(err) {
		return false
	}

	return true
}

// Get a token for audience from the store
func (t *tokenStore) get(audience string) (*tokenData, error) {
	tokenPath := t.filePath(audience)

	bytes, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("error reading token %q: %w", tokenPath, err)
	}

	tokenData := &tokenData{}

	err = json.Unmarshal(bytes, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("error decoding json token %q: %w", tokenPath, err)
	}

	return tokenData, nil
}

// Put a token from audience into the store
func (t *tokenStore) put(audience string, token tokenData) error {
	tokenPath := t.filePath(audience)

	bytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error serializing token: %w", err)
	}

	err = ioutil.WriteFile(tokenPath, bytes, 0600)
	if err != nil {
		return fmt.Errorf("error writing token to store %q: %w", tokenPath, err)
	}

	return nil
}

// Return token filepath for audience
func (t *tokenStore) filePath(audience string) string {
	return filepath.Join(t.path, audience)
}
