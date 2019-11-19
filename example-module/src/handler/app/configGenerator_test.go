package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// create tests for module configuration

func copyConfig(t *testing.T, filepath string) {
	err := os.MkdirAll(configSourcePath, 0777)
	assert.Nil(t, err, "MkdirAll")
	err = copyFile(filepath, path.Join(configSourcePath, "module.json"))
	assert.Nil(t, err, "copyFile")
}

func prepareTestLuaConfTemplate(t *testing.T, srcPath string) {
	err := os.MkdirAll("/opt/section", 0777)
	assert.Nil(t, err, "MkdirAll")
	err = copyFile(srcPath, luaConfTemplatePath)
	assert.Nil(t, err, "copyFile")
	initTemplate()
}

func prepareTestSchemaJSON(t *testing.T, srcPath string) {
	err := os.MkdirAll("/opt/section", 0777)
	assert.Nil(t, err, "MkdirAll")
	err = copyFile(srcPath, moduleSchemaPath)
	assert.Nil(t, err, "copyFile")
}

func TestReadProxyFeaturesJson(t *testing.T) {
	copyConfig(t, "testdata/valid-module.json")
	prepareTestSchemaJSON(t, "testdata/module_schema.json")
	config, err := readProxyFeaturesJSON()
	assert.Nil(t, err, "readProxyFeaturesJSON")
	assert.True(t, config.Enabled, "config.Enabled")
}

func TestReadInvalidProxyFeaturesJson(t *testing.T) {
	copyConfig(t, "testdata/invalid-module.json")
	prepareTestSchemaJSON(t, "testdata/module_schema.json")
	_, err := readProxyFeaturesJSON()
	assert.Error(t, err, "readProxyFeaturesJSON")
}

func TestWritingConfigFile(t *testing.T) {
	prepareTestLuaConfTemplate(t, "testdata/environment_variables.lua.gotemplate")
	prepareTestSchemaJSON(t, "testdata/module_schema.json")
	copyConfig(t, "testdata/valid-module.json")
	generateLuaFile()
	fileContents, err := ioutil.ReadFile(path.Join(luaLibPath, luaLibFilename))
	assert.Nil(t, err, "ioutil.ReadFile")

	expected := `local _M={_VERSION = '1.0.0'}

_M.ENABLED=true;
_M.API_KEY="fdsafdsafdsafdsafdsafdsafdsa0123";
_M.EXAMPLE="some_parameter";

return _M
`
	assert.Equal(t, expected, string(fileContents))
}
