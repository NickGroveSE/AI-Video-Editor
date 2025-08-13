package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"ai-video-editor/processing/video"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputDir    string
	clipDuration string
	maxClips     int
	quality      string
	skipAudio    bool
)

var processCmd = &cobra.Command{
	Use:   "process [video file] [prompt]",
	Short: "Process a video file and extract clips based on your prompt",
	Long: `Process analyzes your video file using AI and extracts relevant short clips
based on your prompt. It performs speech-to-text transcription, content analysis,
and generates captioned clips ready for social media or other uses.`,
	Args: cobra.ExactArgs(2),
	Example: `  # Extract funny moments as 30-second clips
  ai-editor process video.mp4 "find funny moments" --duration 30s
  
  # Get educational highlights with high quality
  ai-editor process lecture.mp4 "key learning points" --quality high --max-clips 5
  
  # Process to specific output directory
  ai-editor process presentation.mp4 "important quotes" --output ./clips`,
	RunE: runProcess,
}

func init() {
	rootCmd.AddCommand(processCmd)

	// Command-specific flags
	processCmd.Flags().StringVarP(&outputDir, "output", "o", "./clips", "output directory for generated clips")
	processCmd.Flags().StringVarP(&clipDuration, "duration", "d", "30s", "target duration for clips (e.g., 15s, 1m)")
	processCmd.Flags().IntVarP(&maxClips, "max-clips", "m", 10, "maximum number of clips to generate")
	processCmd.Flags().StringVarP(&quality, "quality", "", "medium", "output quality (low, medium, high)")
	processCmd.Flags().BoolVar(&skipAudio, "skip-audio", false, "skip audio processing and use video only")

	// Bind flags to viper for config file support
	viper.BindPFlag("output", processCmd.Flags().Lookup("output"))
	viper.BindPFlag("duration", processCmd.Flags().Lookup("duration"))
	viper.BindPFlag("max-clips", processCmd.Flags().Lookup("max-clips"))
	viper.BindPFlag("quality", processCmd.Flags().Lookup("quality"))
}

func runProcess(cmd *cobra.Command, args []string) error {
	videoFile := args[0]
	prompt := args[1]

	// Validate input file
	if err := validateVideoFile(videoFile); err != nil {
		return fmt.Errorf("invalid video file: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Display processing info
	if !viper.GetBool("quiet") {
		fmt.Printf("🎬 AI Video Editor\n")
		fmt.Printf("Input: %s\n", videoFile)
		fmt.Printf("Prompt: %s\n", prompt)
		fmt.Printf("Output: %s\n", outputDir)
		fmt.Printf("Duration: %s\n", clipDuration)
		fmt.Printf("Max clips: %d\n", maxClips)
		fmt.Printf("Quality: %s\n", quality)
		fmt.Println()
	}

	// TODO: Implement the actual processing pipeline
	fmt.Println("🔄 Starting video processing...")

	// Placeholder for processing steps
	// steps := []string{
	// 	"Analyzing video metadata",
	// 	"Extracting audio track",
	// 	"Performing speech-to-text transcription",
	// 	"Running AI content analysis",
	// 	"Identifying clip segments",
	// 	"Extracting video clips",
	// 	"Generating captions",
	// 	"Finalizing output files",
	// }

	// for i, step := range steps {
	// 	if !viper.GetBool("quiet") {
	// 		fmt.Printf("📋 Step %d/%d: %s...\n", i+1, len(steps), step)
	// 	}

	// 	// Simulate processing time
	// 	time.Sleep(500 * time.Millisecond)

	// 	if viper.GetBool("verbose") {
	// 		fmt.Printf("   ✅ %s completed\n", step)
	// 	}

	// 	if i == 0 {
	video.Analyze(videoFile)
	// 	}
	// }

	// if !viper.GetBool("quiet") {
	// 	fmt.Printf("\n🎉 Processing complete! Generated clips saved to: %s\n", outputDir)
	// }

	return nil
}

func validateVideoFile(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv"}

	for _, validExt := range validExts {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("unsupported file format: %s (supported: %s)",
		ext, strings.Join(validExts, ", "))
}
