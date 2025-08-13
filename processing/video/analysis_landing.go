package video

import (
	"fmt"
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Analyze(inputPath string) {

	// Set FFprobe path using environment variable (correct approach for u2takey/ffmpeg-go)
	os.Setenv("FFPROBE_PATH", "C:\\ffmpeg\\bin\\ffprobe.exe")
	
	// Also ensure PATH includes FFmpeg directory
	currentPath := os.Getenv("PATH")
	os.Setenv("PATH", "C:\\ffmpeg\\bin;"+currentPath)
	
	// Check if file path was provided as argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-video-file>")
	}

	filePath := os.Args[1]

	fmt.Printf("📹 Analyzing video: %s\n\n", filePath)

	// Use ffprobe to get video metadata
	_, err := ffmpeg.Probe(inputPath)
	if err != nil {
		log.Fatalf("❌ Error probing video file: %v", err)
	}

	ExtractAudio(inputPath)
	ExtractVideo(inputPath)

	/*
	fmt.Println("✅ Successfully probed video file!")

	// Parse the JSON response
	type ProbeData struct {
		Format struct {
			Filename       string `json:"filename"`
			FormatName     string `json:"format_name"`
			FormatLongName string `json:"format_long_name"`
			Duration       string `json:"duration"`
			Size           string `json:"size"`
			BitRate        string `json:"bit_rate"`
		} `json:"format"`
		Streams []struct {
			Index              int    `json:"index"`
			CodecName          string `json:"codec_name"`
			CodecLongName      string `json:"codec_long_name"`
			CodecType          string `json:"codec_type"`
			Width              int    `json:"width"`
			Height             int    `json:"height"`
			PixelFormat        string `json:"pix_fmt"`
			Duration           string `json:"duration"`
			BitRate            string `json:"bit_rate"`
			SampleRate         string `json:"sample_rate"`
			Channels           int    `json:"channels"`
		} `json:"streams"`
	}

	var probeData ProbeData
	err = json.Unmarshal([]byte(data), &probeData)
	if err != nil {
		log.Fatalf("❌ Error parsing probe data: %v", err)
	}

	// Print basic file information
	fmt.Println("📊 File Information:")
	fmt.Printf("   Filename: %s\n", probeData.Format.Filename)
	fmt.Printf("   Format: %s (%s)\n", probeData.Format.FormatName, probeData.Format.FormatLongName)
	fmt.Printf("   Duration: %s seconds\n", probeData.Format.Duration)
	fmt.Printf("   Size: %s bytes\n", probeData.Format.Size)
	fmt.Printf("   Bit Rate: %s bps\n", probeData.Format.BitRate)

	// Print stream information
	fmt.Printf("\n🎬 Streams (%d total):\n", len(probeData.Streams))
	
	for i, stream := range probeData.Streams {
		fmt.Printf("\n   Stream %d (%s):\n", i, stream.CodecType)
		fmt.Printf("      Codec: %s (%s)\n", stream.CodecName, stream.CodecLongName)
		
		if stream.CodecType == "video" {
			fmt.Printf("      Resolution: %dx%d\n", stream.Width, stream.Height)
			if stream.PixelFormat != "" {
				fmt.Printf("      Pixel Format: %s\n", stream.PixelFormat)
			}
		}
		
		if stream.CodecType == "audio" {
			fmt.Printf("      Sample Rate: %s Hz\n", stream.SampleRate)
			fmt.Printf("      Channels: %d\n", stream.Channels)
		}
		
		if stream.Duration != "" {
			fmt.Printf("      Duration: %s seconds\n", stream.Duration)
		}
		if stream.BitRate != "" {
			fmt.Printf("      Bit Rate: %s bps\n", stream.BitRate)
		}
	}

	fmt.Println("\n🎉 Analysis complete!")
	fmt.Println("💡 u2takey/ffmpeg-go wrapper is working correctly!")
	*/
}