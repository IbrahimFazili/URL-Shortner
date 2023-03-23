USAGE="Usage: $0 IP1 IP2 IP3 ..."

if [ "$#" == "0" ]; then
        echo "$USAGE"
        exit 1
fi
while (( "$#" )); do
        sshpass -p hhhhiotwwg ssh -t student@$1 "sudo rm -rf data/"
        echo "destroyed data/cassandra folder in $1"
        shift
done
