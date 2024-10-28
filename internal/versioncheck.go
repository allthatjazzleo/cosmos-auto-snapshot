package internal

import (
	"fmt"
	"log"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/rootmulti"
	dbm "github.com/cosmos/cosmos-db"
)

// CheckVersion gets the height of this node
func CheckVersionAndDB(
	dataDir string,
) (int64, dbm.BackendType, error) {

	var err error
	var db dbm.DB

	// try all backends
	for _, backend := range []string{"goleveldb", "pebbledb", "rocksdb", "memdb"} {
		db, err = dbm.NewDB("application", getBackend(backend), dataDir)
		if err != nil {
			continue
		}
		log.Printf("Using backend: %s\n", backend)

		store := rootmulti.NewStore(db, cosmossdklog.NewNopLogger(), nil)

		height := store.LatestVersion()

		db.Close()
		return height, dbm.BackendType(backend), nil
	}
	return 0, dbm.BackendType(""), fmt.Errorf("failed to open db: %w", err)
}

func getBackend(backend string) dbm.BackendType {
	switch backend {
	case "goleveldb":
		return dbm.GoLevelDBBackend
	case "memdb":
		return dbm.MemDBBackend
	case "rocksdb":
		return dbm.RocksDBBackend
	case "pebbledb":
		return dbm.PebbleDBBackend
	default:
		panic(fmt.Errorf("unknown backend %s", backend))
	}
}
