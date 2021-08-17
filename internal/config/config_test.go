package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	configFileWithPath    = "testdata/test_config_with_path.txt"
	configFileTwoIntegers = "testdata/test_config_with_two_integers.txt"
	nonExistingFile       = "testdata/non_existent_file.txt"
)

func TestReadConfigFile_WhenExistingFilenameIsProvided_ShouldReadFileData(t *testing.T) {
	got, err := ReadConfigFile(configFileWithPath)

	require.NoError(t, err, "No error should be returned for valid filename")
	assert.Equal(t, "path = \"$HOME/app\"\n", string(got), "The file was not read properly")
}

func TestReadConfigFile_WhenWrongFilenameIsProvided_ShouldReturnError(t *testing.T) {
	_, err := ReadConfigFile(nonExistingFile)

	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "open testdata/non_existent_file.txt: no such file or directory",
		err.Error(), "Incorrect error message")
}

func TestReadMultipleConfigFiles_WhenExistingFilenamesAreProvided_ShouldReadAllFileData(t *testing.T) {
	filenames := []string{
		configFileWithPath,
		configFileTwoIntegers,
	}

	data, err := ReadMultipleConfigFiles(filenames)
	got := convertBytesToStrings(data)

	require.NoError(t, err, "No error should be returned for valid filenames")
	expected := []string{
		"path = \"$HOME/app\"\n",
		"a = 1\nb = 2\n",
	}
	assert.Equal(t, expected, got, "The file was not read properly")
}

func TestReadMultipleConfigFiles_WhenErrorOccursDuringFileProcessing_ShouldReturnError(t *testing.T) {
	filenames := []string{
		nonExistingFile,
		configFileWithPath,
	}

	_, err := ReadMultipleConfigFiles(filenames)

	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "open testdata/non_existent_file.txt: no such file or directory",
		err.Error(), "Incorrect error message")
}

func convertBytesToStrings(data [][]byte) []string {
	result := make([]string, len(data))

	for i, bytes := range data {
		result[i] = string(bytes)
	}

	return result
}
