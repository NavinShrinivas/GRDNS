package main

import (
    "fmt"
    "github.com/miekg/dns"
    "github.com/gomodule/redigo/redis"
)

func FlushToDB(Record ResponseStruct) bool {
    _,err := c.Do("HSET", record_number,"name", Record.Name, "ttl",Record.Rawttl, "class", Record.Rawclass, "type", Record.Rawrrtype, "reply", Record.Rawstr,"length",Record.Rawrdlength)
    if err != nil {
        checkError(err)
        return false
    }
    domain_map[Record.Name] = append(domain_map[Record.Name],record_number)
    record_number++;
    fmt.Println("Record inserted to Database!")
    return true
}


func EntryExists(domain string)bool{
    mod_dom := domain+"."
    if len(domain_map[mod_dom]) == 0{
        return false 
    }else{
        return true
    }
}


func ReturnWithAnswers(domain string)*dns.Msg{
    mod_dom := domain+"."
    fmt.Println("Resolving from Cache!")
    var res = new(dns.Msg)
    res.MsgHdr.Response = true;
    res.MsgHdr.Rcode = 0; //No error handling :(
    for _,record_number := range domain_map[mod_dom]{
        raw_str,err :=  redis.String(c.Do("HGET",record_number,"reply"))
        checkError(err)
        a,err := dns.NewRR(raw_str)
        checkError(err)
        fmt.Println(a)
        res.Answer = append(res.Answer,a)
    }
    return res
}
