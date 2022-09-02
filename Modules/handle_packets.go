package Modules

import (
    "fmt"
    "net"
    "github.com/miekg/dns"
);

func Serverstart(Conn *net.UDPConn){
    //Load balanced thread pool allocater

    //The way this fucntion is designed saves time and resources spawing thread everytime!.
    //Although it does lead to some dropped packets but thats and extreme

    for i:=0 ;i<int( System_State.FreeThreads/16 );i++{
        go UDPConnHandlers(Conn);
    }
    for i:=0 ;i<int( System_State.FreeThreads );i++{
        go request_handle_thread(Thread_channels[i]) //spawnning threads
        System_State.FreeThreads = System_State.FreeThreads - 1;
    }
    go LoadBalancer(LoadBalancerChannel); 
}


func UDPConnHandlers(Conn *net.UDPConn){
    for{
        buffer := make([]byte,10000)
        _,CAddr,err := Conn.ReadFromUDP(buffer)
        CheckError(err)
        new_job := Job{ //Creating new Job
            buffer : buffer,
            Conn : Conn,
            Caddr : CAddr,
        }
        LoadBalancerChannel <- new_job;
    }

}


func LoadBalancer(LoadBalancerChannel chan Job){
    var min_buffer_mod = System_State.FreeThreads;
    thread_counter := 0;
    for{
        new_job := <- LoadBalancerChannel
        Thread_channels[thread_counter%int(min_buffer_mod)]<-new_job;
        thread_counter++;

    }
}

func request_handle_thread(job chan Job){

    for{ //The thread continuously checks for any job
        new_job := <- job
        packet_buff := new_job.buffer;
        UDPConn := new_job.Conn;
        UDPCaddr := new_job.Caddr;
        DNSpacketObj := new(dns.Msg);
        DNSpacketObj.Unpack(packet_buff);

        for _,it:=range DNSpacketObj.Question{
            DNSpacketObj.MsgHdr.Response = true;
            DNSpacketObj.MsgHdr.Rcode = 0;//meaning successfull
            if EntryExists(string(it.Name)){
                ReturnWithAnswers(string(it.Name),DNSpacketObj)
            }else{
                resolve(DNSpacketObj,string(it.Name))
            }

            if DNSpacketObj!=nil{         


                resbuf,_ := DNSpacketObj.Pack()

                //Writing back to client
                _, err := UDPConn.WriteToUDP(resbuf, UDPCaddr)
                CheckError(err)
            }
        }
    }
}

func error_ret_fail(DNSpacketObj *dns.Msg){
    DNSpacketObj.MsgHdr.Rcode = 2 //Meaning failed on server side...only error check as for now
    return
}

func resolve(DNSpacketObj *dns.Msg,Name string){ //inserts records to map and databse
    //First check in redis server 

    //if not found resolve using 8.8.8.8/1.1.1.1 (for now)
    //Root server : 198.41.0.4
    msg := new(dns.Msg)
    msg.SetQuestion(dns.Fqdn(Name),dns.TypeA) //FQDN : fully qualified Domain name
    in, err := dns.Exchange(msg, "8.8.8.8:53")
    if err!=nil{
        error_ret_fail(DNSpacketObj)
        return
    }
    if in!=nil{
        //Without this we get nil dereference errors
        response_handlers(DNSpacketObj,in)
    }
    return
}


func response_handlers(DNSpacketObj *dns.Msg ,res *dns.Msg){

    //Fetching answer records [MOST IMP] : 
    for i:=0;i<len(res.Answer);i++{
        it := res.Answer[i]
        DNSpacketObj.Answer = append(DNSpacketObj.Answer,res.Answer[i])
        res_struct := get_fields_whitespace(it.String())

        if res_struct.Typ == "CNAME"{
            resolve(DNSpacketObj, res_struct.Reply)
        }

        //Preparing to flus to db
        res_struct.Rawname = it.Header().Name
        res_struct.Rawclass = it.Header().Class
        res_struct.Rawrdlength = it.Header().Rdlength
        res_struct.Rawstr = it.String()
        res_struct.Rawrrtype = it.Header().Rrtype
        res_struct.Rawttl = it.Header().Ttl 
        res := FlushToDB(res_struct)
        if res==false{
            fmt.Println("Something wrong with redis server!")
        }
    }
    for i:=0;i<len(res.Ns);i++{
        //Cant handle caching Auth section as of now 
        DNSpacketObj.Ns = append(DNSpacketObj.Ns,res.Ns[i])

    }
    return 
}
