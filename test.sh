timeout 180s go run client/client.go -maddr="10.10.1.1" -writes=0.4 -rmws=0.20 -c=0 -T=15
python3 client_metrics.py
