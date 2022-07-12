package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
    "GRDNS/Modules"
);


func main(){

    //Get number of threads for thread pool, have to parse JSON, to be deprecated soon
    dat, err := os.ReadFile("./system.conf")
    Modules.CheckError(err)
    json.Unmarshal(dat,&Modules.System_State)
    fmt.Print("Recieved from conf file :\n",string(dat),"\n",Modules.System_State)

    if Modules.System_State.FreeThreads > 1024{
        fmt.Println("Thats one monster of a Modules.System_State you got there, sadly our software has constrains..mind getting your engineering team on this?")
        return
    }

    if Modules.System_State.FreeThreads <= 15 {
        fmt.Println("Atleast 16 soft threads are needed to be spawned to get this program working.")
        return
    }

    Modules.Record_number = 0//Starting record numbers at 0,meaning redis server needs to be flushed during every start.
    Conn,err:= net.ListenUDP("udp",&Modules.Addr)//Main UDP connection socket
    if err!=nil{
        fmt.Println("Error listening to port : ",err);
    }else{
        fmt.Println("Listening to port 53");
    }

    //Initializing
    //These are thre buffered channels that the thread threads that are part of the thread pool 
    //Keep checking for any jobs

    for i:=0;i<int(Modules.System_State.FreeThreads);i++{
        Modules.Thread_channels[i] = make(chan Modules.Job,7000); //No more than 7000 jobs can be buffered at a time
    } 
    Modules.UpdateMapBuffer = make(chan Modules.InsertRecordJob,1000000) //No more than 1000000 record inserts buffered at one time
    Modules.LoadBalancerChannel = make(chan Modules.Job, 1000000);

    Modules.Work_Group.Add(1)//making the program wait for other threads, no logic here. hard coded.
    Modules.Serverstart(Conn) //First we spawnt the load balancer thread.
    Modules.Work_Group.Wait() //As long as Work_Groups value is greater than 0, ii will wait for some child thread to make it to 0.
    fmt.Println()
}
