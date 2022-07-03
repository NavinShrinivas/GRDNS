package Modules

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
    defer c.Close()
    defer pool.Close()

    _,err := c.Do("HSET", Record_number,"name", Record.Name, "ttl",Record.Rawttl, "class", Record.Rawclass, "type", Record.Rawrrtype, "reply", Record.Rawstr,"length",Record.Rawrdlength) //Record number acts as the key in db
    if err != nil {
        CheckError(err)
        return false
    }
    domain_map[Record.Name] = append(domain_map[Record.Name],Record_number) //Inserting rescord number to memory map, key is the domain question
    Record_number++;
    fmt.Println("Record inserted to Database!")
    return true
}


func EntryExists(domain string)bool{
    mod_dom := domain+"."
    if len(domain_map[mod_dom]) == 0{ //Checking map if there is any record numbers the memory map where the name is the key
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
    defer pool.Close()

    mod_dom := []string{domain+"."}
    for j:=0;j<len(mod_dom);j++{
        for i:=0;i<len(domain_map[mod_dom[j]]);i++{ //Should not be converted to iterrator based loops
            record_number := domain_map[mod_dom[j]][i]
            raw_str,err :=  redis.String(c.Do("HGET",record_number,"reply"))
            CheckError(err)
            a,err := dns.NewRR(raw_str)
            CheckError(err)
            if a.Header().Rrtype == 5{
                //Meaning its reply is CNAME record
                //Need to fetch details for CNAME record now
                //Need to refactor later : 
                response_struct := get_fields_whitespace(raw_str);
                cname_response_domain := response_struct.Reply
                mod_dom = append(mod_dom,cname_response_domain)
            }
            res.Answer = append(res.Answer,a)
        }
    }
    return
}
