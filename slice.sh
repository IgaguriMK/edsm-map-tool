#! /bin/bash

WIDTH=3000

COUNT=0
TOP=$WIDTH

mkdir -p slices

while [ $TOP -gt -$WIDTH ]; do
	if [ $(($TOP % 200)) -eq 0 ]; then
		echo mid slice $BOTTOM
		./transform -o tmp.bin :cut-y $(($TOP - 400)) $TOP
	fi

	BOTTOM=$(($TOP - 40))
	./transform -i tmp.bin :cut-y $BOTTOM $TOP :add -42300 0 -16900 :add 40600 0 65700
	./imaging -o slices/`printf %05d $COUNT`.png -i trans.bin -hs 3 -ht colorful_noback -multof 4
	COUNT=$(($COUNT + 1))
	TOP=$(($TOP - 10))
done

rm tmp.bin
rm trans.bin

ffmpeg -r 30 -i slices/%05d.png -vcodec libx264 -pix_fmt yuv420p -crf 20 slices.mp4
