package video

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// VideoExtractor handles video extraction from video files
type VideoExtractor struct {
	TempDir string
}

// NewVideoExtractor creates a new video extractor
func NewVideoExtractor(tempDir string) *VideoExtractor {
	return &VideoExtractor{
		TempDir: tempDir,
	}
}

// ExtractVideo extracts video from video and returns path to video file
func (ae *VideoExtractor) ExtractVideoPath(inputPath string) (string, error) {
	// Create unique temporary file name
	timestamp := time.Now().Unix()
	videoFileName := fmt.Sprintf("video_%d.mp4", timestamp)
	videoPath := filepath.Join(ae.TempDir, videoFileName)

	// Ensure temp directory exists
	if err := os.MkdirAll(ae.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract video using ffmpeg-go
	err := ffmpeg.Input(inputPath).

		Output(videoPath, ffmpeg.KwArgs{"an": ""}).
		OverWriteOutput().
		Run()

	if err != nil {
		return "", fmt.Errorf("failed to extract video: %w", err)
	}

	return videoPath, nil
}

/*
// ExtractvideoChunk extracts a specific time segment from video
func (ae *videoExtractor) ExtractvideoChunk(videoPath string, startTime, duration float64) (string, error) {
	// Create unique temporary file name
	timestamp := time.Now().UnixNano()
	videoFileName := fmt.Sprintf("video_chunk_%d.wav", timestamp)
	videoPath := filepath.Join(ae.TempDir, videoFileName)

	// Ensure temp directory exists
	if err := os.MkdirAll(ae.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract specific video segment
	err := ffmpeg.Input(videoPath, ffmpeg.KwArgs{
		"ss": startTime, // Start time in seconds
		"t":  duration,  // Duration in seconds
	}).
		Output(videoPath, ffmpeg.KwArgs{
			"vn":     "",           // No video
			"acodec": "pcm_s16le",  // 16-bit PCM codec
			"ar":     16000,        // 16kHz sample rate
			"ac":     1,            // Mono video
			"f":      "wav",        // WAV format
		}).
		OverWriteOutput().
		Silent(true).
		Run()

	if err != nil {
		return "", fmt.Errorf("failed to extract video chunk: %w", err)
	}

	return videoPath, nil
}
*/

/*
// GetvideoDuration gets the duration of the video in the video
func (ae *videoExtractor) GetvideoDuration(videoPath string) (float64, error) {
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
*/

// CleanupvideoFile removes the temporary video file
func (ae *VideoExtractor) CleanupVideoFile(videoPath string) error {
	return os.Remove(videoPath)
}


// ProcessVideoForHuggingFace extracts video and prepares it for API
func (ae *VideoExtractor) ExtractVideoBytes(videoPath string) ([]byte, string, error) {
	// Extract video to temporary file
	videoPath, err := ae.ExtractVideoPath(videoPath)
	if err != nil {
		return nil, "", fmt.Errorf("video extraction failed: %w", err)
	}

	// Read video file into bytes for API
	videoBytes, err := os.ReadFile(videoPath)
	if err != nil {
		ae.CleanupVideoFile(videoPath) // Clean up on error
		return nil, "", fmt.Errorf("failed to read video file: %w", err)
	}

	return videoBytes, videoPath, nil
}

/*
// ProcessLongVideoInChunks processes video in chunks for long videos
func (ae *videoExtractor) ProcessLongVideoInChunks(videoPath string, chunkDuration float64) ([]videoChunk, error) {
	// First get total duration (you'd implement proper duration detection)
	// totalDuration, err := ae.GetvideoDuration(videoPath)
	// For now, using placeholder
	totalDuration := 300.0 // 5 minutes example

	var chunks []videoChunk
	currentTime := 0.0

	for currentTime < totalDuration {
		remainingTime := totalDuration - currentTime
		actualDuration := chunkDuration
		if remainingTime < chunkDuration {
			actualDuration = remainingTime
		}

		// Extract this chunk
		videoPath, err := ae.ExtractvideoChunk(videoPath, currentTime, actualDuration)
		if err != nil {
			// Clean up any successful chunks on error
			for _, chunk := range chunks {
				ae.CleanupvideoFile(chunk.videoPath)
			}
			return nil, fmt.Errorf("failed to extract chunk at %f: %w", currentTime, err)
		}

		// Read chunk into bytes
		videoBytes, err := os.ReadFile(videoPath)
		if err != nil {
			ae.CleanupvideoFile(videoPath)
			return nil, fmt.Errorf("failed to read chunk file: %w", err)
		}

		chunks = append(chunks, videoChunk{
			videoPath:   videoPath,
			videoBytes:  videoBytes,
			StartTime:   currentTime,
			Duration:    actualDuration,
			ChunkIndex:  len(chunks),
		})

		currentTime += chunkDuration
	}

	return chunks, nil
}
*/

/*
// videoChunk represents a segment of video from the video
type videoChunk struct {
	videoPath   string
	videoBytes  []byte
	StartTime   float64
	Duration    float64
	ChunkIndex  int
}
*/

// Example usage
func ExtractVideo(inputPath string) {
	extractor := NewVideoExtractor("./temp")
	
	// Simple extraction
	videoBytes, videoPath, err := extractor.ExtractVideoBytes(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer extractor.CleanupVideoFile(videoPath)

	fmt.Printf("Extracted %d bytes of video to %s\n", len(videoBytes), videoPath)
	
	// Now you can send videoBytes to your Hugging Face API client
	// transcription := yourHuggingFaceClient.Transcribe(videoBytes)
}