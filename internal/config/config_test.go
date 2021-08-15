package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfigFile_WhenExistingFilenameIsProvided_ShouldReadFileData(t *testing.T) {
	filename := "testdata/test_config_1.txt"

	got, err := ReadConfigFile(filename)

	assert.Nil(t, err, "No error should be returned for valid filename")
	assert.Equal(t, "path = \"$HOME/app\"\n", string(got), "The file was not read properly")
}

func TestReadConfigFile_WhenWrongFilenameIsProvided_ShouldReturnError(t *testing.T) {
	filename := "testdata/wrong_filename.txt"

	got, err := ReadConfigFile(filename)

	assert.Nil(t, got, "No result should be returned when error occurs")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "open testdata/wrong_filename.txt: no such file or directory",
		err.Error(), "Incorrect error message")
}

func TestReadMultipleConfigFiles_WhenExistingFilenamesAreProvided_ShouldReadAllFileData(t *testing.T) {
	filenames := []string{
		"testdata/test_config_1.txt",
		"testdata/test_config_2.txt",
	}

	data, err := ReadMultipleConfigFiles(filenames)
	got := convertBytesToStrings(data)

	assert.Nil(t, err, "No error should be returned for valid filenames")
	expected := []string{
		"path = \"$HOME/app\"\n",
		"a = 1\nb = 2\n",
	}
	assert.Equal(t, expected, got, "The file was not read properly")
}

func TestReadMultipleConfigFiles_WhenErrorOccursDuringFileProcessing_ShouldReturnError(t *testing.T) {
	filenames := []string{
		"testdata/wrong_filename.txt",
		"testdata/test_config_2.txt",
	}

	got, err := ReadMultipleConfigFiles(filenames)

	assert.Nil(t, got, "No result should be returned when error occurs")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "open testdata/wrong_filename.txt: no such file or directory",
		err.Error(), "Incorrect error message")
}

func convertBytesToStrings(data [][]byte) []string {
	result := make([]string, len(data))

	for i, bytes := range data {
		result[i] = string(bytes)
	}

	return result
}
