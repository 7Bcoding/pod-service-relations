package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTarGzip(t *testing.T) {
	tg := NewTarGzip()
	tempDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error("create tmp dir error")
	}
	defer os.RemoveAll(tempDir)

	tempDir2, err := ioutil.TempDir(tempDir, "")
	tempFile, err := ioutil.TempFile(tempDir2, "test-file")
	if err != nil {
		t.Error("create tmp file error")
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte{1, 2, 3, 4})
	if err != nil {
		t.Error("write to temp file error")
	}
	var buf bytes.Buffer
	if err := tg.Compress(tempDir2, &buf); err != nil {
		t.Error("compress error")
	}

	fileToWrite, err := os.OpenFile(filepath.Join(tempDir, "test.tar.gz"), os.O_CREATE|os.O_RDWR, os.FileMode(777))
	if err != nil {
		t.Error("create tar file error")
	}

	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		t.Error("compress io copy error")
	}
	//fi, err := os.Open("/Users/luoyongchang/Desktop/tmp/cloud_service_temp_dir/baidu-bce-cds_2.0.77.14460.tar.gz")
	//if err != nil {

	//if err != nil {
	//	panic(err)
	//}
	//err = tg.Uncompress(fi, "/Users/luoyongchang/Desktop/tmp/cloud_service_temp_dir")

	f, err := os.Open(filepath.Join(tempDir, "test.tar.gz"))
	if err != nil {
		t.Error("open file error")
		return
	}

	if err := tg.Uncompress(f, tempDir); err != nil {
		t.Error(fmt.Sprintf("uncompress errro, err %s", err.Error()))
	}
}
