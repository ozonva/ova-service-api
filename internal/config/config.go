package config

import (
	"io"
	"os"
)

func ReadConfigFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		closeErr := f.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(file)

	data := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, readErr := file.Read(buffer)

		if n > 0 {
			data = append(data, buffer[:n]...)
		}

		switch readErr {
		case nil:
			break
		case io.EOF:
			return data, nil
		default:
			return nil, readErr
		}
	}
}

func ReadMultipleConfigFiles(filenames []string) ([][]byte, error) {
	data := make([][]byte, 0)

	for _, filename := range filenames {
		fileData, err := ReadConfigFile(filename)

		if err != nil {
			return nil, err
		}

		data = append(data, fileData)
	}

	return data, nil
}
