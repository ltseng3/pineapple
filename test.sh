timeout 180s go run client/client.go -maddr="10.10.1.1" -writes=.5 -rmws=0 -c=0 -T=1
python3 client_metrics.py
