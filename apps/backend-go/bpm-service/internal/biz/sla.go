package biz

import (
	"time"
)

type SLAPriority string

const (
	SLANormal  SLAPriority = "ON_TIME"
	SLAWarning SLAPriority = "WARNING"
	SLAOverdue SLAPriority = "OVERDUE"
)

type SLAConfig struct {
	ID           string
	ProcessName  string
	StepName     string
	DurationHours int
	WarningPercent int // e.g., 80% of duration
}

type SLACalculator struct {
	// Add config repo here
}

func NewSLACalculator() *SLACalculator {
	return &SLACalculator{}
}

func (c *SLACalculator) CalculateStatus(startTime time.Time, config SLAConfig) (SLAPriority, time.Duration) {
	deadline := startTime.Add(time.Duration(config.DurationHours) * time.Hour)
	now := time.Now()

	if now.After(deadline) {
		return SLAOverdue, now.Sub(deadline)
	}

	warningDuration := time.Duration(float64(config.DurationHours*int(time.Hour)) * (float64(config.WarningPercent) / 100.0))
	warningTime := startTime.Add(warningDuration)

	if now.After(warningTime) {
		return SLAWarning, deadline.Sub(now)
	}

	return SLANormal, deadline.Sub(now)
}
