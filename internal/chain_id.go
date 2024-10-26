package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type genesis struct {
	ChainID string `json:"chain_id"`
}

// GetChainID returns the chain ID of the chain
func GetChainID(homeDir string) (string, error) {
	genesisFile, err := os.Open(filepath.Join(homeDir, "config", "genesis.json"))
	if err != nil {
		return "", err
	}
	defer genesisFile.Close()

	var g genesis
	err = json.NewDecoder(genesisFile).Decode(&g)
	if err != nil {
		return "", err
	}
	return g.ChainID, nil
}
