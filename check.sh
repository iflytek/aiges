#!/bin/bash

words="iat tts dts ent aue auf language domain accent prs acp saddr dsc app maxcnt scn eos ptt clg dwa vascn rsh wbest scence svad sch"
dir="*"
if [ ! -z "$1" ]; then
	echo "set directory to $1"
	dir=$1
fi

for word in $words; do
	if [[ -z "$IGNORE" ]]; then
		grep -E  "[^a-z|0-9]$word[^a-z]" $dir -rn |grep  -E  "[^a-z|0-9]$word[^a-z]" --color
	else
		#echo $IGNORE	
		grep -E  "[^a-z|0-9]$word[^a-z]" $dir -rn | grep -v -E $IGNORE |grep  -E  "[^a-z|0-9]$word[^a-z]" --color
	fi
done
