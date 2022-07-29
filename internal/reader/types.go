/*
Copyright © 2022 Aurelio Calegari, et al.

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

package reader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	parentPath = ".loggo"
	paramsPath = "params"
)

type SavedParams struct {
	From     string `yaml:"from,omitempty"`
	Filter   string `yaml:"filter,omitempty"`
	Project  string `yaml:"project,omitempty"`
	Template string `yaml:"template,omitempty"`
}

func (spw *SavedParams) Print() {
	if len(spw.From) > 0 {
		fmt.Println()
		fmt.Printf(`          - From:     %s`, spw.From)
	}
	if len(spw.Project) > 0 {
		fmt.Println()
		fmt.Printf(`          - Project:  %s`, spw.Project)
	}
	if len(spw.Template) > 0 {
		fmt.Println()
		fmt.Printf(`          - Template: %s`, spw.Template)
	}
	if len(spw.Filter) > 0 {
		fmt.Println()
		fmt.Printf(`          - Filter:   %s`, spw.Filter)
	}
}

type SavedParamsWrapper struct {
	Name   string      `yaml:"name"`
	Params SavedParams `yaml:"params"`
}

func (spw *SavedParamsWrapper) Print() {
	fmt.Println()
	fmt.Printf(`Name:     %s`, spw.Name)
	fmt.Println()
	fmt.Printf(`Execute:  loggo gcp-stream --params-load "%s"`, spw.Name)
	fmt.Println()
	fmt.Printf(`Contents:`)
	spw.Params.Print()
	fmt.Println()
	fmt.Println()

}

func Save(name string, s *SavedParams) error {
	paramsDir, err := paramDirectory()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(paramsDir, os.ModePerm); err != nil {
		return err
	}
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	n, fDir, err := validateName(name)
	f, err := os.Create(fDir)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	sw := SavedParamsWrapper{
		Name:   n,
		Params: *s,
	}

	sw.Print()

	fmt.Println()
	fmt.Printf("Saved at %v", fDir)
	return nil
}

func Load(name string) (*SavedParams, error) {
	_, paramsDir, err := validateName(name)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(paramsDir)
	if err != nil {
		return nil, err
	}
	sp := SavedParams{}
	if err := yaml.Unmarshal(b, &sp); err != nil {
		return nil, err
	}
	return &sp, nil
}

func validateName(name string) (string, string, error) {
	paramsDir, err := paramDirectory()
	if err != nil {
		return "", "", err
	}
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	if len(name) == 0 {
		return "", "", fmt.Errorf("invalid name")
	}
	return name, path.Join(paramsDir, name), nil
}

func Remove(name string) error {
	name, paramsDir, err := validateName(name)
	if err == nil {
		err = os.Remove(paramsDir)
	}
	return err
}

func List() ([]*SavedParamsWrapper, error) {
	paramsDir, err := paramDirectory()
	fmt.Println()
	fmt.Println("Listing files at ", paramsDir)
	fmt.Println()
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(paramsDir)
	if err != nil {
		return nil, err
	}

	list := make([]*SavedParamsWrapper, 0)
	for _, file := range files {
		if !file.IsDir() {
			sp, err := Load(file.Name())
			if err == nil {
				list = append(list, &SavedParamsWrapper{
					Name:   file.Name(),
					Params: *sp,
				})
			}
		}

	}
	sort.SliceStable(list, func(i, j int) bool {
		return strings.Compare(list[i].Name, list[j].Name) < 0
	})
	return list, nil
}

func paramDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	paramsDir := path.Join(home, parentPath, paramsPath)
	return paramsDir, nil
}
