#! /bin/bash

DUMPFILE="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates.json"
DUMPFILE_7DAY="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates7days.json"

set -eu

echo "Creating all.bin"
./systemCoord.exe -o all.bin $DUMPFILE
echo "Creating 7day.bin"
./systemCoord.exe -o 7day.bin $DUMPFILE_7DAY

#echo "Creating all system image"
#./imaging.exe -i all.bin -p xz -o xz.png
#./imaging.exe -i all.bin -p xz -ht opaque -o xz_opaque.png
#./imaging.exe -i all.bin -p xy -o xy.png
#./imaging.exe -i all.bin -p xy -ht opaque -o xy_opaque.png
#./imaging.exe -i all.bin -p zy -o zy.png
#./imaging.exe -i all.bin -p zy -ht opaque -o zy_opaque.png
#
#echo "Creating updated system image"
#./imaging.exe -i 7day.bin -p xz -o xz_7day.png
#./imaging.exe -i 7day.bin -p xy -o xy_7day.png
#./imaging.exe -i 7day.bin -p zy -o zy_7day.png
