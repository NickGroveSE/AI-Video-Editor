package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// AudioExtractor handles audio extraction from video files
type AudioExtractor struct {
	TempDir string
}

// NewAudioExtractor creates a new audio extractor
func NewAudioExtractor(tempDir string) *AudioExtractor {
	return &AudioExtractor{
		TempDir: tempDir,
	}
}

// ExtractAudio extracts audio from video and returns path to audio file
func (ae *AudioExtractor) ExtractAudio(videoPath string) (string, error) {
	// Create unique temporary file name
	timestamp := time.Now().Unix()
	audioFileName := fmt.Sprintf("audio_%d.wav", timestamp)
	audioPath := filepath.Join(ae.TempDir, audioFileName)

	// Ensure temp directory exists
	if err := os.MkdirAll(ae.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract audio using ffmpeg-go
	err := ffmpeg.Input(videoPath).
		Output(audioPath, ffmpeg.KwArgs{
			"vn":       "",           // No video
			"acodec":   "pcm_s16le",  // 16-bit PCM codec
			"ar":       16000,        // 16kHz sample rate
			"ac":       1,            // Mono audio
			"f":        "wav",        // WAV format
		}).
		OverWriteOutput(). // Overwrite if file exists
		Silent(true).      // Suppress ffmpeg output
		Run()

	if err != nil {
		return "", fmt.Errorf("failed to extract audio: %w", err)
	}

	return audioPath, nil
}

// ExtractAudioChunk extracts a specific time segment from video
func (ae *AudioExtractor) ExtractAudioChunk(videoPath string, startTime, duration float64) (string, error) {
	// Create unique temporary file name
	timestamp := time.Now().UnixNano()
	audioFileName := fmt.Sprintf("audio_chunk_%d.wav", timestamp)
	audioPath := filepath.Join(ae.TempDir, audioFileName)

	// Ensure temp directory exists
	if err := os.MkdirAll(ae.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract specific audio segment
	err := ffmpeg.Input(videoPath, ffmpeg.KwArgs{
		"ss": startTime, // Start time in seconds
		"t":  duration,  // Duration in seconds
	}).
		Output(audioPath, ffmpeg.KwArgs{
			"vn":     "",           // No video
			"acodec": "pcm_s16le",  // 16-bit PCM codec
			"ar":     16000,        // 16kHz sample rate
			"ac":     1,            // Mono audio
			"f":      "wav",        // WAV format
		}).
		OverWriteOutput().
		Silent(true).
		Run()

	if err != nil {
		return "", fmt.Errorf("failed to extract audio chunk: %w", err)
	}

	return audioPath, nil
}

// GetAudioDuration gets the duration of the audio in the video
func (ae *AudioExtractor) GetAudioDuration(videoPath string) (float64, error) {
	// Use ffprobe to get duration
	data, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, fmt.Errorf("failed to probe video: %w", err)
	}

	// Parse duration from probe data
	// Note: You'll need to parse the JSON response from ffmpeg.Probe
	// This is a simplified version - you'd want proper JSON parsing
	return 0, fmt.Errorf("duration parsing not implemented")
}

// CleanupAudioFile removes the temporary audio file
func (ae *AudioExtractor) CleanupAudioFile(audioPath string) error {
	return os.Remove(audioPath)
}

// ProcessVideoForHuggingFace extracts audio and prepares it for API
func (ae *AudioExtractor) ProcessVideoForHuggingFace(videoPath string) ([]byte, string, error) {
	// Extract audio to temporary file
	audioPath, err := ae.ExtractAudio(videoPath)
	if err != nil {
		return nil, "", fmt.Errorf("audio extraction failed: %w", err)
	}

	// Read audio file into bytes for API
	audioBytes, err := os.ReadFile(audioPath)
	if err != nil {
		ae.CleanupAudioFile(audioPath) // Clean up on error
		return nil, "", fmt.Errorf("failed to read audio file: %w", err)
	}

	return audioBytes, audioPath, nil
}

// ProcessLongVideoInChunks processes video in chunks for long videos
func (ae *AudioExtractor) ProcessLongVideoInChunks(videoPath string, chunkDuration float64) ([]AudioChunk, error) {
	// First get total duration (you'd implement proper duration detection)
	// totalDuration, err := ae.GetAudioDuration(videoPath)
	// For now, using placeholder
	totalDuration := 300.0 // 5 minutes example

	var chunks []AudioChunk
	currentTime := 0.0

	for currentTime < totalDuration {
		remainingTime := totalDuration - currentTime
		actualDuration := chunkDuration
		if remainingTime < chunkDuration {
			actualDuration = remainingTime
		}

		// Extract this chunk
		audioPath, err := ae.ExtractAudioChunk(videoPath, currentTime, actualDuration)
		if err != nil {
			// Clean up any successful chunks on error
			for _, chunk := range chunks {
				ae.CleanupAudioFile(chunk.AudioPath)
			}
			return nil, fmt.Errorf("failed to extract chunk at %f: %w", currentTime, err)
		}

		// Read chunk into bytes
		audioBytes, err := os.ReadFile(audioPath)
		if err != nil {
			ae.CleanupAudioFile(audioPath)
			return nil, fmt.Errorf("failed to read chunk file: %w", err)
		}

		chunks = append(chunks, AudioChunk{
			AudioPath:   audioPath,
			AudioBytes:  audioBytes,
			StartTime:   currentTime,
			Duration:    actualDuration,
			ChunkIndex:  len(chunks),
		})

		currentTime += chunkDuration
	}

	return chunks, nil
}

// AudioChunk represents a segment of audio from the video
type AudioChunk struct {
	AudioPath   string
	AudioBytes  []byte
	StartTime   float64
	Duration    float64
	ChunkIndex  int
}

// Example usage
func main() {
	extractor := NewAudioExtractor("./temp")
	
	// Simple extraction
	audioBytes, audioPath, err := extractor.ProcessVideoForHuggingFace("./video.mp4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer extractor.CleanupAudioFile(audioPath)

	fmt.Printf("Extracted %d bytes of audio to %s\n", len(audioBytes), audioPath)
	
	// Now you can send audioBytes to your Hugging Face API client
	// transcription := yourHuggingFaceClient.Transcribe(audioBytes)
}