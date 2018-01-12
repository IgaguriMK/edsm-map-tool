#! /bin/bash

./transform.exe -i all.bin -o test.bin :cut-x -500 500
./imaging.exe -i test.bin -o test_zy.png -p zy -ht opaque
