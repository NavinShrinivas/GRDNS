package main 

import (
    "fmt"
    "net"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
);

func serverstart(Conn *net.UDPConn){
    for{
        buffer := make([]byte,10000)
        no_bits,CAddr,_ := Conn.ReadFromUDP(buffer)
        fmt.Println("Connection from : ",CAddr)
        fmt.Println("Size of Recieved packet : ",no_bits)
        go handle_request(buffer);
    }

}


func handle_request(buffer []byte){
        packetlayers := gopacket.NewPacket(buffer,layers.LayerTypeDNS,gopacket.Default) 
        //Above gives a set of layer of the packet revieved
        //Where the DNS layer is filled with our Recieved bits
        DNSlayer := packetlayers.Layer(layers.LayerTypeDNS)
        //Above only extracts the DNS layer from set of layers 
        //with above layer we can create an object :) 
        DNSpacketObj := DNSlayer.(*layers.DNS);
        fmt.Println(string( DNSpacketObj.Questions[0].Name ));

}
