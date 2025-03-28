package csv

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrCantFind       = errors.New("Can't find metadata")
	ErrCantFindBucket = errors.New("Can't find bucket with this name")
	objects           = "objects.csv"
)

func GetContentType(path, bucketName, objectName string) (string, error) {
	found, err := IsExistBucket(path, bucketName)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrCantFindBucket
	}
	fullpath := strings.Join([]string{path, bucketName, objects}, "/")
	var csvparser CSVParser = &CSV{}
	file, err := os.Open(fullpath)
	if err != nil {
		return "", err
	}
	for {
		_, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		name, err := csvparser.GetField(0)
		if err != nil {
			return "", err
		}
		if name == objectName {
			content_type, err := csvparser.GetField(2)
			if err != nil {
				return "", err
			}
			return content_type, nil
		}
	}
	return "", ErrCantFind
}

func UpdateObject(path, bucketName, objectName, ContentType string, ContentLength int64) error {
	found, err := IsExistBucket(path, bucketName)
	if err != nil {
		return err
	}
	if !found {
		return ErrCantFindBucket
	}
	fullPath := strings.Join([]string{path, bucketName, objects}, "/")
	file, err := os.OpenFile(fullPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	var csvparser CSVParser = &CSV{}
	lines := make([]byte, 0)
	found = false

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		name, err := csvparser.GetField(0) // object name
		if err != nil {
			return err
		}
		if name == objectName {
			LastModifiedTime := time.Now().Format(time.RFC3339)
			ContentLength := strconv.Itoa(int(ContentLength))

			lines = append(lines, []byte(strings.Join([]string{objectName, ContentLength, ContentType, LastModifiedTime + "\n"}, ","))...)
			found = true
		} else {
			lines = append(lines, []byte(line+"\n")...)
		}
	}
	if !found {
		LastModifiedTime := time.Now().Format(time.RFC3339)
		ContentLength := strconv.Itoa(int(ContentLength))

		lines = append(lines, []byte(strings.Join([]string{objectName, ContentLength, ContentType, LastModifiedTime + "\n"}, ","))...)
	}
	err = os.WriteFile(fullPath, lines, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func DeleteObject(path, bucketName, objectName string) (bool, error) {
	found, err := IsExistBucket(path, bucketName)
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrCantFindBucket
	}
	fullPath := strings.Join([]string{path, bucketName, objects}, "/")
	file, err := os.OpenFile(fullPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return false, err
	}

	var csvparser CSVParser = &CSV{}

	lines := make([]byte, 0)
	found = false
	lineCount := 0
	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, err
		}
		name, err := csvparser.GetField(0) // object name
		if err != nil {
			return false, err
		}
		if name != objectName {
			lines = append(lines, []byte(line+"\n")...)
		} else {
			found = true
		}
		lineCount++
	}
	if !found {
		return false, ErrCantFind
	}
	err = os.WriteFile(fullPath, lines, os.ModePerm)
	if err != nil {
		return false, err
	}
	if lineCount <= 2 {
		return true, nil
	}
	return false, nil
}
