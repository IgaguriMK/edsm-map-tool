#! /bin/bash

DUMPFILE="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates.json"
DUMPFILE_7DAY="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates7days.json"

set -eu

echo "Creating all.bin"
./systemCoord.exe -o all.bin $DUMPFILE
echo "Creating 7day.bin"
./systemCoord.exe -o 7day.bin $DUMPFILE_7DAY
