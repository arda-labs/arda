package sla

import (
	"time"
)

// CalculateSLAExpiry calculates the expiration time based on start time and duration in hours.
// In a real banking system, this would account for business hours and holidays.
// For now, it implements a basic 24/7 calculation.
func CalculateSLAExpiry(startTime time.Time, durationHours int) time.Time {
	return startTime.Add(time.Duration(durationHours) * time.Hour)
}

// IsOverdue checks if the current time is past the SLA expiration time.
func IsOverdue(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}
