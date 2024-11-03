package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App   AppConfig
		Redis RedisConfig
		Log   LogConfig
	}

	// 應用程式設定
	AppConfig struct {
		Port             int
		PidFile          string
		SearchGameCode   []string
		ScheduleCleanRTP string
	}

	// Redis設定
	RedisConfig struct {
		ReadTimeout        int
		WriteTimeout       int
		MaxRetries         int
		DialTimeout        int
		PoolSize           int
		PoolTimeout        int
		IdleTimeout        int
		IdleCheckFrequency int
		Hosts              []string
	}

	// Log設定
	LogConfig struct {
		Enable       bool   // 是否啟用
		Level        string // log level
		FileSizeMega int    // log rotate 的檔案大小 (MB)
		FileCount    int    // log 檔的保留數量
		KeepDays     int    // log 檔名日期的保留天數
		Path         string // log 路徑; 若為空字串, 則不輸出到檔案
	}
)

// LoadConfig loads the configuration from the specified filename.
func LoadConfig(filename string) (Config, error) {
	// Create a new Viper instance.
	v := viper.New()

	// Set the configuration file name, path, and environment variable settings.
	v.SetConfigName(fmt.Sprintf("./%s", filename))
	v.AddConfigPath(".")
	v.SetConfigType("json")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file.
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return Config{}, err
	}

	// Unmarshal the configuration into the Config struct.
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return Config{}, err
	}

	return config, nil
}

// LoadConfigPath loads the configuration from the specified path.
func LoadConfigPath(path string) (Config, error) {
	// Create a new Viper instance.
	v := viper.New()

	// Set the configuration file name, path, and environment variable settings.
	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file.
	if err := v.ReadInConfig(); err != nil {
		// Handle the case where the configuration file is not found.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	// Parse the configuration into the Config struct.
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}
