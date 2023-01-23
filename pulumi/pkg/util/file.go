package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func WriteTempFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", fmt.Sprintf("%s-", filepath.Base(os.Args[0])))
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = tmpFile.WriteString(content)
	return tmpFile.Name(), err
}

func Template(data any, templateName, templateContent string) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
