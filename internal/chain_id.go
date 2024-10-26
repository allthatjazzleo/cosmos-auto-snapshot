package internal

import (
	"os"
	"path/filepath"

	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// GetChainID returns the chain ID of the chain
func GetChainID(homeDir string) (string, error) {
	genesis, err := os.Open(filepath.Join(homeDir, "config", "genesis.json"))
	if err != nil {
		return "", err
	}
	defer genesis.Close()
	chainID, err := genutiltypes.ParseChainIDFromGenesis(genesis)
	if err != nil {
		return "", err
	}
	return chainID, nil
}
