package csv

import (
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

var (
	buckets       = "buckets.csv"
	deleteAble    = "deleteAble"
	notDeleteAble = "notDeletable"
	ErrCantDelete = errors.New("Can't delete bucket contain data")
)

func UpdateBucket(path, bucketName string, canDelete bool) error {
	fullPath := strings.Join([]string{path, buckets}, "/")
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	var csvparser CSVParser = &CSV{}
	lines := make([]byte, 0)
	found := false

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		name, err := csvparser.GetField(0)
		if err != nil {
			return err
		}
		CreatedTime, err := csvparser.GetField(1)
		if err != nil {
			return err
		}
		if name == bucketName {
			LastModifiedTime := time.Now().Format(time.RFC3339)
			found = true
			// if i am updating bucket, its mean, i can't delete it any time
			if canDelete {
				lines = append(lines, []byte(bucketName+","+CreatedTime+","+LastModifiedTime+","+deleteAble+"\n")...)
			} else {
				lines = append(lines, []byte(bucketName+","+CreatedTime+","+LastModifiedTime+","+notDeleteAble+"\n")...)
			}
		} else {
			lines = append(lines, []byte(line+"\n")...)
		}
	}

	if !found {
		currentTime := time.Now().Format(time.RFC3339)
		LastModifiedTime := time.Now().Format(time.RFC3339)
		lines = append(lines, []byte(bucketName+","+currentTime+","+LastModifiedTime+","+deleteAble+"\n")...)
	}

	err = os.WriteFile(fullPath, lines, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBucket(path, bucketName string) error {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	var csvparser CSVParser = &CSV{}
	found := false
	canDelete := false
	lines := make([]byte, 0)

	for {
		line, err := csvparser.ReadLine(f)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		s, err := csvparser.GetField(0)
		if err != nil {
			return err
		}
		delet, err := csvparser.GetField(3)
		if err != nil {
			return err
		}
		if s == bucketName {
			found = true
			if delet == deleteAble {
				canDelete = true
			} else {
				lines = append(lines, []byte(line+"\n")...)
			}
		} else {
			lines = append(lines, []byte(line+"\n")...)
		}

	}

	if !found {
		return ErrCantFind
	}

	if !canDelete {
		return ErrCantDelete
	}

	err = os.WriteFile(path, lines, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func IsExistBucket(path, bucketName string) (bool, error) {
	fullPath := strings.Join([]string{path, "buckets.csv"}, "/")
	file, err := os.Open(fullPath)
	if err != nil {
		return false, err
	}

	var csvparser CSVParser = &CSV{}
	for {
		_, err := csvparser.ReadLine(file)
		if err != nil {
			if err != io.EOF {
				break
			}
			return false, err
		}
		name, err := csvparser.GetField(0)
		if err != nil {
			return false, err
		}
		if name == bucketName {
			return true, nil
		}
	}
	return false, nil
}

func IsExistObject(path, bucketName, objectName string) (bool, error) {
	fullPath := strings.Join([]string{path, bucketName, objects}, "/")
	file, err := os.Open(fullPath)
	if err != nil {
		return false, err
	}

	var csvparser CSVParser = &CSV{}
	for {
		_, err := csvparser.ReadLine(file)
		if err != nil {
			if err != io.EOF {
				break
			}
			return false, err
		}
		name, err := csvparser.GetField(0)
		if err != nil {
			return false, err
		}
		if name == objectName {
			return true, nil
		}
	}
	return false, nil
}
