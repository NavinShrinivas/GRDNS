#perf/stability test systems :
echo "this will test any dns server running on localhost. As of now, all test mode can do this"
FILE=5m_query
if [ "$1" = "--5mtest" ];then
    if test -f "$FILE";then
        dnsperf -d 5m_query -s 0.0.0.0
        if [ $? -ne "0" ];then
            echo "please install dnsperf to run this script"
            exit
        fi
    else
        unxz -k 5m_query.xz
        dnsperf -d 5m_query -s 0.0.0.0
        if [ $? -ne "0" ];then
            echo "please install dnsperf to run this script"
            exit
        fi
    fi
fi

