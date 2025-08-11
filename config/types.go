package config

type G3u8Config struct {
	Version string       `json:"version,omitempty"`
	Worker  WorkerConfig `json:"worker"`
	Proxy   ProxyConfig  `json:"proxy"`
	Http    HttpConfig   `json:"http"`
}

func NewG3u8Config() *G3u8Config {
	return &G3u8Config{
		Version: Version,
		Worker: WorkerConfig{
			MaxWorkers:        Default_Workers,
			QueueSize:         Default_Queue_Size,
			SaveStateDuration: Save_State_Durition,
		},
		Proxy: ProxyConfig{
			Enable: false,
		},
		Http: HttpConfig{
			Timeout: Http_Timeout,
			Headers: []HttpHeader{
				HttpHeader{
					Name:  "User-Agent",
					Value: Default_UserAgent,
				},
			},
		},
	}
}

type WorkerConfig struct {
	MaxWorkers        int `json:"max_workers"`
	QueueSize         int `json:"queue_size"`
	SaveStateDuration int `json:"save_state_duration"`
}

type ProxyConfig struct {
	Proxy  *string `json:"proxy"`
	Enable bool    `json:"enable"`
}

type HttpConfig struct {
	Timeout int          `json:"timeout"` //second
	Headers []HttpHeader `json:"headers"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
