# GRDNS : The easy to setup DNS caching server 

GRDNS aims to be the program to run on yout DNS server when one needs fast resolve times. Consider major ISP who have their own DNS servers (Which may sometimes be used for tracking, but thats not the point of this project) who want to server throusands of users simultaneously. GRDNS does this by using in memeory key value pairs and extreme concurrency.

## Basic Info

- Built using GoLang and redis.
- Clean and effective code 
- Hash based database storing records (redis)
- Multi threading done right
- Using low overhead libraries wherever possible 
- Using in memory data as much as possible when also maintaining stability
- using newer languages such as GO lang

### Building locally : 

> Note : this project is built for the linux system, compatibility with windows is not confirnmed(This can be an issue, propose if you wish to).
To run this porject, you need to have golang and redis installed. After getting those two 
```bash
$ ./run.sh
```
Should compile and start the DNS server.

### Testing locally :

- A single resolve test can be done like so : 
```bash
$ dig @0.0.0.0 google.com
```
It's best you run the above command twice. This is cus the first resolve doesn't hit the cache, second will.
- A performance test can be done like so [needs dnsperf installed before hand]:
```
./test.sh --10mtest
```

### Maintainers in charge : 

[Navin Shrinivas](https://github.com/NavinShrinivas)
[Mukund Deepak](https://github.com/mukunddeepak)

## Performance : 

- Down below we see comparision to cloudflare and google dns server (which note runs on very powerful hardware, mine runs simply on my computer) : 
- GRDNS : ~110QP/s || google's DNS : ~360QP/s || cloudflare's DNS : 220QP/s
- images : 

![GRDNS](https://user-images.githubusercontent.com/42774281/178130007-9eac0476-3eb6-407f-afcb-fb4a12be8c4a.png)

![google dns](https://user-images.githubusercontent.com/42774281/178130001-3d135c32-a9e3-4cf5-8602-844a13040db6.png)

![cloudflare dns](https://user-images.githubusercontent.com/42774281/178130010-9983f49a-2bb7-4d8d-a231-840545824b5f.png)



