package conf

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	Token           string  `json:"TOKEN"`
	BatPath         string  `json:"BAT_PATH"`
	StoreFilePath   string  `json:"STORE_FILE_PATH"`
	GitServerPath   string  `json:"GIT_SERVER_PATH"`
	MemePath        string  `json:"MEME_PATH"`
	GoodMorningPath string  `json:"GOOD_MORING_PATH"`
	WeekendPath     string  `json:"WEEKEND_PATH"`
	MusicPath       string  `json:"MUSIC_PATH"`
	SaveFilePath    string  `json:"SAVE_FILE_PATH"`
	SaveFileName    string  `json:"SAVE_FILE_NAME"`
	BatName         string  `json:"BAT_NAME"`
	TestGroupId     int64   `json:"TEST_GROUP_ID"`
	GroupsId        []int64 `json:"GROUPS_ID"`
}

var (
	instance *Config
	once     sync.Once
)

func Init() {
	loadConfig()
}

func loadConfig() (*Config, error) {
	file, err := os.Open("tgbot.conf")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetInstance() *Config {
	var err error
	once.Do(func() {
		instance, err = loadConfig()
	})
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return instance
}
