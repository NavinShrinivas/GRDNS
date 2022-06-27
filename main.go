package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
    "GRDNS/Modules"
);


func main(){

    //Get number of threads for thread pool, have to parse JSON 
    dat, err := os.ReadFile("./system.conf")
    Modules.CheckError(err)
    json.Unmarshal(dat,&Modules.System_State)
    fmt.Print("Recieved from conf file :\n",string(dat),"\n",Modules.System_State)

    if Modules.System_State.FreeThreads > 1024{
        fmt.Println("Thats one monster of a Modules.System_State you got there, sadly our software has constrains..mind getting your engineering team on this?")
        return
    }

    if Modules.System_State.FreeThreads <= 2 {
        fmt.Println("Atleast 3 threads are needed to be spawned to get this program working.")
        return
    }

    Modules.Record_number = 0;
    Conn,err:= net.ListenUDP("udp",&Modules.Addr)
    if err!=nil{
        fmt.Println("Error listening to port : ",err);
    }else{
        fmt.Println("Listening to port 53");
    }

    //Making those channels buffered 
    for i:=0;i<int(Modules.System_State.FreeThreads);i++{
        Modules.Thread_channels[i] = make(chan Modules.Job,10000); //No more than 10000 jobs can be buffered at a time
    }

    Modules.Work_Group.Add(1)
    Modules.System_State.FreeThreads = Modules.System_State.FreeThreads - 1;
    Modules.Serverstart(Conn) //One of the threads is given to Load balanced thread pool allocater
    Modules.Work_Group.Wait()
    fmt.Println()
}
