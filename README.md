# fixgpx

This program repairs corrupted gpx files for one specific failure
mode.  It solves the common case in which the timestamps jump into the
future vs the actual start time and subsequently track normally.  This
is the only corruption case I have seen on my device.  The package
could be used to find and fix other types of corruption given an example.

# Usage
% fixgpx -gpxin corrupted-gpx -gpxout repaired-gpx


