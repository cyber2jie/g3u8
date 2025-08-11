package download

import (
	"github.com/bytedance/sonic"
	"os"
	"path"
)

const (
	Cache_Dir     = ".g3u8Cache"
	Ts_Dir        = "ts"
	M3u8_File     = "m3u8.txt"
	Manifest_File = "manifest.json"
	Download_File = "download.json"
)

type Manifest struct {
	M3u8Url  string `json:"m3u8_url"`
	Out      string `json:"out"`
	CreateAt string `json:"create_at"`
}

type DownloadCache struct {
	Manifest  *Manifest   `json:"manifest"`
	Complete  bool        `json:"complete"`
	PlayLists []*PlayList `json:"play_lists"`
}
type PlayList struct {
	Index    int         `json:"index"`
	Url      string      `json:"url"`
	Duration float64     `json:"duration"`
	Key      *EncryptKey `json:"key"`
	Download bool        `json:"download"`
	WorkerId string      `json:"-"`
}

type EncryptKey struct {
	Method            string `json:"method"`
	Keyformat         string `json:"keyformat"`
	Keyformatversions string `json:"keyformatversions"`
	KeyHex            string `json:"keyhex"`
	IVHex             string `json:"ivhex"`
}

func ResolveCache(dir string) (*DownloadCache, error) {
	if dir != "" {
		cacheDir := path.Join(dir, Cache_Dir)
		if IsDirExists(cacheDir) {
			downloadFile := path.Join(cacheDir, Download_File)
			if IsFileExists(downloadFile) {

				b, err := os.ReadFile(downloadFile)

				if err != nil {
					return nil, err
				}

				downloadCache := &DownloadCache{}

				err = sonic.Unmarshal(b, downloadCache)

				if err != nil {
					return nil, err
				}
				return downloadCache, nil
			}

		}
	}
	return nil, nil
}

func SaveCache(dir string, cache *DownloadCache) error {
	cacheDir := path.Join(dir, Cache_Dir)
	downloadFile := path.Join(cacheDir, Download_File)

	b, err := sonic.Marshal(cache)

	if err != nil {
		return err
	}
	return os.WriteFile(downloadFile, b, os.ModePerm)
}

func CleanCache(dir string) error {
	if dir != "" {
		return os.RemoveAll(path.Join(dir, Cache_Dir))
	}
	return nil
}
func IsDirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ResetCacheDir(dir string, m3u8List string) error {
	cacheDir := path.Join(dir, Cache_Dir)
	m3u8File := path.Join(cacheDir, M3u8_File)
	err := os.WriteFile(m3u8File, []byte(m3u8List), os.ModePerm)

	if err != nil {
		return err
	}
	return nil
}

func InitCacheDir(dir string, m3u8List string, manifest Manifest) error {
	cacheDir := path.Join(dir, Cache_Dir)

	tsDir := path.Join(cacheDir, Ts_Dir)
	err := os.MkdirAll(tsDir, os.ModePerm)
	if err != nil {
		return err
	}

	m3u8File := path.Join(cacheDir, M3u8_File)
	err = os.WriteFile(m3u8File, []byte(m3u8List), os.ModePerm)

	if err != nil {
		return err
	}

	manifestFile := path.Join(cacheDir, Manifest_File)

	b, err := sonic.Marshal(manifest)
	if err != nil {
		return err
	}
	return os.WriteFile(manifestFile, b, os.ModePerm)
}
