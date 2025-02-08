package utils

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gzip "github.com/klauspost/pgzip"
)

type TarGzip interface {
	Compress(src string, buf io.Writer) error
	Uncompress(src io.Reader, dst string) error
}

type tarGzip struct {
}

func NewTarGzip() TarGzip {
	return &tarGzip{}
}

func (tg tarGzip) Compress(src string, buf io.Writer) error {
	zw := gzip.NewWriter(buf)
	tw := tar.NewWriter(zw)
	defer zw.Close()
	defer tw.Close()

	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	mode := fi.Mode()

	if mode.IsRegular() {
		header, err := tar.FileInfoHeader(fi, src)
		if err != nil {
			return err
		}
		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// get content
		data, err := os.Open(src)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}
	} else if mode.IsDir() {
		err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// must provide real name
			// (see https://golang.org/src/archive/tar/common.go?#L626)
			//header.Name = filepath.ToSlash(file)
			relPath, err := filepath.Rel(src, file)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(relPath)
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				defer data.Close()
				if _, err := io.Copy(tw, data); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("error: file type not supported")
	}

	return nil
}

func (tg tarGzip) Uncompress(src io.Reader, dst string) error {
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	defer zr.Close()

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// validate name against path traversal
		if !tg.validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q\n", header.Name)
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}
	return nil
}

func (tg tarGzip) validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
