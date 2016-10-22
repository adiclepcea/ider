NO_OF_SERVERS=50
set +e
for i in $(seq 1 $NO_OF_SERVERS)
do
  port=$(printf "8%03g" $i)
  curl http://localhost:$port/stop
  sleep 1
  docker rm -f "iti$i"
done

docker rm -f itiClient
docker rm -f itiSorter

docker rmi ider-client
docker rmi ider-sorter
docker rmi ider-acceptance-test
set -e
