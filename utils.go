package main  
  
import (  
 "fmt"  
 "os"
 "net"
 fqdn "github.com/Showmax/go-fqdn"
)  

func checkStringForIpOrHostname(host string) {  
 addr := net.ParseIP(host)  
 if addr == nil {  
  fmt.Println("Given String is a Domain Name")  
  
 } else {  
  fmt.Println("Given String is a Ip Address")  
 }  
}  

func main() {  
 hostname, error := os.Hostname()  
 if error != nil {  
  fmt.Println("error : ", error)
  panic(error)
 }  
  fmt.Println("hostname:", hostname)
  fmt.Println(fqdn.Get())  
 checkStringForIpOrHostname("google.com")
 checkStringForIpOrHostname("2001:0db8:85a3:0000:0000:8a2e:0370:7334") // 192.168.1.0

}  