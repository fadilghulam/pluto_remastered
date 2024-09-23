package controllers

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"pluto_remastered/helpers"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func isImageMimeType(mimeType string) bool {
	imageMimeTypes := []string{"image/jpg", "image/jpeg", "image/jpe", "image/png", "image/jif", "image/jfif"}
	for _, t := range imageMimeTypes {
		if strings.EqualFold(mimeType, t) {
			return true
		}
	}
	return false
}

func initializeS3() (*s3.S3, string, error) {
	s3Region := os.Getenv("AWS_REGION")
	s3Bucket := os.Getenv("AWS_BUCKET")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if s3Region == "" || s3Bucket == "" || awsAccessKey == "" || awsSecretKey == "" {
		return nil, "", fmt.Errorf("AWS configuration is incomplete")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)
	return s3Client, s3Bucket, nil
}

func generateUniqueFilename(originalFilename string) string {
	extension := path.Ext(originalFilename)
	filename := uuid.New().String()
	return filename + extension
}

func uploadFileToS3(fileHeader *multipart.FileHeader, s3Client *s3.S3, bucketName string, directory string, fileName string) (string, string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	fileHeaderTemp := fileHeader.Header
	contentType := fileHeaderTemp.Get("Content-Type")

	if !isImageMimeType(contentType) {
		return "", "", "", fmt.Errorf("file is not a valid image")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", "", "", err
	}

	// Generate a unique filename

	var key, uniqueFilename string

	if fileName == "" {
		uniqueFilename = generateUniqueFilename(fileHeader.Filename)
		key = path.Join(directory, uniqueFilename)
	} else {
		fileName = fileName + path.Ext(fileHeader.Filename)
		uniqueFilename = fileName
		key = path.Join(directory, fileName)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          bytes.NewReader(buf.Bytes()),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(http.DetectContentType(buf.Bytes())),
	})
	if err != nil {
		return "", "", "", err
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com%s", bucketName, aws.StringValue(s3Client.Config.Region), key)
	return fileURL, uniqueFilename, helpers.TrimLeftChar(key), nil
}

func deleteFileFromS3(s3Client *s3.S3, bucketName string, fileKey string) error {
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return fmt.Errorf("error occurred while waiting for object to be deleted: %v", err)
	}

	return nil
}

func DoUpload(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	c.Accepts("multipart/form-data")

	s3Client, s3Bucket, err := initializeS3()
	if err != nil {
		log.Fatalf("Error initializing S3: %v", err)
	}

	file, err := c.FormFile("file")
	if err == nil {
		directory := c.FormValue("path")

		fileName := file.Filename

		fileURL, fileName, fileKey, err := uploadFileToS3(file, s3Client, s3Bucket, directory, fileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to upload file to S3",
				"success": false,
				"error":   err.Error(),
			})
		}

		fmt.Println(fileURL, fileName, fileKey)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "File berhasil diunggah",
			"success": true,
		})
	}

	// Check for multiple file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("Error reading form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).SendString("No files to upload")
	}

	directory := c.FormValue("path")

	for _, file := range files {

		fileName := file.Filename

		fileURL, fileName, fileKey, err := uploadFileToS3(file, s3Client, s3Bucket, directory, fileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to upload file to S3",
				"success": false,
				"error":   err.Error(),
			})
		}

		fmt.Println(fileURL, fileName, fileKey)

		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"message":  "File uploaded successfully",
		// 	"success":  true,
		// 	"fileURL":  fileURL,
		// 	"fileName": fileName,
		// 	"fileKey":  fileKey,
		// })
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File berhasil diunggah",
		"success": true,
	})
}

func DoDelete(c *fiber.Ctx) error {

	s3Client, s3Bucket, err := initializeS3()
	if err != nil {
		log.Fatalf("Error initializing S3: %v", err)
	}

	fileKey := c.FormValue("fileKey")
	if fileKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "fileKey query parameter is required",
			"success": false,
		})
	}

	err = deleteFileFromS3(s3Client, s3Bucket, fileKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete file from S3",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted successfully",
		"success": true,
	})
}
