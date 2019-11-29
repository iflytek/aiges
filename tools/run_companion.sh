#!/bin/bash

docker run -itd --net host --name companion -v /log/athena/componion/:/log/ --rm littlescw00/companion:latest sh watchdog.sh $*
