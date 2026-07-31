package confdir

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"tinyrdm/backend/consts"
)

var confdirOnceValue = sync.OnceValue(func() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("get user config dir failed: %v, fallback to home dir", err)
		dir, _ = os.UserHomeDir()
		return dir
	}
	if runtime.GOOS == "darwin" {
		migratedDir, err := migrateDarwinConfigDir(dir)
		if err != nil {
			log.Printf("migrate legacy config dir failed: %v, fallback to %s", err, migratedDir)
		}
		dir = migratedDir
	}
	return dir
})

func GetConfigDir() string {
	return confdirOnceValue()
}

func migrateDarwinConfigDir(configDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return configDir, err
	}

	legacyConfigDir := filepath.Join(homeDir, "Library", "Preferences")
	if filepath.Clean(legacyConfigDir) == filepath.Clean(configDir) {
		return configDir, nil
	}

	legacyAppDataDir := filepath.Join(legacyConfigDir, consts.APP_DATA_FOLDER)
	appDataDir := filepath.Join(configDir, consts.APP_DATA_FOLDER)
	if err = migrateAppDataDir(legacyAppDataDir, appDataDir); err != nil {
		return legacyConfigDir, err
	}
	return configDir, nil
}

func migrateAppDataDir(legacyAppDataDir, appDataDir string) error {
	if filepath.Clean(legacyAppDataDir) == filepath.Clean(appDataDir) {
		return nil
	}

	legacyInfo, err := os.Stat(legacyAppDataDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !legacyInfo.IsDir() {
		return nil
	}

	if _, err = os.Stat(appDataDir); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(appDataDir), 0700); err != nil {
		return err
	}
	if err = os.Rename(legacyAppDataDir, appDataDir); err != nil {
		return err
	}
	log.Printf("migrated legacy config dir from %s to %s", legacyAppDataDir, appDataDir)
	return nil
}
