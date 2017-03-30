/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2017 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tagsfilter

import (
	"fmt"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	// Name of the plugin
	Name = "tags-filter"
	// Version of the plugin
	Version = 1

	allowSuffix     = ".allow"
	denySuffix      = ".deny"
	valuesSeparator = ","
)

// Tags filter processor implementation
type TFProcessor struct {
}

type rule struct {
	allowedValues []string
	deniedValues  []string
}

// Get new tags filter processor plugin instance
func NewTFProcessor() *TFProcessor {
	return &TFProcessor{}
}

// GetConfigPolicy returns plugin's ConfigPolicy
func (tf *TFProcessor) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	// Config keys are dynamic - not possible to include them in ConfigPolicy
	return *plugin.NewConfigPolicy(), nil
}

// Process processes metrics
func (tf *TFProcessor) Process(mts []plugin.Metric, cfg plugin.Config) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}

	rules, err := parseRules(cfg)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return mts, nil
	}

MetricsLoop:
	for _, mt := range mts {
		adding := false

		for tag, value := range mt.Tags {
			rl, ok := rules[tag]
			if !ok {
				continue
			}

			for _, deniedValue := range rl.deniedValues {
				if value == deniedValue {
					continue MetricsLoop
				}
			}

			if !adding {
				for _, allowedValue := range rl.allowedValues {
					if value == allowedValue {
						adding = true
					}
				}
			}
		}

		// Add only if iterated over every tag to make sure we checked denied values
		if adding {
			metrics = append(metrics, mt)
		}
	}

	return metrics, nil
}

func parseRules(config plugin.Config) (map[string]*rule, error) {
	rules := make(map[string]*rule)

	for key := range config {
		allowed := false
		var suffix string
		if strings.HasSuffix(key, allowSuffix) {
			allowed = true
			suffix = allowSuffix
		} else if strings.HasSuffix(key, denySuffix) {
			suffix = denySuffix
		} else {
			continue
		}

		tag := key[:len(key)-len(suffix)]
		if len(tag) == 0 {
			return nil, fmt.Errorf("Config key must contain tag name: %s", key)
		}

		sRule, err := config.GetString(key)
		if err != nil {
			return nil, err
		}
		values := strings.Split(sRule, valuesSeparator)

		// There must be at least one value per rule
		if len(values) == 0 {
			return nil, fmt.Errorf("Rule must contain at least one value: %s", sRule)
		}

		// Check if any of values is empty
		for _, value := range values {
			if len(value) == 0 {
				return nil, fmt.Errorf("Value cannot be empty: %s", sRule)
			}
		}

		if _, ok := rules[tag]; !ok {
			rules[tag] = &rule{[]string{}, []string{}}
		}

		if allowed {
			rules[tag].allowedValues = values
		} else {
			rules[tag].deniedValues = values
		}
	}

	return rules, nil
}
