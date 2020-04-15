package aws

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	protobuf "github.com/golang/protobuf/proto"
	pulsetestinfopb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/pulsetestinfo/v1"
)

// PulseTestRawDataS3Client 可以上传和下载波形数据
type PulseTestRawDataS3Client interface {
	Upload(data pulsetestinfopb.PulseTestRawInfo) (*s3.PutObjectOutput, error)
	Download(key string) (*pulsetestinfopb.PulseTestRawInfo, error)
}

// Client 与 aws 通信
type Client struct {
	options *Options
	sess    *session.Session
	s3      *s3.S3
}

// NewClient 返回一个新的 aws 连接
func NewClient(opts ...Option) (*Client, error) {
	client := new(Client)
	client.options = newOptions(opts...)
	creds := credentials.NewStaticCredentials(client.options.AccessKeyID, client.options.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(client.options.Region),
	})
	if err != nil {
		return nil, err
	}
	client.sess = sess
	client.s3 = s3.New(sess)
	return client, nil
}

// UploadFile 存储一个输入流到指定路径
func (client *Client) UploadFile(readerSeeker io.ReadSeeker, key string, acl string) (*s3.PutObjectOutput, error) {
	svc := client.s3
	return svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(client.options.BucketName),
		Key:    aws.String(key),
		Body:   readerSeeker,
		ACL:    aws.String(acl),
	})
}

// DownloadFile 从存储桶读取文件
func (client *Client) DownloadFile(key string) ([]byte, error) {
	svc := client.s3
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(client.options.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close() // nolint: errcheck
	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Upload 存储波形数据
func (client *Client) Upload(data pulsetestinfopb.PulseTestRawInfo) (*s3.PutObjectOutput, error) {

	key := GenerateS3Key(client.options.PulseTestRawDataS3KeyPrefix, int(data.RecordId))
	datas, errMarshal := protobuf.Marshal(&data)
	if errMarshal != nil {
		return nil, errMarshal
	}
	return client.UploadFile(readWaveDataSendToAWS(datas), path.Join(client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key), s3.ObjectCannedACLPrivate)
}

// Download 下载波形数据
func (client *Client) Download(key string) (*pulsetestinfopb.PulseTestRawInfo, error) {
	fullKey := path.Join(client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key)
	downloadFile, errDownloadFile := client.DownloadFile(fullKey)
	if errDownloadFile != nil {
		return nil, errDownloadFile
	}
	return ParsePulseTestData(downloadFile, key)
}

// ParsePulseTestData 解析脉冲测试数据
func ParsePulseTestData(downloadFile []byte, key string) (*pulsetestinfopb.PulseTestRawInfo, error) {
	isSpecV2 := strings.HasSuffix(key, ".pbd")
	if isSpecV2 {
		return ParsePulseTestDataV2(downloadFile)
	}
	return ParsePulseTestDataV1(downloadFile)
}

// ParsePulseTestDataV2 解析V2波形数据
func ParsePulseTestDataV2(waveDataFile []byte) (*pulsetestinfopb.PulseTestRawInfo, error) {
	var data pulsetestinfopb.PulseTestRawInfo
	err := protobuf.Unmarshal(waveDataFile, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// ParsePulseTestDataV1 解析波形数据
func ParsePulseTestDataV1(data []byte) (*pulsetestinfopb.PulseTestRawInfo, error) {
	return &pulsetestinfopb.PulseTestRawInfo{
		Spec:     1,
		Payloads: data,
	}, nil
}

// readWaveDataSendToAWS 将波形数据格式化为可读文本流
func readWaveDataSendToAWS(data []byte) io.ReadSeeker {
	buf := bytes.NewBuffer(data)
	return bytes.NewReader(buf.Bytes())
}

// GenerateS3Key 生成S3的key
func GenerateS3Key(pulseTestRawDataS3KeyPrefix string, recordID int) string {
	t := time.Now().UTC()
	return path.Join(pulseTestRawDataS3KeyPrefix, fmt.Sprintf("/%d/%d/%d/%d.pbd", t.Year(), t.Month(), t.Day(), recordID))
}
