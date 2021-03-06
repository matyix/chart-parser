package utils

import (
	"github.com/kubernetes/helm/pkg/proto/hapi/chart"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
	"errors"
)

// ApiVersionV1 is the API version number for version 1.
//
// This is ApiVersionV1 instead of APIVersionV1 to match the protobuf-generated name.
const ApiVersionV1 = "v1" // nolint

// UnmarshalChartfile takes raw Chart.yaml data and unmarshals it.
func UnmarshalChartfile(data []byte) (*chart.Metadata, error) {
	y := &chart.Metadata{}
	err := yaml.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	return y, nil
}

// LoadChartfile loads a Chart.yaml file into a *chart.Metadata.
func LoadChartfile(filename string) (*chart.Metadata, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshalChartfile(b)
}

// SaveChartfile saves the given metadata as a Chart.yaml file at the given path.
//
// 'filename' should be the complete path and filename ('foo/Chart.yaml')
func SaveChartfile(filename string, cf *chart.Metadata) error {
	out, err := yaml.Marshal(cf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0644)
}

// IsChartDir validate a chart directory.
//
// Checks for a valid Chart.yaml.
func IsChartDir(dirName string) (bool, error) {
	if fi, err := os.Stat(dirName); err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, fmt.Errorf("%q is not a directory", dirName)
	}

	chartYaml := filepath.Join(dirName, "Chart.yaml")
	if _, err := os.Stat(chartYaml); os.IsNotExist(err) {
		return false, fmt.Errorf("no Chart.yaml exists in directory %q", dirName)
	}

	chartYamlContent, err := ioutil.ReadFile(chartYaml)
	if err != nil {
		return false, fmt.Errorf("cannot read Chart.Yaml in directory %q", dirName)
	}

	chartContent, err := UnmarshalChartfile(chartYamlContent)
	if err != nil {
		return false, err
	}
	if chartContent == nil {
		return false, errors.New("chart metadata (Chart.yaml) missing")
	}
	if chartContent.Name == "" {
		return false, errors.New("invalid chart (Chart.yaml): name must not be empty")
	}

	return true, nil
}