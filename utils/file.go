package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"swan-lib/constants"
	"swan-lib/logs"
)

func IsFileExists(filePath, fileName string) bool {
	fileFullPath := filepath.Join(filePath, fileName)
	_, err := os.Stat(fileFullPath)

	if err != nil {
		logs.GetLogger().Info(err)
		return false
	}

	return true
}

func IsFileExistsFullPath(fileFullPath string) bool {
	_, err := os.Stat(fileFullPath)

	if err != nil {
		logs.GetLogger().Info(err)
		return false
	}

	return true
}

func IsPathFile(dirFullPath string) (*bool, error) {
	fi, err := os.Stat(dirFullPath)

	if err != nil {
		logs.GetLogger().Info(err)
		return nil, err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		isFile := false
		return &isFile, nil
	case mode.IsRegular():
		isFile := true
		return &isFile, nil
	default:
		err := fmt.Errorf("unknown path type")
		logs.GetLogger().Error(err)
		return nil, err
	}
}

func GetPathType(dirFullPath string) int {
	fi, err := os.Stat(dirFullPath)

	if err != nil {
		logs.GetLogger().Info(err)
		return constants.PATH_TYPE_NOT_EXIST
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return constants.PATH_TYPE_DIR
	case mode.IsRegular():
		return constants.PATH_TYPE_FILE
	default:
		return constants.PATH_TYPE_UNKNOWN
	}
}

func RemoveFile(filePath, fileName string) {
	fileFullPath := filepath.Join(filePath, fileName)
	err := os.Remove(fileFullPath)
	if err != nil {
		logs.GetLogger().Error(err.Error())
	}
}

func GetFileSize(fileFullPath string) int64 {
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		logs.GetLogger().Info(err)
		return -1
	}

	return fi.Size()
}

func GetFileSize2(dir, fileName string) int64 {
	fileFullPath := filepath.Join(dir, fileName)
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		logs.GetLogger().Info(err)
		return -1
	}

	return fi.Size()
}

func CopyFile(srcFilePath, destFilePath string) (int64, error) {
	sourceFileStat, err := os.Stat(srcFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		err = errors.New(srcFilePath + " is not a regular file")
		logs.GetLogger().Error(err)
		return 0, err
	}

	source, err := os.Open(srcFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}

	defer source.Close()

	destination, err := os.Create(destFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}

	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}

	return nBytes, err
}

func CreateFileWithContents(filepath string, lines []string) (int, error) {
	f, err := os.Create(filepath)

	if err != nil {
		logs.GetLogger().Error(err)
		return 0, nil
	}

	defer f.Close()

	bytesWritten := 0
	for _, line := range lines {
		bytesWritten1, err := f.WriteString(line + "\n")
		if err != nil {
			logs.GetLogger().Error(err)
			return 0, nil
		}
		bytesWritten = bytesWritten + bytesWritten1
	}

	if err != nil {
		logs.GetLogger().Error(err)
		return 0, nil
	}

	logs.GetLogger().Info(filepath, " generated.")
	return bytesWritten, nil
}

func ReadAllLines(dir, filename string) ([]string, error) {
	fileFullPath := filepath.Join(dir, filename)

	file, err := os.Open(fileFullPath)

	if err != nil {
		logs.GetLogger().Error("failed opening file: ", fileFullPath)
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func ReadFile(filePath string) (string, []byte, error) {
	sourceFileStat, err := os.Stat(filePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", nil, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		err = errors.New(filePath + " is not a regular file")
		logs.GetLogger().Error(err)
		return "", nil, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logs.GetLogger().Error("failed reading data from file: ", filePath)
		logs.GetLogger().Error(err)
		return "", nil, err
	}

	return sourceFileStat.Name(), data, nil
}

func IsDirExists(dir string) bool {
	if len(dir) == 0 {
		err := fmt.Errorf("dir is not provided")
		logs.GetLogger().Error(err)
		return false
	}

	if GetPathType(dir) != constants.PATH_TYPE_DIR {
		err := fmt.Errorf("%s is not a directory", dir)
		logs.GetLogger().Error(err)
		return false
	}

	return true
}

func CreateDir(dir string) error {
	if len(dir) == 0 {
		err := fmt.Errorf("dir is not provided")
		logs.GetLogger().Info(err)
		return err
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		err := fmt.Errorf("%s, failed to create output dir:%s", err.Error(), dir)
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func CallGenerateFile() {
	fmt.Println("usage: generateFile filepath filename filesizeInGigabyte")
	filepath := os.Args[1]
	filename := os.Args[2]
	filesizeInGigabyte, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	GenerateFile(filepath, filename, filesizeInGigabyte)
}

func GenerateFile(filepath, filename string, filesize int64) {
	filefullpath := filepath + "/" + filename
	file, err := os.Create(filefullpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesizeInByte := filesize * 100000000
	var i int64
	for i = 0; i < filesizeInByte; i++ {
		_, err := file.WriteString("Hello World")
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}
		//fmt.Println(l, "bytes written successfully")
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
