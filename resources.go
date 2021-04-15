package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	MiB = 1024 * 1024
	KiB = 1024
)

type Metrics struct {
	CPU    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

type FloatMetrics struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type Container struct {
	Name    string   `json:"name"`
	Metrics *Metrics `json:"metrics"`
}

type Microservice struct {
	Timestamp  jsonTime     `json:"timestamp"`
	Name       string       `json:"name"`
	Containers []*Container `json:"containers"`
	Metrics    *Metrics     `json:"metrics"`
}

func (m *Microservice) ToJSONString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal json: %v", err)
	}
	return string(bytes), nil
}

type AverageMetrics struct {
	Timestamp jsonTime      `json:"timestamp"`
	Name      string        `json:"name"`
	Metrics   *FloatMetrics `json:"metrics"`
}

func (am *AverageMetrics) ToJSONString() (string, error) {
	bytes, err := json.Marshal(am)
	if err != nil {
		return "", fmt.Errorf("failed to marshal json: %v", err)
	}
	return string(bytes), nil
}

func NewContainer(name string, m *Metrics) *Container {
	c := &Container{Name: name, Metrics: m}
	return c
}

func NewMicroservice(name string, containers []*Container, timestamp time.Time) *Microservice {
	var jt jsonTime
	jt.Time = timestamp
	totalMetrics, err := calcTotalResource(containers)
	if err != nil {
		fmt.Printf("%s: %v", name, err)
	}
	m := &Microservice{Name: name, Containers: containers, Timestamp: jt, Metrics: totalMetrics}
	return m
}

func NewAverageMetrics(name string, timestamp time.Time, fm *FloatMetrics) *AverageMetrics {
	var jt jsonTime
	jt.Time = timestamp
	return &AverageMetrics{Name: name, Metrics: fm, Timestamp: jt}
}

func calcTotalResource(containers []*Container) (*Metrics, error) {
	var cpu int64
	var memory int64
	for _, container := range containers {
		if container.Metrics.CPU > math.MaxInt64-cpu {
			return nil, fmt.Errorf("detect overflow")
		}
		if container.Metrics.Memory > math.MaxInt64-memory {
			return nil, fmt.Errorf("detect overflow")
		}
		cpu += container.Metrics.CPU
		memory += container.Metrics.Memory
	}
	return &Metrics{CPU: cpu, Memory: memory}, nil
}
