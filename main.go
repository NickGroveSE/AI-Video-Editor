package main

import (
	"ai-video-editor/cmd"
)

func main() {
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