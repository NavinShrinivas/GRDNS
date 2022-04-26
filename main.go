package main

import (
	"fmt"
    "sync"
    "net"
);

var addr = net.UDPAddr{
    Port: 53,
    IP:   net.ParseIP("0.0.0.0"),
}

var wg sync.WaitGroup;

func main(){
    Conn,err:= net.ListenUDP("udp",&addr)
    if err!=nil{
        fmt.Println("Error listening to port : ",err);
    }else{
        fmt.Println("Listening to port 53");
    }
    wg.Add(1)
    go serverstart(Conn)
    wg.Wait()
}
