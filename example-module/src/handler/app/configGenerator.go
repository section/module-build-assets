package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

const (
	configSourcePath    = "/opt/proxy_config"
	luaLibPath          = "/usr/local/openresty/lualib/section"
	luaLibFilename      = "environment_variables.lua"
	luaConfTemplatePath = "/opt/section/environment_variables.lua.gotemplate"
	moduleSchemaPath    = "/opt/section/module_schema.json"
)

var luaConfTemplate *template.Template

// struct mapping the module config json object
type moduleConfig struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
	Example string `json:"example"`
}

// read the configuration file for the module
// this file is commited by the user in their Section configuration directory
// the configuration is mounted to the path /opt/proxy_config
func readProxyFeaturesJSON() (*moduleConfig, error) {
	// the file is of the format <module-name>.json
	var proxyFeaturesJSONPath = path.Join(configSourcePath, "module.json")

	data, err := ioutil.ReadFile(proxyFeaturesJSONPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read '%s'", proxyFeaturesJSONPath)
	}

	// load and validate against the json schema https://github.com/section-io/module-build-assets/blob/692592d5d77da03a8bd4eac285a83dd2789cde3f/example-module/src/module_schema.json
	var schemaLoaderFilePath = path.Join("file://", moduleSchemaPath)
	var documentLoaderFilePath = path.Join("file://", proxyFeaturesJSONPath)

	schemaLoader := gojsonschema.NewReferenceLoader(schemaLoaderFilePath)
	documentLoader := gojsonschema.NewReferenceLoader(documentLoaderFilePath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	// returns the configuration json object
	if result.Valid() {
		var config moduleConfig

		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	var errorMessages = []string{}

	for _, desc := range result.Errors() {
		var errorMessage = []string{desc.Field(), desc.Description()}
		errorMessages = append(errorMessages, strings.Join(errorMessage, ": "))
	}

	return nil, errors.New(strings.Join(errorMessages, "\n"))
}

// parse the file environment_variables.lua.gotemplate
func initTemplate() {
	mustMkdirAll(luaLibPath, 0755)

	rawTemplate, err := ioutil.ReadFile(luaConfTemplatePath)
	if err != nil {
		log.Panicf("Failed to read template file '%s': %v", luaConfTemplatePath, err)
	}

	luaConfTemplate = template.
		New(luaLibFilename)

	_, err = luaConfTemplate.Parse(string(rawTemplate))
	if err != nil {
		log.Panicf("Failed to parse template '%s': %v", luaConfTemplatePath, err)
	}
}

// create the lua file based on the go template provided : https://github.com/section-io/module-build-assets/blob/692592d5d77da03a8bd4eac285a83dd2789cde3f/example-module/src/environment_variables.lua.gotemplate
func generateLuaFile() {
	config, err := readProxyFeaturesJSON()
	if err != nil {
		log.Panic(err)
	}

	if luaConfTemplate == nil {
		initTemplate()
	}

	// the created lua file is loaded as a module and can be used to alter the behaviour of the module : https://github.com/section-io/module-build-assets/blob/692592d5d77da03a8bd4eac285a83dd2789cde3f/example-module/src/proxy/content.lua#L1
	luaConfPath := path.Join(luaLibPath, luaLibFilename)
	file, err := os.Create(luaConfPath)
	if err != nil {
		log.Panicf("Failed to create '%s': %#v", luaConfPath, err)
	}
	defer mustClose(file)

	err = luaConfTemplate.Execute(file, config)
	if err != nil {
		log.Panic(err)
	}

}
