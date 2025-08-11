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
ai-editor process video.mp4 "find funny moments"

-- Advanced options
ai-editor process video.mp4 "educational highlights" --duration 30s --output ./clips --quality high

-- Configuration
ai-editor config set-api-key openai sk-xxx
ai-editor config set whisper-model large

*/