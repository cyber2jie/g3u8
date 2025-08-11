package download

import (
	"errors"
	"g3u8/m3u8"
	urlutil "net/url"
	"os"
	"strings"
)

func ParseFromUrl(url string) (*m3u8.MediaPlaylist, error) {
	resp, err := GetDefaultM3u8DownloadContext().Client.Get(url)
	if err != nil {
		return nil, err
	}
	return ParseFromString(string(resp.Body))
}

func ParseFromFile(path string) (*m3u8.MediaPlaylist, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseFromString(string(bytes))
}

func ParseFromString(m3u8String string) (*m3u8.MediaPlaylist, error) {
	playList, listType, err := m3u8.DecodeFrom(strings.NewReader(m3u8String), false)
	if err != nil {
		return nil, err
	}

	switch listType {
	case m3u8.MEDIA:
		return playList.(*m3u8.MediaPlaylist), nil
	}
	return nil, errors.New("Can't detect playlist type")
}

func ParseListFromUrl(url string) (*m3u8.MediaPlaylist, error) {

	resp, err := GetDefaultM3u8DownloadContext().Client.Get(url)

	if err != nil {
		return nil, err
	}

	playList, listType, err := m3u8.DecodeFrom(strings.NewReader(string(resp.Body)), false)
	if err != nil {
		return nil, err
	}

	switch listType {
	case m3u8.MASTER:
		return ParseFromMasterList(url, playList.(*m3u8.MasterPlaylist))
	case m3u8.MEDIA:
		return playList.(*m3u8.MediaPlaylist), nil
	}
	return nil, errors.New("Can not parse m3u8 list")
}
func ParseFromMasterList(base string, masterList *m3u8.MasterPlaylist) (*m3u8.MediaPlaylist, error) {

	if masterList != nil {
		for _, variant := range masterList.Variants {
			// only handle one  variant
			if variant.URI != "" {
				url := variant.URI
				if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
					baseUrl := GetDefaultM3u8DownloadContext().Options.BaseUrl

					if baseUrl == "" {

						uri, err := urlutil.Parse(base)
						if err != nil {
							return nil, err
						}
						baseUrl = uri.Scheme + "://" + uri.Host

					}
					_url, err := urlutil.JoinPath(baseUrl, url)
					if err != nil {
						return nil, err
					}
					url = _url
				}
				return ParseListFromUrl(url)
			}
		}
	}
	return nil, errors.New("Can not parse m3u8 list")
}
