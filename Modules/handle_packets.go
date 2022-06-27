package Modules

import (
    "fmt"
    "net"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/miekg/dns"
);

func Serverstart(Conn *net.UDPConn){
    //Load balanced thread pool allocater

    //The way this fucntion is designed saves time and resources spawing thread everytime!.
    //Although it does lead to some dropped packets but thats and extreme

    for i:=0 ;i<int( System_State.FreeThreads-1 );i++{ //-1 cus one thread taken up by this function
        go request_handle_thread(Thread_channels[i]) //spwanning threads
        System_State.FreeThreads = System_State.FreeThreads - 1;
    }
    fmt.Println("Threads spwaned!")
    for{
        buffer := make([]byte,10000)
        _,CAddr,err := Conn.ReadFromUDP(buffer)
        CheckError(err)
        fmt.Println("Connection from : ",CAddr)
        new_job := Job{
            buffer : buffer,
            Conn : Conn,
            Caddr : CAddr,
        }
        var min_buffer = 0;
        var min_buffer_len = len(Thread_channels[0]);
        for i:=0;i<int(System_State.FreeThreads-1);i++{
            if(len(Thread_channels[i]) < min_buffer_len){
                min_buffer_len = len(Thread_channels[i]);
                min_buffer = i;
            }
        }
        fmt.Println("Job given to thread",min_buffer)
        Thread_channels[min_buffer]<-new_job;
    }
}

func request_handle_thread(job chan Job){

    for{
        new_job := <- job
        packet_buff := new_job.buffer;
        UDPConn := new_job.Conn;
        UDPCaddr := new_job.Caddr;
        //Fetching DNS layers and parsing to object, Can be converted to handle with dns package later 
        packetlayers := gopacket.NewPacket(packet_buff,layers.LayerTypeDNS,gopacket.Default) 
        DNSlayer := packetlayers.Layer(layers.LayerTypeDNS)
        DNSpacketObj := DNSlayer.(*layers.DNS)

        fmt.Println("Questions Recieved : ")
        for i,it:=range DNSpacketObj.Questions{
            fmt.Println("\t Question",i+1,":",string(it.Name))

            req_id := DNSpacketObj.ID; //Used by All DNS System_States to ensure authenticity
            var response = new(dns.Msg);
            if EntryExists(string(it.Name)){
                response.MsgHdr.Response = true;
                response.MsgHdr.Rcode = 0; //No error handling :(
                response.MsgHdr.RecursionDesired = true;
                l := new(dns.Msg)
                l.Unpack(packet_buff)
                response.Question = l.Question;
                ReturnWithAnswers(string(it.Name),response)

            }else{
                response = resolve(string(it.Name))
            }


            if response!=nil{         
                response.MsgHdr.Id = req_id;
                resbuf,_ := response.Pack()

                //Writing back to client
                _, err := UDPConn.WriteToUDP(resbuf, UDPCaddr)
                CheckError(err)
            }
        }
    }
}



func resolve(Name string) *dns.Msg{
    //First check in redis server 

    //if not found resolve using 8.8.8.8/1.1.1.1 (for now)
    //Root server : 198.41.0.4
    msg := new(dns.Msg)
    msg.SetQuestion(dns.Fqdn(Name),dns.TypeA) //FQDN : fully qualified Domain name
    in, err := dns.Exchange(msg, "1.1.1.1:53")
    CheckError(err)
    if in!=nil{
        //Without this we get nil dereference errors
        response_handlers(in)
        return in
    }
    return nil
}


func response_handlers(res *dns.Msg){

    //Fetching answer records [MOST IMP] : 
    for i:=0;i<len(res.Answer);i++{
        it := res.Answer[i]
        res_struct := get_fields_whitespace(it.String())
        fmt.Print(res_struct.Name," ")
        fmt.Print(res_struct.Typ," ")
        fmt.Print(res_struct.Class," ")
        fmt.Print(res_struct.Ttl," ")
        fmt.Println(res_struct.Reply)
        if res_struct.Typ == "CNAME"{
            new_record := resolve(res_struct.Reply)
            res.Answer = append(res.Answer,new_record.Answer...)

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
}
