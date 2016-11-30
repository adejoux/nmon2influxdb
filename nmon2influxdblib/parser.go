// nmon2influxdb
// author: adejoux@djouxtech.net

package nmon2influxdblib

import (
	"fmt"
	"regexp"
)

//Tags array
type Tags []Tag

// Tag is a struct to store additional tags
type Tag struct {
	Name   string
	Value  string
	Regexp *regexp.Regexp `toml:",skip"`
}

// TagParsers access struct : TagParsers[mesurement][tag name]
type TagParsers map[string]map[string]Tags

// ParseInputs process user inputs and compile regular expressions
func ParseInputs(inputs Inputs) TagParsers {
	tagParsers := make(TagParsers)

	for _, input := range inputs {
		tagRegexp, RegCompErr := regexp.Compile(input.Match)
		if RegCompErr != nil {
			fmt.Printf("could not compile Config Input match parameter %s\n", input.Match)
		}

		// Intialize if empty struct
		if _, ok := tagParsers[input.Measurement]; !ok {
			tagParsers[input.Measurement] = make(map[string]Tags)
		}

		tagger := tagParsers[input.Measurement][input.Name]

		for _, tag := range input.Tags {
			tag.Regexp = tagRegexp
			tagger = append(tagger, tag)
		}

		tagParsers[input.Measurement][input.Name] = tagger

	}
	return tagParsers
}
