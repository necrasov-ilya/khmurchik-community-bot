package timeutil

import "time"

func LoadLocation() (*time.Location, error) {
	return time.LoadLocation("Europe/Minsk")
}
