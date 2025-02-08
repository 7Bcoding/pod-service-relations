package client

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"pod-service-relations/logging"
	"strings"
)

// GetFilesAndDirs 获取指定目录下的所有文件和目录
func GetFilesAndDirs(rootPath string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, nil, err
	}

	pathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fileInfo := range dir {
		fiPath := rootPath + pathSep + fileInfo.Name()
		if fileInfo.IsDir() { // 是目录, dfs递归遍历
			dirs = append(dirs, fiPath)
			GetFilesAndDirs(fiPath)
		} else { // 是文件
			// 过滤指定格式
			ok := strings.HasSuffix(fileInfo.Name(), ".go")
			if ok {
				files = append(files, fiPath)
			}
		}
	}

	return files, dirs, nil
}

// GetAllFilesRecursion 递归的获取指定目录下的所有文件,包含子目录下的文件
func GetAllFilesRecursion(rootPath string) (files []string, err error) {
	var dirs []string
	dir, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	pathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fileInfo := range dir {
		fiPath := rootPath + pathSep + fileInfo.Name()
		if fileInfo.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, fiPath)
			_, err2 := GetAllFilesRecursion(fiPath)
			if err2 != nil {
				return nil, err2
			}
		} else {
			files = append(files, fiPath)
		}
	}

	// 读取子目录下文件
	for _, subDir := range dirs {
		tempFiles, _ := GetAllFilesRecursion(subDir)
		for _, tempFile := range tempFiles {
			files = append(files, tempFile)
		}
	}

	return files, nil
}

// GetAllFiles 获取一个目录下的文件的绝对路径path
func GetAllFiles(rootPath string) (files []string, err error) {
	dir, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	pathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fileInfo := range dir {
		fiPath := rootPath + pathSep + fileInfo.Name()
		files = append(files, fiPath)
	}
	return files, nil

}

// SplitListToStrByPattern 按照surrPattern包裹每个字段，然后按照splitPattern分割各字段
func SplitListToStrByPattern(elemlist []string, splitPattern string, surrPattern string) string {
	retStr := ""
	for i, elem := range elemlist {
		if i == len(elemlist)-1 {
			retStr += surrPattern + elem + surrPattern
		} else {
			retStr += surrPattern + elem + surrPattern + splitPattern
		}
	}
	return retStr
}

// PathExists 检查一个文件夹/文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Tar 按tar格式压缩打包文件或文件夹
func Tar(ctx context.Context, src string, dstTar string, failIfExist bool) (err error) {
	// 清理路径字符串
	src = path.Clean(src)

	// 判断要打包的文件或目录是否存在
	if !Exists(src) {
		return errors.New("要打包的文件或目录不存在：" + src)
	}

	// 判断目标文件是否存在
	if FileExists(dstTar) {
		if failIfExist { // 不覆盖已存在的文件
			return errors.New("目标文件已经存在：" + dstTar)
		} else { // 覆盖已存在的文件
			err := os.Remove(dstTar)
			if err != nil {
				return err
			}
		}
	}

	// 创建空的目标文件
	fw, err := os.Create(dstTar)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 创建 tar.Writer，执行打包操作
	tw := tar.NewWriter(fw)
	defer func() {
		// 这里要判断 tw 是否关闭成功，如果关闭失败，则 .tar 文件可能不完整
		if err := tw.Close(); err != nil {
			logging.GetLogger().Errorln("tw closed failed, got error: %s", err)
		}
	}()

	// 获取文件或目录信息
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 获取要打包的文件或目录的所在位置和名称
	srcBase, srcRelative := path.Split(path.Clean(src))

	// 开始打包
	if fi.IsDir() {
		err = TarDir(srcBase, srcRelative, tw, fi)
		if err != nil {
			return err
		}
	} else {
		err = TarFile(srcBase, srcRelative, tw, fi)
		if err != nil {
			return err
		}
	}

	return nil
}

// TarDir 打包文件夹(因为要执行遍历操作，所以要单独创建一个函数)
func TarDir(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 在结尾添加 "/"
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	// 获取 srcFull 下的文件或子目录列表
	fis, err := ioutil.ReadDir(srcFull)
	if err != nil {
		return err
	}

	// 开始遍历
	for _, fi := range fis {
		if fi.IsDir() {
			if err := TarDir(srcBase, srcRelative+fi.Name(), tw, fi); err != nil {
				return err
			}
		} else {
			if err := TarFile(srcBase, srcRelative+fi.Name(), tw, fi); err != nil {
				return err
			}
		}
	}

	// 写入目录信息
	if len(srcRelative) > 0 {
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		hdr.Name = srcRelative

		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}
	}

	return nil
}

// TarFile 打包file,因为要在 defer 中关闭文件，所以要单独创建一个函数
func TarFile(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 写入文件信息
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	hdr.Name = srcRelative

	if err = tw.WriteHeader(hdr); err != nil {
		return err
	}

	// 打开要打包的文件，准备读取
	fr, err := os.Open(srcFull)
	if err != nil {
		return err
	}
	defer fr.Close()

	// 将文件数据写入 tw 中
	if _, err = io.Copy(tw, fr); err != nil {
		return err
	}
	return nil
}

// UnTar 解压tar格式文件
func UnTar(srcTar string, dstDir string) (err error) {
	// 清理路径字符串
	dstDir = path.Clean(dstDir) + string(os.PathSeparator)

	// 打开要解包的文件
	fr, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	defer fr.Close()

	// 创建 tar.Reader，准备执行解包操作
	tr := tar.NewReader(fr)

	// 遍历包中的文件
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return err
		}

		// 获取文件信息
		fi := hdr.FileInfo()

		// 获取绝对路径
		dstFullPath := dstDir + hdr.Name

		if hdr.Typeflag == tar.TypeDir {
			// 创建目录
			if err := os.MkdirAll(dstFullPath, fi.Mode().Perm()); err != nil {
				return err
			}
			// 设置目录权限
			if err := os.Chmod(dstFullPath, fi.Mode().Perm()); err != nil {
				return err
			}
		} else {
			// 创建文件所在的目录
			if err := os.MkdirAll(path.Dir(dstFullPath), os.ModePerm); err != nil {
				return err
			}
			// 将 tr 中的数据写入文件中
			if err := UnTarFile(dstFullPath, tr); err != nil {
				return err
			}
			// 设置文件权限
			if err := os.Chmod(dstFullPath, fi.Mode().Perm()); err != nil {
				return err
			}
		}
	}
	return nil
}

// UnTarFile 因为要在 defer 中关闭文件，所以要单独创建一个函数
func UnTarFile(dstFile string, tr *tar.Reader) error {
	// 创建空文件，准备写入解包后的数据
	fw, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 写入解包后的数据
	_, err = io.Copy(fw, tr)
	if err != nil {
		return err
	}

	return nil
}

// Exists 判断文件是否存在
func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// FileExists 判断文件是否存在
func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

// DirExists 判断目录是否存在
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

// GetFileMd5 计算文件的md5值
func GetFileMd5(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(m.Sum(nil)), err
}
