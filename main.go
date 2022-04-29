package main

import (
    "fmt"
    "net"
    "os"
    "encoding/json"
);


func main(){

    //Get number of threads for thread pool, have to parse JSON 
    dat, err := os.ReadFile("./system.conf")
    checkError(err)
    json.Unmarshal(dat,&system)
    fmt.Print("Recieved from conf file :\n",string(dat),"\n",system)

    if system.FreeThreads > 1024{
        fmt.Println("Thats one monster of a system you got there, sadly our software has constrains..mind getting your engineering team on this?")
    }

    record_number = 0;
    Conn,err:= net.ListenUDP("udp",&addr)
    if err!=nil{
        fmt.Println("Error listening to port : ",err);
    }else{
        fmt.Println("Listening to port 53");
    }

    wg.Add(1)
    system.FreeThreads = system.FreeThreads - 1;
    serverstart(Conn) //One of the threads is given to Load balanced thread pool allocater
    wg.Wait()
}
