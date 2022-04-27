package main

import (
    "fmt"
    "sync"
    "net"
);


//----------Gloabls----------

var record_number int64; //Very very bad idea, will definetly break down the line! 
//But I cant see any other way of tracking mutiple records for the same domain.

var domain_map = make(map[string][]int64) //Will act as a map to the respective records for each domain


var addr = net.UDPAddr{
    Port: 53,
    IP:   net.ParseIP("0.0.0.0"),
}

var wg sync.WaitGroup;

type ResponseStruct struct{
    Name string  //For printing purposes only 
    Typ string //"
    Class string //"
    Reply string //"
    Ttl string //"
    Rawstr string //For creating packets in the end 
    Rawname string  //"
    Rawrrtype uint16  //"
    Rawclass uint16 //"
    Rawttl uint32 //"
    Rawrdlength uint16 //"
}


func main(){
    record_number = 0;
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
