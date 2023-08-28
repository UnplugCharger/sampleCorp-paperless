package utils

import "time"

func ShortenTimestamp(timestamp time.Time) (string, error) {
	// Load the East Africa Time (EAT) time zone
	location, err := time.LoadLocation("Africa/Nairobi") // Nairobi is in East Africa Time Zone
	if err != nil {
		return "", err
	}

	// Change the time zone of the timestamp object to EAT
	timestampInEAT := timestamp.In(location)

	// Define the desired output format
	outputFormat := "2006-01-02 15:04:05"

	// Format the time object in the desired way
	shortenedTimestamp := timestampInEAT.Format(outputFormat)

	return shortenedTimestamp, nil
}
