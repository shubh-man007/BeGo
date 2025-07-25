package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/joho/godotenv"
)

func extractFrames(videoPath, outputDir string, fps int) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPattern := filepath.Join(outputDir, "frame_%04d.jpg")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("fps=%d", fps), outputPattern)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running command: %v\n", cmd.Args)
	return cmd.Run()
}

func uploadFileToS3(ctx context.Context, client *s3.Client, bucket, key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload %s to S3: %w", filePath, err)
	}
	return nil
}

func uploadDirToS3(ctx context.Context, client *s3.Client, bucket, prefix, dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			key := filepath.Join(prefix, info.Name())
			fmt.Printf("Uploading %s to s3://%s/%s...\n", path, bucket, key)
			if err := uploadFileToS3(ctx, client, bucket, key, path); err != nil {
				log.Printf("Failed to upload %s: %v", path, err)
			} else {
				fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", path, bucket, key)
			}
		}
		return nil
	})
	return err
}

func main() {
	godotenv.Load()

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <video_path> <output_dir> <fps>")
		os.Exit(1)
	}

	videoPath := os.Args[1]
	outputDir := os.Args[2]
	fps := 1
	fmt.Sscanf(os.Args[3], "%d", &fps)

	if err := extractFrames(videoPath, outputDir, fps); err != nil {
		fmt.Printf("Error extracting frames: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Frame extraction complete.")

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "fly-brain-img01"
	}
	prefix := "frames"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	if err := uploadDirToS3(ctx, client, bucket, prefix, outputDir); err != nil {
		log.Fatalf("Error uploading frames to S3: %v", err)
	}
	fmt.Println("All frames uploaded to S3.")
}

// go run main.go input/input1.mp4 output 1
