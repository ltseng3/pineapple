timeout 180s go run client/client.go -maddr="10.10.1.1" -writes=0.375 -rmws=0.25 -c=0 -T=3
python3 client_metrics.py
