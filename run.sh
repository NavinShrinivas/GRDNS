echo "Needed Dependencies : redis "
echo "Which distro are u running me on?"
echo "1. Debian base distros [Ubunut, pop, ...]"
echo "2. Arch base distro [EOS, Manjaro, ..]"
echo "3. MacOS (Darwin)"
echo -n "Enter option :"
read option

notify_problems() {
    if [ $? -ne 0 ]; then
        echo "oh oh, something went wrong, possible causes : "
        for cause in "$@"; do
            echo $cause
        done
        exit 1
    fi
}

case $option in
1)
    sudo systemctl start redis &>/dev/null
    ;;
2)
    sudo systemctl start redis-serve &>/dev/null
    ;;
3)
    brew --version &>/dev/null
    notify_problems "brew is not installed"
    brew services run redis &>/dev/null
    notify_problems "redis is not installed"
    ;;
esac

redis-cli FLUSHALL

go build .
notify_problems "No go toochains installed"

./GRDNS
