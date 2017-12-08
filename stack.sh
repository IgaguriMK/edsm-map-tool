#! /bin/bash

WIDTH=3000

COUNT=0
TOP=-$WIDTH

mkdir -p stack || true

while [ $TOP -lt $WIDTH ]; do
	if [ $(($TOP % 200)) -eq 0 ]; then
		echo mid slice $TOP
		./transform -o tmp.bin :cut-y -$WIDTH $(($TOP + 400))
	fi

	./transform -i tmp.bin :cut-y -$WIDTH $TOP :add -42300 0 -16900 :add 40600 0 65700
	./imaging -o stack/`printf %05d $COUNT`.png -i trans.bin -hs 0.1 -ht colorful_noback
	COUNT=$(($COUNT + 1))
	TOP=$(($TOP + 10))
done

rm tmp.bin
rm trans.bin
