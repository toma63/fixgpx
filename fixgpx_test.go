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
