timeout 180s go run client/client.go -maddr="10.10.1.1" -writes=0.505 -rmws=0 -c=25 -T=80
python3 client_metrics.py
