timeout 180s go run client/client.go -maddr="10.10.1.1" -writes=0.425 -rmws=0.15 -c=0 -T=300
python3 client_metrics.py
