echo "Needed Dependencies : redis "
echo "Which distro are u running me on?"
echo "1.Debian base distros [Ubuntu, pop, ...]"
echo "2.Arch base distro [EOS, Manjaro, ..]"
echo "3.MacOs"
echo -n "Enter option :"
read option


case $option in
    1)
    # INSTALL REDIS FOR UBUNTU
        curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

        echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

        sudo apt-get update
        sudo apt-get install redis
        ;;

    2)
    # INSTALL FOR ARCH
        sudo snap install redis
        ;;

    3)
    # INSTALL FOR MAC
        brew install redis
        ;;
esac



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



