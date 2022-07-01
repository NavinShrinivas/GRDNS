#perf/stability test systems :
echo "this will test any dns server running on localhost. As of now, all test mode can do this"
FILE=10m_query10
if [ "$1" = "--10mtest" ];then
    if test -f "$FILE";then
        dnsperf -d 10m_query10 -s 0.0.0.0
        if [ $? -ne "0" ];then
            echo "please install dnsperf to run this script"
            exit
        fi
    else
        unxz -k 10m_query10.xz
        dnsperf -d 10m_query10 -s 0.0.0.0
        if [ $? -ne "0" ];then
            echo "please install dnsperf to run this script"
            exit
        fi
    fi
fi


