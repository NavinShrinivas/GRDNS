echo "Make sure go toolchains are installed properly"
echo "also make sure you have redis server installed"
sudo systemctl restart redis
redis-cli FLUSHALL
sudo go run .
