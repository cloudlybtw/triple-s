package internal

import (
	"encoding/xml"
	"io"
	"os"
	"triple-s/internal/csv"
)

type ErrorXml struct {
	Code int    `xml:"code"`
	Text string `xml:"error"`
}

type Bucket struct {
	XMLName      xml.Name `xml:"Bucket"`
	CreationDate string   `xml:"CreationDate"`
	Name         string   `xml:"Name"`
}

type ListAllMyBucketsResult struct {
	Bucket []Bucket `xml:"Buckets"`
}

type ResponceXml struct {
	Code int    `xml:"code"`
	Text string `xml:"text"`
}

func MakeResponceXml(code int, text string) ([]byte, error) {
	responceXml := ResponceXml{Code: code, Text: text}

	return xml.MarshalIndent(responceXml, " ", " ")
}

func MakeErrorXml(code int, text string) ([]byte, error) {
	errorXml := ErrorXml{Code: code, Text: text}

	return xml.MarshalIndent(errorXml, " ", " ")
}

func MakeListAllMyBucketsResult(path string) ([]byte, error) {
	buckets := ListAllMyBucketsResult{
		Bucket: make([]Bucket, 0),
	}
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvparser csv.CSVParser = &csv.CSV{}
	firstLine := true

	for {
		_, err := csvparser.ReadLine(f)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if firstLine {
			firstLine = false
			continue
		}

		name, err := csvparser.GetField(0) // Name of bucket
		if err != nil {
			return nil, err
		}

		creationTime, err := csvparser.GetField(1) // creation time of bucket
		if err != nil {
			return nil, err
		}

		buckets.Bucket = append(buckets.Bucket, Bucket{
			Name:         name,
			CreationDate: creationTime,
		})
	}

	return xml.MarshalIndent(buckets, " ", " ")
}
