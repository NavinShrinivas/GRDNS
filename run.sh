echo "Needed Dependencies : redis "
echo "Which distro are u running me on?"
echo "1.Debian base distros [Ubuntu, pop, ...]"
echo "2.Arch base distro [EOS, Manjaro, ..]"
echo "3.MacOS (darwin)"
echo -n "Enter option :"
read option

case $option in
1)
    sudo systemctl start redis
    notify_problems "redis is not installed"
    ;;
2)
    sudo systemctl start redis-server
    notify_problems "redis is not installed"
    ;;
3)
    brew services start redis
    notify_problems "brew is not installed" "redis is not installed"
    ;;
*)
    echo "Invalid option"
    ;;
esac

go build .
notify_problems "No go toolchains installed"

./GRDNS
