echo "Needed Dependencies : redis "
echo "Which distro are u running me on?"
echo "1. Debian base distros [Ubunut, pop, ...]"
echo "2. Arch base distro [EOS, Manjaro, ..]"
echo "3. MacOS (Darwin)"
echo -n "Enter option :"
read option

notify_redis_absence() {
    if [ $? -ne 0 ]; then
        echo "oh oh, something went wrong, possible causes : "
        echo "No instance of redis is installed in this system"
        exit 1
    fi
}

case $option in
1)
    sudo systemctl start redis &>/dev/null
    notify_redis_absence

    ;;
2)
    sudo systemctl start redis-serve &>/dev/null
    notify_redis_absence
    ;;
3)
    brew --version &>/dev/null
    if [ $? -ne 0 ]; then
        echo "Requires brew to start redis"
        echo "Please make sure you have brew installed"
        echo "Then run 'brew install redis'"
        exit 1
    fi
    brew services run redis &>/dev/null
    notify_redis_absence
    ;;
esac

redis-cli FLUSHALL

go build .
if [ $? -ne "0" ]; then
    echo "oh oh, something went wrong, possible causes : "
    echo "No go toochains installed"
    exit 1
fi

sudo ./GRDNS
