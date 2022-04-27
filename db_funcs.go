package main

import (
    "fmt"
    "github.com/miekg/dns"
    "github.com/gomodule/redigo/redis"
)
func newPool() *redis.Pool {
    return &redis.Pool{
        // Maximum number of idle connections in the pool.
        MaxIdle: 800,
        // max number of connections
        MaxActive: 12000,
        // Dial is an application supplied function for creating and
        // configuring a connection.
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", ":6379")
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    }
}
func FlushToDB(Record ResponseStruct) bool {

    var pool = newPool()
    var c = pool.Get()

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

func ReturnWithAnswers(domain string,res *dns.Msg){

    fmt.Println("Resolving from Cache!",domain)
    var pool = newPool()
    var c = pool.Get()
    defer c.Close()

    mod_dom := domain+"."
    for _,record_number := range domain_map[mod_dom]{
        raw_str,err :=  redis.String(c.Do("HGET",record_number,"reply"))
        checkError(err)
        a,err := dns.NewRR(raw_str)
        checkError(err)
        fmt.Println(a)
        res.Answer = append(res.Answer,a)
        /*
         *if get_fields_whitespace(raw_str).Typ == "CNAME"{
         *    ReturnWithAnswers(get_fields_whitespace(raw_str).Reply[:len(get_fields_whitespace(raw_str).Reply)-1],res)
         *    break
         *}
         */
    }
    return
}
