package config

import (
	"github.com/bytedance/sonic"
	"log"
	"os"
	"path"
)

const (
	Version            = "0.1.0"
	Default_Workers    = 8
	Max_Workers        = 32
	Default_Queue_Size = 16
	Max_Queue_Size     = 64
	Http_Timeout       = 60

	Save_State_Durition = 15

	Config_Filename   = "config.json"
	Default_UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0"
)

var (
	User_Config_Dir string
	Config          *G3u8Config
)

func init() {

	home_dir, _ := os.UserHomeDir()

	User_Config_Dir = path.Join(home_dir, ".g3u8")

	if _, err := os.Stat(User_Config_Dir); os.IsNotExist(err) {
		os.MkdirAll(User_Config_Dir, os.ModePerm)
	}

	configFile := GetConfigFilePath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Creating config file,%s", configFile)
		Config = NewG3u8Config()
		b, err := sonic.Marshal(Config)

		if err != nil {
			log.Printf("Error: %v", err)
		}

		err = os.WriteFile(configFile, b, os.ModePerm)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	} else {
		log.Printf("Loading from  config file,%s", configFile)
		b, err := os.ReadFile(configFile)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		Config = &G3u8Config{}
		err = sonic.Unmarshal(b, Config)
		if err != nil {
			log.Printf("Error: %v", err)
			Config = NewG3u8Config()
		}

		if Config.Worker.MaxWorkers > Max_Workers {
			Config.Worker.MaxWorkers = Max_Workers
		}
		if Config.Worker.MaxWorkers < 1 {
			Config.Worker.MaxWorkers = Default_Workers
		}
		if Config.Worker.QueueSize > Max_Queue_Size {
			Config.Worker.QueueSize = Max_Queue_Size
		}
		if Config.Worker.QueueSize < 1 {
			Config.Worker.QueueSize = Default_Queue_Size
		}

		if Config.Worker.SaveStateDuration < 1 {
			Config.Worker.SaveStateDuration = Save_State_Durition
		}
	}

}

func GetConfigFilePath() string {
	return path.Join(User_Config_Dir, Config_Filename)
}
