// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package tracer

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"strings"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/internal/globalconfig"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/log"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/version"
)

const (
	unknown = "unknown"
)

type startupInfo struct {
	Date                  string                 `json:"date"`    // ISO 8601 date and time of start
	OSName                string                 `json:"os_name"` // Windows, Darwin, Debian, etc.
	OSVersion             string                 `json:"os_version"`
	Version               string                 `json:"version"`              // Tracer version
	Lang                  string                 `json:"lang"`                 // "Go"
	LangVersion           string                 `json:"lang_version"`         // Go version, e.g. go1.13
	Env                   string                 `json:"env"`                  // Tracer env
	Service               string                 `json:"service"`              // Tracer Service
	AgentURL              string                 `json:"agent_url"`            // The address of the agent
	AgentError            error                  `json:"agent_error"`          // Any error that occurred trying to connect to agent
	Debug                 bool                   `json:"debug"`                // Whether debug mode is enabled
	AnalyticsEnabled      bool                   `json:"analytics_enabled"`    // True if there is a global analytics rate set
	SampleRate            float64                `json:"sample_rate"`          // The default sampling rate for the priority sampler
	SamplingRules         []SamplingRule         `json:"sampling_rules"`       // Rules used by the rules sampler
	SamplingRulesError    error                  `json:"sampling_rules_error"` // Any errors that occurred while parsing sampling rules
	Tags                  map[string]interface{} `json:"tags"`                 // Global tags
	RuntimeMetricsEnabled bool                   `json:"runtime_metrics_enabled"`
	HealthMetricsEnabled  bool                   `json:"health_metrics_enabled"`
	ApplicationVersion    string                 `json:"dd_version"`     // Version of the user's application
	Architecture          string                 `json:"architecture"`   // Architecture of host machine
	GlobalService         string                 `json:"global_service"` // Global service string. If not-nil should be same as Service. (#614)
}

func agentReachable(t *tracer) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/v0.4/traces", resolveAddr(t.config.agentAddr)), strings.NewReader("[]"))
	if err != nil {
		return fmt.Errorf("cannot create http request: %v", err)
	}

	req.Header.Set(traceCountHeader, "0")
	req.Header.Set("Content-Length", "2")
	response, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	if code := response.StatusCode; code != 200 {
		// error, check the body for context information and
		// return a nice error.
		msg := make([]byte, 1000)
		n, _ := response.Body.Read(msg)
		response.Body.Close()
		txt := http.StatusText(code)
		if n > 0 {
			return fmt.Errorf("%s (Status: %s)", msg[:n], txt)
		}
		return fmt.Errorf("%s", txt)
	}
	return nil
}

func newStartupInfo(t *tracer) *startupInfo {
	if !envBool("DD_TRACE_STARTUP_LOGS", true) {
		return &startupInfo{}
	}
	return &startupInfo{
		Date:                  time.Now().Format(time.RFC3339),
		OSName:                osName(),
		OSVersion:             osVersion(),
		Version:               version.Tag,
		Lang:                  "Go",
		LangVersion:           runtime.Version(),
		Env:                   t.config.env,
		Service:               t.config.serviceName,
		AgentURL:              t.config.agentAddr,
		AgentError:            agentReachable(t),
		Debug:                 t.config.debug,
		AnalyticsEnabled:      !math.IsNaN(globalconfig.AnalyticsRate()),
		SampleRate:            t.prioritySampling.defaultRate,
		SamplingRules:         t.rulesSampling.rules,
		Tags:                  t.globalTags,
		RuntimeMetricsEnabled: t.config.runtimeMetrics,
		HealthMetricsEnabled:  t.config.runtimeMetrics,
		ApplicationVersion:    t.config.version,
		Architecture:          runtime.GOARCH,
		GlobalService:         globalconfig.ServiceName(),
	}
}

func logStartup(info *startupInfo) {
	if !boolEnv("DD_TRACE_STARTUP_LOGS", true) {
		return
	}
	bs, err := json.Marshal(info)
	if err != nil {
		log.Error("Failed to serialize json for startup log: %#v\n", info)
		return
	}
	log.Info("Startup: %s\n", string(bs))
}
