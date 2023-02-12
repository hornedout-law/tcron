#!/bin/bash

crontab -l; echo "$1 rm -r $2 ;crontab -l|grep -v $2|sort -u -|crontab - \\n" 2>/dev/null | sort -| crontab -