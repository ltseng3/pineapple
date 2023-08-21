timeout 180s bin/client -saddr=$1 -sport=$2 -serverID=$3 -writes=.50 -rmws=0 -c=0 -T=80
python3 client_metrics.py
