package termrecorder

import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/service/s3"
import "github.com/aws/aws-sdk-go/aws/session"
import "os"

type AwsUploader struct {
    region string
    bucket string
}

func (uploader *AwsUploader) Upload(
    user string,
    filename string,
    source *os.File) {

    sess := session.New(&aws.Config{
        Region: aws.String("ap-southeast-2"),
    })

    svc := s3.New(sess)

    params := &s3.PutObjectInput{
        Bucket: aws.String(uploader.bucket),
        Key: aws.String(user + "/" + filename),
        Body: source,
    }

    svc.PutObject(params)
}
