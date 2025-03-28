package csv

import (
	"errors"
	"io"
)

type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

type CSV struct {
	currentLine string
	fields      []string
	numOfFields int
}

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
	ErrComma      = errors.New("the last char cant be a comma")
)

func (c *CSV) ReadLine(r io.Reader) (string, error) {
	var line string
	var fields []string
	var countQuotes int
	buf := make([]byte, 1)
	inQuotes := false

	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF && n == 0 {
			if line == "" {
				return "", io.EOF
			}
			break
		}

		if buf[0] == '"' {
			countQuotes++
			inQuotes = !inQuotes
		}

		if buf[0] == '\n' || buf[0] == '\r' {
			if !inQuotes {
				break
			}
		}

		line += string(buf[0])

	}
	if len(line) > 0 && line[len(line)-1] == ',' {
		return "", ErrComma
	}

	if countQuotes%2 != 0 {
		return "", ErrQuote
	}

	c.currentLine = line
	fields, err := c.parseFields(line)
	if err != nil {
		return "", err
	}
	c.fields = fields
	c.numOfFields = len(c.fields)

	return line, nil
}

func (c *CSV) GetField(n int) (string, error) {
	if n < 0 || n >= c.numOfFields {
		return "", ErrFieldCount
	}
	return c.fields[n], nil
}

func (c *CSV) GetNumberOfFields() int {
	return c.numOfFields
}

func (c *CSV) parseFields(line string) ([]string, error) {
	var fields []string
	var currentField []byte
	inQuotes := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		if char == '"' {
			if inQuotes {
				if i+1 < len(line) && line[i+1] == '"' {
					currentField = append(currentField, char)
					i++
				} else {
					inQuotes = false
				}
			} else {
				inQuotes = true
			}
		} else if char == ',' && !inQuotes {
			fields = append(fields, string(currentField))
			currentField = []byte{}
		} else {
			currentField = append(currentField, char)
		}
	}

	if len(currentField) > 0 {
		fields = append(fields, string(currentField))
	}

	return fields, nil
}
