package main

import (
    "fmt"
    "net"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/miekg/dns"
);

func serverstart(Conn *net.UDPConn){
    for{
        buffer := make([]byte,10000)
        no_bits,CAddr,err := Conn.ReadFromUDP(buffer)
        checkError(err)
        fmt.Println("Connection from : ",CAddr)
        fmt.Println("Size of Recieved packet : ",no_bits)
        go handle_request(buffer,CAddr,Conn);
    }

}

func checkError(err error){
    if err!=nil{
        fmt.Println("Something went wrong : ",err)
    }
}

func handle_request(buffer []byte,Caddr *net.UDPAddr,Conn *net.UDPConn){
    packetlayers := gopacket.NewPacket(buffer,layers.LayerTypeDNS,gopacket.Default) 
    //Above gives a set of layer of the packet revieved
    //Where the DNS layer is filled with our Recieved bits
    DNSlayer := packetlayers.Layer(layers.LayerTypeDNS)
    //Above only extracts the DNS layer from set of layers 
    //with above layer we can create an object :) 
    DNSpacketObj := DNSlayer.(*layers.DNS)
    fmt.Println("Questions Recieved : ")
    for i,it:=range DNSpacketObj.Questions{
        fmt.Println("Question",i+1,":",string(it.Name))

        req_id := DNSpacketObj.ID; //Used by All DNS systems to ensure authenticity
        var response = new(dns.Msg);
        if EntryExists(string(it.Name)){
            response.MsgHdr.Response = true;
            response.MsgHdr.Rcode = 0; //No error handling :(
            response.MsgHdr.RecursionDesired = true;
            l := new(dns.Msg)
            l.Unpack(buffer)
            response.Question = l.Question;
            ReturnWithAnswers(string(it.Name),response)

        }else{
           response = resolve(string(it.Name))
        }
        

        if response!=nil{         
            response.MsgHdr.Id = req_id;
            resbuf,_ := response.Pack()

            //Writing back to client
            _, err := Conn.WriteToUDP(resbuf, Caddr)
            checkError(err)
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
    checkError(err)
    if in!=nil{
        //Without this we get nil dereference errors
        response_handlers(in)
        return in
        fmt.Println()
    }
    return nil
}


func response_handlers(res *dns.Msg){
    //gets the possible fields and also pushed to database
    //Note sure if all DNS server do this, but we are caching A, CNAME, NS 
    //But will be returning only CNAME and A type 
    //Why so you ask? Well we have here a recursive DNS looker, meaning if it doesnt have 
    //The needed domain in cache it is gonna go looking for it, meaning we will be reaching 
    //Thso authoritative nameserver time and again.
    //Fetching authoritative records :
    fmt.Println("New Auth records : ")
    for count,it := range res.Ns{
        res_struct := get_fields_whitespace(it.String())
        fmt.Println("Record : ",count+1)
        fmt.Print(res_struct.Name," ")
        fmt.Print(res_struct.Typ," ")
        fmt.Print(res_struct.Class," ")
        fmt.Print(res_struct.Ttl," ")
        fmt.Println(res_struct.Reply)
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
    //Fetching answer records [MOST IMP] : 
    fmt.Println("New Answer records : " )
    for count,it := range res.Answer{
        res_struct := get_fields_whitespace(it.String())
        fmt.Println("Record : ",count+1)
        fmt.Print(res_struct.Name," ")
        fmt.Print(res_struct.Typ," ")
        fmt.Print(res_struct.Class," ")
        fmt.Print(res_struct.Ttl," ")
        if res_struct.Typ == "CNAME"{
            resolve(res_struct.Reply)
        }
        fmt.Println(res_struct.Reply)
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
