package fixgpx

import "testing"

func TestEpoch2iso(t *testing.T) {

	testEpoch := int64(1488727327)
	goldenIso := "2017-03-05T15:22:07.000Z"

	result := epoch2iso(testEpoch)

	if result != goldenIso {
		t.Errorf("epoch2iso result %s should be %s", result, goldenIso)
	}
	
}

func TestIso2epoch(t *testing.T) {

	testIso := "2017-03-05T15:22:07.000Z"
	goldenEpoch := int64(1488727327)

	result, err := iso2epoch(testIso)
	if err != nil {
		t.Errorf("iso2epoch unexpected error on %s: %v", testIso, err)
	}
	if result != goldenEpoch {
		t.Errorf("iso2epoch result %d should be %d", result, goldenEpoch)
	}

}

func TestLoadGPXIn(t *testing.T) {

	testFile := "./testdata/sample.gpx"

	lineBufUT, err := loadGPXIn(testFile)
	if err != nil {
		t.Errorf("LoadGPXIn returned an unexpected error opening %s: %v", testFile, err)
	}

	// check an assortment of line including first and last
	// map of line#, expected string
	testLines := map[int]string {
		0: `<?xml version="1.0" encoding="UTF-8"?>`,
		len(lineBufUT) - 1: "</gpx>",
		10: "    <time>2017-02-22T21:28:02.000Z</time>",
		96: `      <trkpt lat="37.401199601590633392333984375" lon="-121.93786279298365116119384765625">`,
		5411: "            <ns3:hr>146</ns3:hr>",
	}

	for idx, expected := range testLines {
		if lineBufUT[idx] != expected {
			t.Errorf("LoadGPXIn line %d does not match expected value: %s, actual: %s", idx, expected, lineBufUT[idx])
		}
	}
	
}

func TestGetTimeDelta(t *testing.T) {

	testLines := []string{"wozzle",
		"gozzle snozzle\n",
		"<sddad>schmoo<dsaUTYTY>",
		"     metadata",
		" <metadata>",
		"bazzle",
		"    <time>2017-02-22T21:28:02.000Z</time>",
		"</metadata>",
		"sdasdad",
		"<><><><>",
		"<trk>",
		"        <time>2017-03-01T23:15:45.000Z</time>",
		"</trk>",
		`<trkpt lat="37.40106205455958843231201171875" lon="-121.93716533482074737548828125">`,
	}

	goldenDelta := int64(611263)

	res, err := GetTimeDelta(testLines)

	if err != nil {
		t.Errorf("GetTimeDelta returned an unexpected error: %v", err)
	}

	if res != goldenDelta {
		t.Errorf("getTimeDelta return %d instead of %d as expected", res, goldenDelta)
	}
}
