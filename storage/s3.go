package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Wrapper struct {
	inst *s3.S3
}

func NewS3Wrapper(endpoint string, access string, secret string) (sw *S3Wrapper, err error) {
	cred := credentials.NewStaticCredentials(access, secret, "")
	config := &aws.Config{
		Region:           aws.String("default"), // 可用区配置
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      cred,
		DisableSSL:       aws.Bool(true),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		sw = &S3Wrapper{s3.New(sess)}
	}
	return
}

// create bucket

func (sw *S3Wrapper) S3Upload(buck string, key string, data []byte) (err error) {
	input := &s3.PutObjectInput{
		Bucket: &buck,
		Key:    &key,
		Body:   bytes.NewReader(data),
	}
	_, err = sw.inst.PutObject(input)
	return
}

func (sw *S3Wrapper) S3Download(buck string, key string) (data []byte, err error) {
	//input := &s3.GetObjectInput{
	//	Bucket: &buck,
	//	Key:&key,
	//}
	//resp, err := sw.inst.GetObject(input)
	//if err != nil {
	//	//var buf bytes.Buffer
	//	//buf.ReadFrom(resp.Body)
	//	//data = buf.Bytes()
	//
	//}
	return
}
