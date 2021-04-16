package main

import (
	"fmt"
	"strings"
)

type Threshold struct {
	threshold float64
}

func (t Threshold) IsDefined() bool {
	return t.threshold > 0
}

func (t Threshold) IsOver(value float64) bool {
	return t.threshold < value
}

type Alert struct {
	Name            string
	CPUThreshold    Threshold
	MemoryThreshold Threshold
	Window          MetricsWindow
	Status          bool
}

type MetricsWindow struct {
	size   int
	window []Metrics
}

func (mw *MetricsWindow) Push(m Metrics) {
	mw.window = append(mw.window, m)
	if len(mw.window) > mw.size {
		mw.window = mw.window[1:]
	}
}

func (mw MetricsWindow) CalcAverage() FloatMetrics {
	var cpuAvg float64
	var memoryAvg float64
	for i, m := range mw.window {
		if i == 0 {
			cpuAvg = float64(m.CPU)
			memoryAvg = float64(m.Memory)
		}
		cpuAvg = cpuAvg + (float64(m.CPU)-cpuAvg)/2.0
		memoryAvg = memoryAvg + (float64(m.Memory)-memoryAvg)/2.0
	}
	return FloatMetrics{CPU: cpuAvg, Memory: memoryAvg}
}

func (a Alert) IsRaised() bool {
	return a.Status
}

func (a *Alert) Alert(now jsonTime, env Env, as AlertSetting) (FloatMetrics, error) {
	fm := a.Window.CalcAverage()
	var errStr []string
	if a.CPUThreshold.IsDefined() && a.CPUThreshold.IsOver(fm.CPU) {
		errStr = append(errStr, fmt.Sprintf("alert: CPU over %f < %f", a.CPUThreshold.threshold, fm.CPU))
	}
	if a.MemoryThreshold.IsDefined() && a.MemoryThreshold.IsOver(fm.Memory) {
		errStr = append(errStr, fmt.Sprintf("alert: Memory over %f < %f", a.MemoryThreshold.threshold, fm.Memory))
	}
	if len(errStr) > 0 {
		if !a.IsRaised() {
			Notify(
				MessageContents{
					MicroserviceName: a.Name,
					DeviceName:       as.Device.Name,
					IPv4:             as.Device.Addr,
					AlertTime:        now,
					AlertLog:         fmt.Sprintf("%#v", errStr),
				},
				env,
			)
		}
		a.Status = true
		return fm, fmt.Errorf("%#v", errStr)
	}
	a.Status = false
	return FloatMetrics{}, nil
}

func NewMetricsWindow(size int) MetricsWindow {
	return MetricsWindow{
		size:   size,
		window: []Metrics{},
	}
}

func NewAlert(name string, size int, cpu float64, memory float64) *Alert {
	return &Alert{
		Name:            name,
		CPUThreshold:    Threshold{threshold: cpu},
		MemoryThreshold: Threshold{threshold: memory},
		Window:          NewMetricsWindow(size),
	}
}

func MakeAlert(as AlertSetting, env Env) []*Alert {
	var alerts []*Alert
	for _, a := range as.Alert {
		alerts = append(alerts, NewAlert(a.Name, env.WindowSize, a.Threshold.CPU, a.Threshold.Memory))
	}
	return alerts
}

func FindAlert(alerts []*Alert, name string) (*Alert, bool) {
	for _, alert := range alerts {
		if strings.HasPrefix(name, alert.Name) {
			return alert, true
		}
	}
	return nil, false
}
