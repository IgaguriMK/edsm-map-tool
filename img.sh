#! /bin/bash

DUMPFILE="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates.json"
DUMPFILE_7DAY="C:\Users\Igaguri\Documents\EliteDangerous\EDTools\systemsWithCoordinates7days.json"

set -eu

echo "Creating all.bin"
./systemCoord.exe -o all.bin $DUMPFILE &
echo "Creating 7day.bin"
./systemCoord.exe -o 7day.bin $DUMPFILE_7DAY %

wait

echo "Creating all system image"
./imaging.exe -i all.bin -p xz -o xz_all.png &
./imaging.exe -i all.bin -p xz -ht colorful -o xz_all_trans.png &
wait
./imaging.exe -i all.bin -p xy -o xy_all.png &
./imaging.exe -i all.bin -p xy -ht colorful -o xy_all_trans.png &
wait
./imaging.exe -i all.bin -p zy -o zy_all.png &
./imaging.exe -i all.bin -p zy -ht colorful -o zy_all_trans.png &
wait

echo "Creating updated system image"
./imaging.exe -i 7day.bin -p xz -o xz_7day.png &
./imaging.exe -i 7day.bin -p xy -o xy_7day.png &
./imaging.exe -i 7day.bin -p zy -o zy_7day.png &
wait
