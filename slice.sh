#! /bin/bash

WIDTH=3000

COUNT=0
BOTTOM=-$WIDTH

mkdir -p slices

while [ $BOTTOM -lt $WIDTH ]; do
	if [ $(($BOTTOM % 200)) -eq 0 ]; then
		echo mid slice $BOTTOM
		./transform -o tmp.bin :cut-y $BOTTOM $(($BOTTOM + 400))
	fi

	TOP=$(($BOTTOM + 40))
	./transform -i tmp.bin :cut-y $BOTTOM $TOP :add -42300 0 -16900 :add 40600 0 65700
	./imaging -o slices/`printf %05d $COUNT`.png -i trans.bin -hs 3 -ht colorful_noback
	COUNT=$(($COUNT + 1))
	BOTTOM=$(($BOTTOM + 10))
done

rm tmp.bin
rm trans.bin
