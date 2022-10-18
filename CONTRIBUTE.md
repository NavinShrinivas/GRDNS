# Contributors Help 

## File stsructure
This file hold the environment variables handling and passes off a UDP connection to Modules/handle_packets.go. 
This file spawns a bunch of UDPConn handlers and a bunch of LoadBalancers. Communication between load balancers and request_handle_thread is a whole lot of channels (each thread has its own channel).

request_handle_thread asks the functions present in mem_func whether this records exists in cache or not. It takes actions accordingly.

mem_func is where all redis and in memory hashmap operations are done.

In this entire project, wherever CNAME reocord is handled, things get messy.

## PR Guidelines 

Make sure to tag the issue number that you are attempting to solve with the PR. Give credits and resources when its due.

## Issues guidelines 

Ones that are trying to propose new issue, make sure you explain the problem clearly. Along with any initial solutions. Also attach any screenshots that might explain the error.
