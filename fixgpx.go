package fixgpx

import (
	"fmt"
	"time"
	"os"
	"regexp"
	"strings"
	"bufio"
)

// create parsing regexps
var timeTagRE = regexp.MustCompile(`<time>(.+)<\/time>`)
var metaStartTagRE = regexp.MustCompile(`^\s*<metadata>`)
var metaEndTagRE = regexp.MustCompile(`^\s*</metadata>`)
var trkStartTagRE = regexp.MustCompile(`^\s*<trk>`)
var trkEndTagRE = regexp.MustCompile(`^\s*</trk>`)
	

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

// load the entire input file into a string slice
// allow multiple passes
func loadGPXIn(gpxin string) ([]string, error) {

	// buffer for lines in the file
	lineBuf := []string{}

	// open the gpx file
	fdin, finerr := os.Open(gpxin)
	if finerr != nil {
		return lineBuf, fmt.Errorf("fixgpx: Error opening gpx file %s: %v\n", gpxin, finerr)
	}
	defer fdin.Close()

	// read the file into a slice
	scanner := bufio.NewScanner(fdin)

	for scanner.Scan() {
		lineBuf = append(lineBuf, scanner.Text())
	}

	return lineBuf, nil
}

// get the time delta
// The time in the first trkpt should match the time in the metadata tag
// return the delta between the metadata tag and the first trkpt as a unix epoch int64
func GetTimeDelta(lines []string) (int64, error) {

	// result variables
	metaEpoch := int64(0)
	firstTrkEpoch := int64(0)

	//state variables
	inMeta := false
	inTrk := false

	for _, line := range lines {

		// entering metadata tag
		if metaStartRm := metaStartTagRE.FindStringSubmatch(line) ; metaStartRm != nil {
			inMeta = true
		}

		// leaving metadata tag
		if metaEndRm := metaEndTagRE.FindStringSubmatch(line) ; metaEndRm != nil {
			inMeta = false
		}

		// entering trk tag
		if trkStartRm := trkStartTagRE.FindStringSubmatch(line) ; trkStartRm != nil {
			inTrk = true
		}

		// leaving trk tag
		if trkEndRm := trkEndTagRE.FindStringSubmatch(line) ; trkEndRm != nil {
			inTrk = false
		}

		// capture time tags
		if timeTagRm := timeTagRE.FindStringSubmatch(line) ; timeTagRm != nil {
			epoch, err := iso2epoch(timeTagRm[1])
			if err != nil {
				return 0, fmt.Errorf("GetTimeDelta: Error: %v", err)
			}

			if inMeta {
				metaEpoch = epoch
			} else if inTrk {
				firstTrkEpoch = epoch

				// just need the first one
				if metaEpoch != 0 {
					return firstTrkEpoch - metaEpoch, nil
				} else {
					return 0, fmt.Errorf("GetTimeDelta: Error: trk time found while 0 meta epoch")
				}

			} else {
				return 0, fmt.Errorf("GetTimeDelta: Error: time tag in unexpected context")
			}
		}
	}
	// return to keep the compiler happy, shouldn't get here
	return 0, fmt.Errorf("getTimeDelta: Error: reached final return and shouldn't have")
}

// correct a time string given an epoch delta
func fixTimeDelta(origTime string, delta int64) (string, error) {
	origEpoch, err := iso2epoch(origTime)
	if err != nil {
		return "", fmt.Errorf("fixTimeDelta: Error converting %s: %v", origTime, err)
	}
	newEpoch := origEpoch - delta
	return epoch2iso(newEpoch), nil
}

// write the fixed gpx file
func WriteFixedGPX(gpxout string, lines []string, delta int64) error {

	// open the output file
	fdout, fderr := os.Create(gpxout)
	if fderr != nil {
		return fmt.Errorf("WriteFixedGPX: Error opening gpx output file %s for write: %v\n", gpxout, fderr)
	}
	defer fdout.Close()

	// state variable
	inTrk := false
	
	// write the new version of the file
	for _, line := range lines {

		fixedLine := line // default to no change

		// entering trk tag
		if trkStartRm := trkStartTagRE.FindStringSubmatch(line) ; trkStartRm != nil {
			inTrk = true
		}

		// leaving trk tag
		if trkEndRm := trkEndTagRE.FindStringSubmatch(line) ; trkEndRm != nil {
			inTrk = false
		}

		// fix all time tags inside the <trk> tag
		if inTrk {
			match := timeTagRE.FindStringSubmatch(line)
			if match != nil{
				origTime := match[1]
				fixedTime, fterr := fixTimeDelta(origTime, delta)
				if fterr != nil {
					return fmt.Errorf("WriteFixedTCX: Error: %v\n", fterr)
				}
				fixedLine = strings.Replace(line, origTime, fixedTime, 1)
			}
		}
		
		_, perr := fmt.Fprintln(fdout, fixedLine)
		if perr != nil {
			return fmt.Errorf("WriteFixedTCX: Error writing line %s: %v\n", fixedLine, perr)
		}
	}
	return nil
}
