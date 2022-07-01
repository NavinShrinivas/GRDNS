echo "Needed Dependencies : redis "
echo "Which distro are u running me on?"
echo "1.Debian base distros [Ubunut, pop, ...]"
echo "2.Arch base distro [EOS, Manjaro, ..]"
echo -n "Enter option :"
read option

if [ $option -ne "1" ];then
    sudo systemctl start redis
   if [ $? -ne "0" ];then
        echo "oh oh, something went wrong, possible causes : "
        echo "No instance of redis is installed in this system"
        exit
    fi
    redis-cli FLUSHALL
    go build .
    if [ $? -ne "0" ];then
        echo "oh oh, something went wrong, possible causes : "
        echo "No go toochains installed"
        exit
    fi
    chmod +x GRDNS
    sudo ./GRDNS
else
    sudo systemctl start redis-server
    if [ $? -ne "0" ];then
        echo "oh oh, something went wrong, possible causes : "
        echo "No instance of redis is installed in this system"
        exit
    fi
    redis-cli FLUSHALL
    go build .
    if [ $? -ne "0" ];then
        echo "oh oh, something went wrong, possible causes : "
        echo "No go toochains installed"
        exit
    fi
    chmod +x GRDNS
    sudo ./GRDNS
fi



