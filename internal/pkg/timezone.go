package pkg

import "time"

func ArgentinaLocation() *time.Location {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return time.FixedZone("ART", -3*60*60)
	}

	return loc
}
