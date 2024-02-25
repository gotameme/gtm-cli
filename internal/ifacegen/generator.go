/*
Copyright (c) 2024 Sebastian Kroczek <me@xbug.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package ifacegen

import (
	"bytes"
	"fmt"
	"github.com/gotameme/gtm-cli/internal"
	"html/template"
	"os"
	"os/exec"
)

const antTemplate = `package main
import(
	"github.com/gotameme/core/ant"
)

type {{ .Name }} struct {
	ant.AntOs
}

func NewAnt(antOS ant.AntOs) ant.Ant {
	return &{{ .Name }}{antOS}
}

func (d *{{ .Name }}) Waits() {
	// Implement me
}

func (d *{{ .Name }}) SeeSugar(sugar ant.Sugar) {
	// Implement me
}

func (d *{{ .Name }}) ReachedSugar(sugar ant.Sugar) {
	// Implement me
}

func (d *{{ .Name }}) SeeFriend(ant ant.Ant) {
	// Implement me
}

func (d *{{ .Name }}) SeeMark(mark ant.Mark) {
	// Implement me
}

func (d *{{ .Name }}) Tick() {
	// Implement me
}
`

func GenerateInterfaceCode(interfaceName string) (code string, err error) {
	if err = validateInterfaceName(interfaceName); err != nil {
		return
	}

	t := template.Must(template.New("ant").Parse(antTemplate))

	data := struct {
		Name string
	}{
		Name: interfaceName,
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, data); err != nil {
		return "", err
	}

	code = tpl.String()
	return
}

func createFile(name, filename, content string) error {
	// Create the directory if it doesn't exist
	if !DirectoryExists(name) {
		if err := CreateDirectory(name); err != nil {
			return err
		}
	}

	// Create the file
	file, err := os.Create(fmt.Sprintf("%s/%s", name, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	// Save the changes
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func CreateMainFile(name string) error {
	content, err := GenerateInterfaceCode(name)
	if err != nil {
		return err
	}
	return createFile(name, "main.go", content)
}

func CreateModFile(name string) error {
	// change into directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}
	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("error changing into directory: %v", err)
	}
	// run go mod init
	if err := exec.Command("go", "mod", "init", name).Run(); err != nil {
		return fmt.Errorf("error running go mod init: %v", err)
	}
	// add dependencies
	coreVersion, err := internal.GetCoreVersion()
	if err != nil {
		return fmt.Errorf("error getting core version: %v", err)
	}
	if err := exec.Command("go", "get", "github.com/gotameme/core@"+coreVersion).Run(); err != nil {
		return fmt.Errorf("error running go get: %v", err)
	}
	// run go mod tidy
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return fmt.Errorf("error running go mod tidy: %v", err)
	}
	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("error changing into directory: %v", err)
	}
	// return error if any
	return nil
}
