package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
);


func main(){

    //Get number of threads for thread pool, have to parse JSON 
    dat, err := os.ReadFile("./system.conf")
    checkError(err)
    json.Unmarshal(dat,&system)
    fmt.Print("Recieved from conf file :\n",string(dat),"\n",system)

    if system.FreeThreads > 1024{
        fmt.Println("Thats one monster of a system you got there, sadly our software has constrains..mind getting your engineering team on this?")
        return
    }

    if system.FreeThreads <= 2 {
        fmt.Println("Atleast 3 threads are needed to be spawned to get this program working.")
        return
    }

    record_number = 0;
    Conn,err:= net.ListenUDP("udp",&addr)
    if err!=nil{
        fmt.Println("Error listening to port : ",err);
    }else{
        fmt.Println("Listening to port 53");
    }

    //Making those channels buffered 
    for i:=0;i<int(system.FreeThreads);i++{
        thread_channels[i] = make(chan Job,10000); //No more than 10000 jobs can be buffered at a time
    }

    wg.Add(1)
    system.FreeThreads = system.FreeThreads - 1;
    serverstart(Conn) //One of the threads is given to Load balanced thread pool allocater
    wg.Wait()
    fmt.Println()
}
