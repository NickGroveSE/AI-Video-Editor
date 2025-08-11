package main

import (
	"fmt"
    "os"
    "path/filepath"
    
    "github.com/spf13/viper"
	"ai-video-editor/cmd"
)

func main() {
	// Set config file name and path
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting home directory:", err)
        os.Exit(1)
    }
    
    viper.SetConfigFile(filepath.Join(home, ".ai-editor.yaml"))
    viper.SetConfigType("yaml")
    
    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; ignore error
            fmt.Println("No config file found, using defaults")
        } else {
            // Config file was found but another error was produced
            fmt.Printf("Error reading config file: %v\n", err)
        }
    }

	cmd.Execute()
}

/* 

Test Commands

-- Basic usage
go run main.go process local/cooking_test_video.mp4 "find funny moments"

-- Advanced options
go run main.go process local/cooking_test_.mp4 "educational highlights" --duration 30s --output ./clips --quality high

-- Configuration
go run main.go config set-api-key openai sk-xxx
go run main.go config set whisper-model large

*/