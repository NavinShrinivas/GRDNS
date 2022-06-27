//This file has all the struct definitions and globals

package Modules 

import(
    "net"
    "sync"
    "fmt"
)


var Record_number int64; //Very very bad idea, will definetly break down the line! 
//But I cant see any other way of tracking mutiple records for the same domain.
//On the other hand, we dont even have that many ip addresses (atleast not in the ip4 space)

var domain_map = make(map[string][]int64) //Will act as a map to the respective records for each domain

var Addr = net.UDPAddr{
    Port: 53,
    IP:   net.ParseIP("0.0.0.0"),
}

var Work_Group sync.WaitGroup;

type ResponseStruct struct{
    Name string  //For ease of use
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

type System struct{
    FreeThreads int64 `json:",string"`
}

var System_State System;

type Job struct{
    buffer []byte 
    Conn *net.UDPConn 
    Caddr *net.UDPAddr
}

var Thread_channels [1024]chan Job //Cant have more than 1024 thread

func CheckError(err error){
    if err!=nil{
        fmt.Println("Something went wrong : ",err)
    }
}



