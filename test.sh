#!/bin/bash

NO_OF_SERVERS=50

set -e

image=$(docker images ider-acceptance-test -q)
imageClient=$(docker images ider-client -q)
imageSorter=$(docker images ider-sorter -q)
dataFolder=$(pwd)/data

if [ -d $dataFolder ]; then
  echo $dataFolder
  sudo rm -rf $dataFolder
fi

if [ -z "$image" ]; then
    docker build -t ider-acceptance-test .
fi

if [ -z "$imageClient" ]; then
    docker build -t ider-client ./acceptance_tests/client
fi

if [ -z "$imageSorter" ]; then
    docker build -t ider-sorter ./acceptance_tests/sorter
fi

for i in $(seq 1 $NO_OF_SERVERS)
do
  port=$(printf "8%03g" $i)
  docker run --name iti$i -v $(pwd)/data:/data -p 127.0.0.1:$port:$port --expose=$port -e "SERVER_ID=$i" -e "SERVER_PORT=$port" -d ider-acceptance-test
done

docker run --name itiClient --link iti1 --link iti2 --link iti3 --link iti4 --link iti5 \
  --link iti6 --link iti7 --link iti8 --link iti9 --link iti10 \
  --link iti11 --link iti12 --link iti13 --link iti14 --link iti15 \
  --link iti16 --link iti17 --link iti18 --link iti19 --link iti20 \
  --link iti21 --link iti22 --link iti23 --link iti24 --link iti25 \
  --link iti26 --link iti27 --link iti28 --link iti29 --link iti30 \
  --link iti31 --link iti32 --link iti33 --link iti34 --link iti35 \
  --link iti36 --link iti37 --link iti38 --link iti39 --link iti40 \
  --link iti41 --link iti42 --link iti43 --link iti44 --link iti45 \
  --link iti46 --link iti47 --link iti48 --link iti49 --link iti50 \
-it ider-client

docker run --name itiSorter -m 1200M -v $(pwd)/data:/data -it ider-sorter





#docker run --name iti111 -v $(pwd)/data:/data -p 8080:8888 -e "SERVER_ID=111" -d ider-integration-test
