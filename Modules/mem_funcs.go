package Modules

import (
    "fmt"
    "github.com/miekg/dns"
    "github.com/gomodule/redigo/redis"
)
func newPool() *redis.Pool {
    return &redis.Pool{
        // Maximum number of idle connections in the pool.
        MaxIdle: 50,
        // max number of connections
        MaxActive: 120000,
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

    temp_arr := FetchMapFunction(Record.Name);
    temp_arr = append(temp_arr,Record_number)
    UpdateMapFunction(Record.Name,temp_arr);

    _,err := c.Do("HSET", Record_number,"name", Record.Name, "ttl",Record.Rawttl, "class", Record.Rawclass, "type", Record.Rawrrtype, "reply", Record.Rawstr,"length",Record.Rawrdlength) //Record number acts as the key in db
    if err != nil {
        CheckError(err)
        return false
    }
    fmt.Println("Record inserted to Database!")
    Record_number++;
    return true
}


func EntryExists(domain string)bool{
    mod_dom := domain
    if len(FetchMapFunction(mod_dom)) == 0{ //Checking map if there is any record numbers the memory map where the name is the key
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
    mod_dom := []string{domain}
    for j:=0;j<len(mod_dom);j++{
        map_value_len := 1;
        for i:=0;i<map_value_len;i++{ //Should not be converted to iterrator based loops
            map_value := FetchMapFunction(mod_dom[j]);
            record_number := map_value[i]
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
                //Something is going wrong, falling back to know server 
                if len(FetchMapFunction(cname_response_domain)) == 0{
                    resolve(res,cname_response_domain)
                    return
                }
                mod_dom = append(mod_dom,cname_response_domain)
            }
            res.Answer = append(res.Answer,a)
            map_value_len = len(map_value)
        }
    }
    return
}
//append(domain_map[new_job.Name],new_job.Record_number) 
func UpdateMapFunction(Name string,value []int64){
    //go lang regular maps in concurrency : https://stackoverflow.com/questions/36167200/how-safe-are-golang-maps-for-concurrent-read-write-operations
    //Needs a rwmutex lock

    //Any read lock tries to aqquire that appear behind the write lock is blocked, this is to ensure that writing also eventually gets a chance at getting the lock.
    if len(FetchMapFunction(Name)) == 0{
        RecordMapRWLock.Lock();
        defer RecordMapRWLock.Unlock();
        domain_map[Name] = value     
    }
}


func FetchMapFunction(Name string) []int64{
    RecordMapRWLock.RLock();
    defer RecordMapRWLock.RUnlock();
    res := domain_map[Name];
    return res
}
