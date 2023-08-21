timeout 180s bin/client -maddr="10.10.1.1" -writes=.50 -rmws=0 -c=0 -T=80
python3 client_metrics.py
