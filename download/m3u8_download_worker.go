package download

import (
	"context"
	"fmt"
	"g3u8/util"
	"log"
	urlutil "net/url"
	"os"
	"path"
	"strings"
)

type Worker struct {
	jobs    chan *PlayList
	context context.Context
	notFull chan string
	Name    string
}

func NewWorker(ctx context.Context, name string) *Worker {
	return &Worker{
		jobs:    make(chan *PlayList, GetDefaultM3u8DownloadContext().Config.Worker.QueueSize),
		context: ctx,
		notFull: make(chan string, 1),
		Name:    name,
	}
}

func (w *Worker) AddJob(job *PlayList) {
	w.jobs <- job
}

func (w *Worker) Start() {
	log.Printf("worker %s start", w.Name)
	for {
		select {
		case job := <-w.jobs:
			if job != nil {

				log.Printf("worker %s start download playList %d", w.Name, job.Index)

				w.DownLoad(job)

				w.notFull <- w.Name
			}
		case <-w.context.Done():
			return
		}
	}
}
func (w *Worker) DownLoad(playList *PlayList) {
	client := GetDefaultM3u8DownloadContext().Client

	url := playList.Url

	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		baseUrl := GetDefaultM3u8DownloadContext().Options.BaseUrl

		if baseUrl == "" {
			uri, err := urlutil.Parse(GetDefaultM3u8DownloadContext().Options.Url)
			if err == nil {
				baseUrl = uri.Scheme + "://" + uri.Host
			}
		}
		if baseUrl == "" {
			log.Println("Can not parse base url")
			playList.WorkerId = ""
			return
		}
		url, _ = urlutil.JoinPath(baseUrl, url)
	}
	if url == "" {
		log.Println("Can not parse url")
		playList.WorkerId = ""
		return
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("download index %d error,url: %s", playList.Index, playList.Url)
		playList.WorkerId = ""
		return
	}

	//decrypt
	b := resp.Body
	if playList.Key != nil {

		method := util.EncryptMethodFromString(playList.Key.Method)
		if method != "" {
			b, err = util.Decrypt(b, method, util.DecryptKey{
				IV:  playList.Key.IVHex,
				Key: playList.Key.KeyHex,
			})

			if err != nil {
				log.Printf("decrypt error,index: %d", playList.Index)
				playList.WorkerId = ""
				return
			}

		}

	}

	dir := GetDefaultM3u8DownloadContext().OutDir
	cacheDir := path.Join(dir, Cache_Dir)
	tsDir := path.Join(cacheDir, Ts_Dir)

	tsFile := path.Join(tsDir, fmt.Sprintf("%d.ts", playList.Index))

	err = os.WriteFile(tsFile, b, os.ModePerm)

	if err != nil {
		log.Printf("write file error,file: %s", tsFile)
		playList.WorkerId = ""
		return
	}

	playList.Download = true

}
func (w *Worker) Wait() <-chan string {
	return w.notFull
}
