package aws

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/pulsetestinfo/v1"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientTestSuite 是 aws 通信测试
type ClientTestSuite struct {
	suite.Suite
	client *Client
}

// SetupSuite 初始化测试
func (suite *ClientTestSuite) SetupSuite() {
	_ = godotenv.Load("./testdata/local.svc-analysis.env")
	suite.client, _ = NewClient(
		BucketName(os.Getenv("X_AWS_BUCKET_NAME")),
		AccessKeyID(os.Getenv("X_AWS_ACCESS_KEY")),
		SecretKey(os.Getenv("X_AWS_SECRET_KEY")),
		Region(os.Getenv("X_AWS_REGION")),
		PulseTestRawDataEnvironmentS3KeyPrefix(os.Getenv("X_WAVE_DATA_KEY_PREFIX")),
		PulseTestRawDataS3KeyPrefix(os.Getenv("X_PULSE_TEST_RAW_DATA_S3_KEY_PREFIX")),
	)
}

// TestUploadFile 测试上传一段文本
func (suite *ClientTestSuite) TestUploadFile() {
	const key = "text-2018-02"
	t := suite.T()
	buf := bytes.NewBufferString("")
	buf.WriteString("1\n2\n3\n")
	_, err := suite.client.UploadFile(bytes.NewReader(buf.Bytes()), path.Join(suite.client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key), s3.ObjectCannedACLPrivate)
	assert.NoError(t, err)
}

// TestDownloadFile 测试下载文本
func (suite *ClientTestSuite) TestDownloadFile() {
	const key = "text-2018-02"
	content := uuid.New().String()
	t := suite.T()
	buf := bytes.NewBufferString(content)
	_, errUploadFile := suite.client.UploadFile(bytes.NewReader(buf.Bytes()), path.Join(suite.client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key), s3.BucketCannedACLPrivate)
	assert.NoError(t, errUploadFile)
	output, err := suite.client.DownloadFile(path.Join(suite.client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key))
	assert.NoError(t, err)
	assert.Equal(t, content, string(output))
}

// TestUploadWaveData 测试上传波形数据
func (suite *ClientTestSuite) TestUploadWaveData() {
	t := suite.T()
	const recordID = 999999
	const testWaveDataLength = 100
	waveData := make([]byte, 0)
	for i := 0; i < testWaveDataLength; i++ {
		waveData = append(waveData, byte(i))
	}
	_, err := suite.client.Upload(proto.PulseTestRawInfo{
		Spec:     uint32(1),
		RecordId: uint32(recordID),
		Payloads: waveData,
	})
	assert.NoError(t, err)
}

// TestDownloadWaveData 测试下载波形数据
func (suite *ClientTestSuite) TestDownloadWaveData() {
	t := suite.T()
	dir := initReportDir()
	key := "spec-v2/229967.pbd"
	fullKey := path.Join(suite.client.options.PulseTestRawDataEnvironmentS3KeyPrefix, key)
	dataDownload, err := suite.client.DownloadFile(fullKey)
	filename := fmt.Sprintf("%s/229967.pbd", dir)
	errWriteFile := ioutil.WriteFile(filename, dataDownload, 0655)
	assert.NoError(t, errWriteFile)
	assert.NoError(t, err)
	assert.NotEqual(t, "", string(dataDownload))
	assert.NotEqual(t, filename, string(dataDownload))
}

// initReportDir 初始化文件夹
func initReportDir() string {
	dir := "./report"
	exist := PathExists(dir)
	if exist {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return dir
}

// PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
