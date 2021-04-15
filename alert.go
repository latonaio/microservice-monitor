package main

import "fmt"

type Threshold struct {
	threshold float64
}

func (t Threshold) IsOver(value float64) bool {
	return t.threshold < value
}

type MetricsWindow struct {
	size            int
	cpuThreshold    *Threshold
	memoryThreshold *Threshold
	window          []Metrics
}

func (mw *MetricsWindow) Push(m Metrics) {
	mw.window = append(mw.window, m)
	if len(mw.window) > mw.size {
		mw.window = mw.window[1:]
	}
}

func (mw *MetricsWindow) CalcAverage() *FloatMetrics {
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
	return &FloatMetrics{CPU: cpuAvg, Memory: memoryAvg}
}

func (mw *MetricsWindow) Alert() (*FloatMetrics, error) {
	fm := mw.CalcAverage()
	var errStr []string
	if mw.cpuThreshold.IsOver(fm.CPU) {
		errStr = append(errStr, fmt.Sprintf("alert: CPU over %f < %f", mw.cpuThreshold.threshold, fm.CPU))
	}
	if mw.memoryThreshold.IsOver(fm.Memory) {
		errStr = append(errStr, fmt.Sprintf("alert: Memory over %f < %f", mw.memoryThreshold.threshold, fm.Memory))
	}
	if len(errStr) > 0 {
		return fm, fmt.Errorf("%#v", errStr)
	}
	return nil, nil
}

func NewMetricsWindow(size int, cpuThreshold float64, memoryThreshold float64) *MetricsWindow {
	return &MetricsWindow{
		size:            size,
		cpuThreshold:    &Threshold{threshold: cpuThreshold},
		memoryThreshold: &Threshold{threshold: memoryThreshold},
		window:          []Metrics{},
	}
}
