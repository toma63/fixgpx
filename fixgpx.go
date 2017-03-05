package fixgpx

import (
	"fmt"
	"time"
	//"os"
	"regexp"
	//"strings"
	//"bufio"
)

// create parsing regexps
var timeTagRe = regexp.MustCompile(`<time>(.+)<\/time>`)

// translate an iso-8601 utc time to unix epoch int64
//  looks like 2006-01-02T15:04:05.000Z
func iso2epoch(iso string) (int64, error) {

	// get rid of 'Z'
	iso = iso[:len(iso) - 1]

	// parse the time, default is utc
	to, err := time.Parse("2006-01-02T15:04:05.000", iso)
	if err != nil {
		return 0, fmt.Errorf("iso2epoch: error parsing time %s: %v", iso, err)
	}
	return to.Unix(), nil
}

// translate an int64 unix epoch time to an iso-8601 string
func epoch2iso(epoch int64) string {
	to := time.Unix(epoch, 0).UTC() // default is local
	iso := to.Format("2006-01-02T15:04:05.000")
	return iso + "Z" // gpx includes a trailing Z for zulu (utc)
}

