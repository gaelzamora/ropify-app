package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3(file *multipart.FileHeader, destiny string) (string, error) {
	// Crea una sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")), // Región configurada en tus variables de entorno
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Abre el archivo
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Genera un nombre único para el archivo
	fileName := fmt.Sprintf("%s/%d-%s", destiny, time.Now().Unix(), file.Filename)

	// Crea un uploader de S3
	uploader := s3manager.NewUploader(sess)

	fmt.Println("Bucket Name:", os.Getenv("AWS_BUCKET_NAME"))
	fmt.Println("AWS Region:", os.Getenv("AWS_REGION"))

	// Sube el archivo a S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")), // Nombre del bucket
		Key:    aws.String(fileName),                     // Ruta del archivo en el bucket
		Body:   src,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Devuelve el URL público del archivo
	return result.Location, nil
}
