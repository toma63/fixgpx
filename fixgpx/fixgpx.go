package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/toma63/fixgpx"
)

func main () {

	gpxin := flag.String("gpxin", "", "gpx file needing repair, read as an input")
	gpxout := flag.String("gpxout", "", "repaired version of the gpx file, written as an output")
	flag.Parse()

	if (*gpxin == "") || (*gpxout == "") {
		fmt.Println("fixgpx: error, input and/or output files not specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// load the input file into a string slice - allow multiple passes
	lineBuf, rerr := fixgpx.LoadGPXIn(*gpxin)
	if rerr != nil {
		fmt.Printf("fixgpx: error loading input file: %v", rerr)
	}

	// compute the time delta
	delta, delerr := fixgpx.GetTimeDelta(lineBuf)
	if delerr != nil {
		fmt.Printf("fixgpx: error computing time delta: %v", delerr)
	}

	fmt.Printf("Time delta is: %d\n", delta)

	// write the repaired file
	fterr := fixgpx.WriteFixedGPX(*gpxout, lineBuf, delta)
	if fterr != nil {
		fmt.Printf("fixgpx: error writing repaired file: %v\n", fterr)
		os.Exit(1)
	}
	os.Exit(0)
}
