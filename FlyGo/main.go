package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/segmentio/kafka-go"
)

type FrameMessage struct {
	VideoID     string `json:"video_id"`
	FrameID     string `json:"frame_id"`
	FrameNumber int    `json:"frame_number"`
	S3Path      string `json:"s3_path"`
}

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

func connectDB() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return pgxpool.New(context.Background(), dbURL)
}

func insertVideo(db *pgxpool.Pool, videoID, name string, length, fps int) error {
	_, err := db.Exec(context.Background(),
		"INSERT INTO videos (video_id, name, length, fps) VALUES ($1, $2, $3, $4)",
		videoID, name, length, fps)
	return err
}

func insertFrame(db *pgxpool.Pool, frameID, videoID string, frameNumber int, s3Path, status string) error {
	_, err := db.Exec(context.Background(),
		"INSERT INTO frames (frame_id, video_id, frame_number, s3_path, status) VALUES ($1, $2, $3, $4, $5)",
		frameID, videoID, frameNumber, s3Path, status)
	return err
}

func sendKafkaMessage(broker, topic string, message interface{}) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return w.WriteMessages(context.Background(),
		kafka.Message{
			Value: msgBytes,
		},
	)
}

func main() {
	godotenv.Load()
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <video_path> <output_dir> <fps> <job_id>")
		os.Exit(1)
	}
	videoPath := os.Args[1]
	outputDir := os.Args[2]
	fps := 1
	fmt.Sscanf(os.Args[3], "%d", &fps)
	jobID := os.Args[4]

	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	videoID := jobID
	videoName := filepath.Base(videoPath)
	videoLength := 0

	if err := insertVideo(db, videoID, videoName, videoLength, fps); err != nil {
		log.Fatalf("Failed to insert video: %v", err)
	}
	fmt.Printf("Inserted video metadata: %s\n", videoID)

	if err := extractFrames(videoPath, outputDir, fps); err != nil {
		fmt.Printf("Error extracting frames: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Frame extraction complete.")

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "your-s3-bucket-name"
	}
	prefix := fmt.Sprintf("frames/%s", videoID)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
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
				frameID := uuid.New().String()
				var frameNumber int
				fmt.Sscanf(info.Name(), "frame_%d.jpg", &frameNumber)
				s3Path := fmt.Sprintf("s3://%s/%s", bucket, key)
				if err := insertFrame(db, frameID, videoID, frameNumber, s3Path, "pending"); err != nil {
					log.Printf("Failed to insert frame: %v", err)
				} else {
					fmt.Printf("Inserted frame metadata: %s\n", frameID)
					// Send Kafka message
					kafkaBroker := "localhost:9092"
					kafkaTopic := "frames-to-process"
					frameMsg := FrameMessage{
						VideoID:     videoID,
						FrameID:     frameID,
						FrameNumber: frameNumber,
						S3Path:      s3Path,
					}
					if err := sendKafkaMessage(kafkaBroker, kafkaTopic, frameMsg); err != nil {
						log.Printf("Failed to send Kafka message: %v", err)
					} else {
						fmt.Printf("Sent Kafka message for frame %d\n", frameNumber)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error uploading frames to S3: %v", err)
	}
	fmt.Println("All frames uploaded to S3 and metadata inserted.")
}

// go run main.go input/input2.mp4 output 1
