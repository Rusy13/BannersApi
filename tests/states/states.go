package states

import "time"

var (
	Banner1TagIDs  = []int{1, 2, 3}
	Banner2TagIDs  = []int{4, 5, 6}
	Banner1Content = map[string]interface{}{
		"title": "New Banner",
		"text":  "Some text",
		"url":   "http://example.com",
	}
	Banner2Content = map[string]interface{}{
		"title": "Another Banner",
		"text":  "More text",
		"url":   "http://example2.com",
	}
	Banner1CreatedAt = time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)
	Banner2CreatedAt = time.Date(2024, time.April, 15, 10, 31, 0, 0, time.UTC)

	Banner1UpdatedAt = time.Date(2024, time.April, 15, 10, 32, 0, 0, time.UTC)
	Banner2UpdatedAt = time.Date(2024, time.April, 15, 10, 33, 0, 0, time.UTC)
)

const (
	Banner1ID = 1
	Banner2ID = 2

	Banner1FeatureID = 1
	Banner2FeatureID = 1

	Banner1IsActive = true
	Banner2IsActive = false
)
