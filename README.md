# fixgpx

This program repairs corrupted gpx files for one specific failure
mode.  It solves the common case in which the timestamps jump into the
future vs the actual start time and subsequently track normally.  This
is the only corruption case I have seen on my device.  The package
could be used to find and fix other types of corruption given an example.

# Build Instructions
Fixgpx can be built with the standard go tools.  With the GOPATH environment variable 
pointing to your workspace, run the folloing commands.

% go get github.com/toma63fixgpx

% cd $GOPATH/src/github.com/toma63/fixgpx

## run unit tests

% go test

## With tests passing, build and install the library

% go build

% go install

## build the spplication

% cd fixgpx

% go build

% go install

The application will now be available in $GOPATH/bin.

I have built it successfully on MacOSX Mavericks, Ubuntu Mate on Rapsberry pi3, and Windows 10.

# Usage
% fixgpx -gpxin corrupted-gpx -gpxout repaired-gpx


