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

import (
	"github.com/gotameme/core/ant" // This imports the ant framework to interact with the simulation.
	"math/rand" // Used for making random decisions.
)

// {{ .Name }} is an example implementation of an ant.
// Think of it as a small agent operating within a virtual world.
type {{ .Name }} struct {
	// We embed the AntOS to enable our ant to interact with the simulation environment.
	ant.AntOs
}

// NewAnt creates a new ant. This is the only function that needs to be visible outside.
// It's where the journey of our ant begins.
func NewAnt(antOS ant.AntOs) ant.Ant {
	return &{{ .Name }}{antOS}
}

// Waits is called when the ant has nothing specific to do.
// It's the ant's downtime, wondering "What's next?"
func (a *{{ .Name }}) Waits() {
	// Choose a random direction because spontaneity keeps things interesting.
	randomDirection := rand.Intn(360) // A number between 0 and 359 for direction.
	a.Turn(randomDirection) // Turn in this random direction.
	// Move forward 60 steps because staying active is key!
	a.GoForward(60)
}

// SeeSugar is called when the ant spots sugar.
// In this world, sugar is akin to treasure.
func (a *{{ .Name }}) SeeSugar(sugar ant.Sugar) {
	if a.GetCurrentLoad() > 0 {
		// Already carrying sugar. No need to go to another one.
		return
	}
	// "Hmm, there's sugar. Let's head over!"
	a.GoToSugar(sugar)
}

// ReachedSugar is called when the ant has reached sugar.
func (a *{{ .Name }}) ReachedSugar(sugar ant.Sugar) {
	// Successfully reached sugar. Now, let's take it.
	_ = a.TakeSugar(sugar)
	// Time to head back to the anthill.
	a.GotToAntHill()
}

// SeeFriend is called when the ant spots another ant.
func (a *{{ .Name }}) SeeFriend(friend ant.Ant) {
	// Interactions when spotting a friend could happen here.
	// Perhaps a simple acknowledgement or sharing information?
	// TODO: implement me
}

// SeeMark is called when the ant spots a mark.
func (a *{{ .Name }}) SeeMark(mark ant.Mark) {
	// Marks are like notes left by other ants.
	// Could be directions or important information.
	// TODO: implement me
}

// Tick is called at every simulation tick.
func (a *{{ .Name }}) Tick() {
	// Each "Tick" is like a heartbeat in our simulation.
	// What does our ant do with every beat?
	// TODO: implement me
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
