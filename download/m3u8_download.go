package download

import (
	"context"
	"errors"
	"fmt"
	"g3u8/config"
	"g3u8/m3u8"
	"g3u8/util"
	"log"
	urlutil "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"
)

type M3u8DownloadOptions struct {
	Url     string
	Out     string
	BaseUrl string
}

func M3u8Download(options M3u8DownloadOptions) error {

	dir := filepath.Dir(options.Out)
	cache, err := ResolveCache(dir)

	if err != nil {
		return err
	}
	//restore old state
	if cache != nil && cache.Manifest != nil {
		if cache.Complete {
			options.Url = cache.Manifest.M3u8Url
			InitContext(context.Background(), config.Config, options, dir)
			NewSimpleMerge(cache.PlayLists).Merge()
			CleanCache(dir)
			log.Printf("merge file,%s ", dir)
			return nil
		} else {
			log.Printf("Resume downloading...")
			options.Url = cache.Manifest.M3u8Url
		}
	}
	InitContext(context.Background(), config.Config, options, dir)

	return recoverOrNewDownload(dir, cache, options)

}

func recoverOrNewDownload(dir string, cache *DownloadCache, options M3u8DownloadOptions) error {
	mediaList, err := ParseListFromUrl(options.Url)

	if err != nil {
		return err
	}
	if mediaList == nil {
		return errors.New("Can not parse m3u8 list")
	}

	if cache != nil {

		log.Printf("reset cache dir")
		ResetCacheDir(dir, mediaList.String())

	} else {

		manifest := Manifest{
			M3u8Url:  options.Url,
			Out:      options.Out,
			CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		InitCacheDir(dir, mediaList.String(), manifest)

		cache = &DownloadCache{
			Complete:  false,
			Manifest:  &manifest,
			PlayLists: mediaPlayListToPlayList(mediaList),
		}
		SaveCache(dir, cache)
	}

	//handle ctrl+c
	exitSignalChan := make(chan os.Signal, 1)
	signal.Notify(exitSignalChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	workerNum := GetDefaultM3u8DownloadContext().Config.Worker.MaxWorkers

	queueSize := GetDefaultM3u8DownloadContext().Config.Worker.QueueSize

	works := make([]*Worker, workerNum)

	for i := 0; i < workerNum; i++ {

		work := NewWorker(ctx, fmt.Sprintf("worker-%d", (i+1)))
		works[i] = work

		for j := 0; j < queueSize; j++ {
			for _, playList := range cache.PlayLists {
				if playList.WorkerId != "" || playList.Download {
					continue
				}
				work.AddJob(playList)
				playList.WorkerId = work.Name
				break
			}
		}
	}

	log.Printf("start workers")
	for i := 0; i < workerNum; i++ {
		go works[i].Start()
	}

	//distribute jobs
	go func() {

		cases := make([]reflect.SelectCase, len(works))

		for i := 0; i < len(works); i++ {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(works[i].Wait()),
			}
		}

		for {

			chose, _, recvOk := reflect.Select(cases)
			if !recvOk {
				log.Printf("worker %s done", works[chose].Name)
				continue
			}

			allDistribute := true
			for _, playList := range cache.PlayLists {
				if playList.WorkerId != "" || playList.Download {
					continue
				}
				works[chose].AddJob(playList)
				playList.WorkerId = works[chose].Name
				allDistribute = false
				break
			}

			if allDistribute {
				allDone := true
				for _, playList := range cache.PlayLists {
					if !playList.Download {
						allDone = false
						break
					}
				}
				if allDone {
					cancel()
					break
				}
			}

			remain := 0
			for _, playList := range cache.PlayLists {
				if !playList.Download {
					remain++
				}
			}
			log.Printf("remain %d playlist files", remain)
		}

	}()

	//wait jobs done
	duration := time.Duration(GetDefaultM3u8DownloadContext().Config.Worker.SaveStateDuration)
	if duration <= 0 {
		duration = time.Duration(config.Save_State_Durition)
	}
	ticker := time.NewTicker(time.Second * duration)

out:
	for {
		select {

		case <-ctx.Done():
			log.Printf("download complete")
			cache.Complete = true
			SaveCache(dir, cache)
			break out
		case <-exitSignalChan:
			log.Printf("save state")
			SaveCache(dir, cache)
			return errors.New("download interrupt")
		case <-ticker.C:
			log.Printf("save state")
			SaveCache(dir, cache)
		}
	}

	//merge
	log.Printf("merge ts files")

	NewSimpleMerge(cache.PlayLists).Merge()
	//CleanCache(dir)

	return nil
}

func mediaPlayListToPlayList(mediaPlayList *m3u8.MediaPlaylist) []*PlayList {
	playLists := make([]*PlayList, 0)

	var key = EncryptKey{}

	for _, segment := range mediaPlayList.Segments {

		if segment == nil {
			continue
		}

		if segment.Key != nil {

			var iv string
			b, err := util.Hex2Byte(segment.Key.IV)

			if err == nil {
				iv = util.Byte2Hex(b)
			}

			key = EncryptKey{
				IVHex:             iv,
				KeyHex:            getKeyHex(segment.Key.URI),
				Keyformat:         segment.Key.Keyformat,
				Keyformatversions: segment.Key.Keyformatversions,
				Method:            segment.Key.Method,
			}
		}

		playLists = append(playLists, &PlayList{
			Index:    int(segment.SeqId),
			Duration: segment.Duration,
			Url:      segment.URI,
			Key:      &key,
			Download: false,
		})
	}
	return playLists
}

func getKeyHex(url string) string {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		baseUrl := GetDefaultM3u8DownloadContext().Options.BaseUrl
		if baseUrl == "" {
			uri, err := urlutil.Parse(GetDefaultM3u8DownloadContext().Options.Url)
			if err == nil {
				baseUrl = uri.Scheme + "://" + uri.Host
			}
		}
		_url, err := urlutil.JoinPath(baseUrl, url)
		if err == nil {
			url = _url
		}

	}
	resp, err := GetDefaultM3u8DownloadContext().Client.Get(url)
	if err == nil {
		return util.Byte2Hex(resp.Body)
	}
	return ""
}
