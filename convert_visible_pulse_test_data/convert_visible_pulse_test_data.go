package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	awsTool "github.com/jinmukeji/jiujiantang-services/analysis/aws"
	handler "github.com/jinmukeji/jiujiantang-services/analysis/handler"
)

var (
	filename string
)

func init() {
	flag.StringVar(&filename, "filename", "", "missing filename")
}

func main() {
	flag.Parse()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}
	pulseTestDataProto, errParsePulseTestData := awsTool.ParsePulseTestData(buf.Bytes(), filename)
	if errParsePulseTestData != nil {
		panic(errParsePulseTestData)
	}
	waveDatas, errParsePayload := handler.ParsePayload(pulseTestDataProto)
	if errParsePayload != nil {
		panic(errParsePayload)
	}
	filename = strings.Replace(filename, "pbd", "txt", -1)
	f, _ := os.Create(filename)
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, waveData := range waveDatas {
		fmt.Fprintln(w, waveData)
	}
	w.Flush()
}
