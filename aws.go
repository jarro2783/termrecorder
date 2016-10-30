package termrecorder

import "fmt"
import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/service/s3"
import "github.com/aws/aws-sdk-go/aws/session"
import "os"

type AwsUploader struct {
    region string
    bucket string
    root string
    subpath string
}

func  MakeAwsUploader(region, bucket, root, subpath string) *AwsUploader {
    aws := &AwsUploader{region, bucket, root, subpath}
    return aws
}

func (uploader *AwsUploader) Upload(
    user string,
    gameid string,
    filename string,
    source *os.File) {

    sess := session.New(&aws.Config{
        Region: aws.String(uploader.region),
    })

    svc := s3.New(sess)

    var key string
    if len(uploader.root) != 0 {
        key += uploader.root + "/"
    }

    key += user

    if len(uploader.subpath) != 0 {
        key += "/" + uploader.subpath
    }

    if len(gameid) != 0 {
        key += "/" + gameid
    }

    key += "/" +  filename

    params := &s3.PutObjectInput{
        Bucket: aws.String(uploader.bucket),
        Key: aws.String(key),
        Body: source,
    }

    output, err := svc.PutObject(params)

    if err != nil {
        fmt.Printf("Error uploading %s for user %s\n", filename, user)
        fmt.Printf(err.Error())
    } else {
        fmt.Printf("Put object %s with etag %s\n", filename, *output.ETag)
    }
}
