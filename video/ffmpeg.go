package video

import (
	"fmt"
	"log"

	"github.com/xfrr/goffmpeg/transcoder"
)

func Analyze(filePath string) {

	fmt.Printf("Analyzing video: %s\n\n", filePath)

	// Create transcoder instance
	trans := new(transcoder.Transcoder)

	// Initialize with input file (no output needed for info)
	err := trans.Initialize(filePath, "")
	if err != nil {
		log.Fatalf("Error initializing transcoder: %v", err)
	}

	// Get the media file object
	mediaFile := trans.MediaFile()
	if mediaFile == nil {
		log.Fatal("Could not get media file information")
	}

	fmt.Println("✅ Successfully loaded video file!")

	// Try to get available information using the wrapper
	// Note: Different versions of goffmpeg may have different available methods
	fmt.Printf("Input file: %s\n", filePath)

	// Check if we can get basic file info
	fmt.Println(mediaFile.Duration())

	// The wrapper is working if we get here without errors
	fmt.Println("🎯 goffmpeg wrapper is functioning correctly!")
	fmt.Println("\nNote: For detailed metadata, the wrapper may need additional")
	fmt.Println("configuration or different methods depending on the version.")
}
