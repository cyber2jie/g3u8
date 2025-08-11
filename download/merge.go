package download

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type M3u8Merge interface {
	Merge()
}

type SimpleMerge struct {
	playList []*PlayList
}

func NewSimpleMerge(playList []*PlayList) *SimpleMerge {
	return &SimpleMerge{
		playList: playList,
	}
}
func GetFileName(p string) string {
	unixName := path.Base(p)
	if unixName == p || !strings.ContainsRune(p, '\\') {
		return unixName
	}
	lastIndex := strings.LastIndex(p, "\\")
	if lastIndex == -1 {
		return p
	}
	return p[lastIndex+1:]
}
func (s *SimpleMerge) Merge() {
	outdir := GetDefaultM3u8DownloadContext().OutDir
	cacheDir := path.Join(outdir, Cache_Dir)
	tsDir := path.Join(cacheDir, Ts_Dir)
	outName := GetFileName(GetDefaultM3u8DownloadContext().Options.Out)
	if outName == "" || !strings.Contains(outName, ".") {
		outName = "out.mp4"
	}
	outName = path.Join(outdir, outName)

	outFile, err := os.Create(outName)
	if err != nil {
		log.Fatalf("create file error: %v", err)
	}
	defer outFile.Close()

	for _, playList := range s.playList {
		tsFile := path.Join(tsDir, fmt.Sprintf("%d.ts", playList.Index))
		file, err := os.Open(tsFile)
		if err != nil {
			log.Fatalf("open file error: %v", err)
		}
		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Fatalf("copy file error: %v", err)
		}
		file.Close()
	}

}
