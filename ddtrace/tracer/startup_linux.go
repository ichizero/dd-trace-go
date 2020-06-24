// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package tracer

import (
	"bufio"
	"os"
	"strings"
)

func osName() string {
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		name := "Linux (Unknown Distribution)"
		for s.Scan() {
			parts := strings.SplitN(s.Text(), "=", 2)
			switch parts[0] {
			case "Name":
				name = strings.Trim(parts[1], "\"")
			}
		}
		return name
	}
	return "Linux (Unknown Distribution)"
}

func osVersion() string {
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		version := "(Unknown Version)"
		for s.Scan() {
			parts := strings.SplitN(s.Text(), "=", 2)
			switch parts[0] {
			case "VERSION":
				version = strings.Trim(parts[1], "\"")
			case "VERSION_ID":
				if version == "" {
					version = strings.Trim(parts[1], "\"")
				}
			}
		}
		return version
	}
	return "(Unknown Version)"
}
