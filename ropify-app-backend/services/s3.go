package services

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3Bytes(imageBytes []byte, destiny, filename string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	fileName := fmt.Sprintf("%s/%d-%s", destiny, time.Now().Unix(), filename)
	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imageBytes),
		ContentType: aws.String("image/png"), // Cambia si usas otro formato
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return result.Location, nil
}

// DeleteFromS3 elimina un objeto de S3 usando su clave
func DeleteFromS3(objectKey string) error {
    // Crear sesión de AWS
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
    })
    if err != nil {
        return fmt.Errorf("failed to create AWS session: %v", err)
    }
    
    // Crear servicio S3
    svc := s3.New(sess)
    
    // Eliminar el objeto
    _, err = svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
        Key:    aws.String(objectKey),
    })
    if err != nil {
        return fmt.Errorf("failed to delete S3 object: %v", err)
    }
    
    // Esperar hasta que el objeto esté realmente eliminado
    err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
        Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
        Key:    aws.String(objectKey),
    })
    if err != nil {
        return fmt.Errorf("error waiting for object to be deleted: %v", err)
    }
    
    return nil
}

// ExtractS3KeyFromURL extrae la clave de un objeto de la URL de S3
func ExtractS3KeyFromURL(url string) (string, error) {
    // Las URLs de S3 suelen tener este formato:
    // https://bucket-ropify.s3.us-east-2.amazonaws.com/garments/users/userid/filename.jpg
    
    parts := strings.Split(url, ".amazonaws.com/")
    if len(parts) != 2 {
        return "", fmt.Errorf("invalid S3 URL format")
    }
    
    return parts[1], nil
}