package main

import (
	"fmt"
	"flag"
	"os"
	"regexp"
	"bufio"
	"io/ioutil"
	"io"
	"strings"
	"time"
	
	"golang.org/x/crypto/ssh"	
)

func main() {

	var dns_file string
	var zone_dir string
	var bast_addr string
	var bast_user string
	var bast_pass string
	var bast_port string
	var bast_key string
	var target_user string
	var target_pass string
	var target_port string
	var target_key string
	var single_file string

	var bulk_check bool
  
	flag.StringVar(&dns_file, "dns-file", "", "DNS configuration file")
	flag.StringVar(&zone_dir, "zone-dir", "", "Zone file directory")

	flag.StringVar(&bast_addr, "bastion-addr", "", "Server address or name for Bastion Host")
	flag.StringVar(&bast_user, "bastion-user", "", "Username for Bastion Host")
	flag.StringVar(&bast_pass, "bastion-pass", "", "Password for Bastion Host")
	flag.StringVar(&bast_port, "bastion-port", "", "Port for Bastion Host")
	flag.StringVar(&bast_key, "bastion-key", "", "Private Key for Bastion Host")

	flag.StringVar(&target_user, "target-user", "", "Username for Target Host")
	flag.StringVar(&target_pass, "target-pass", "", "Password for Target Host")
	flag.StringVar(&target_port, "target-port", "", "Port for Target Host")
	flag.StringVar(&target_key, "target-key", "", "Private Key for Target Host")
	
	flag.BoolVar(&bulk_check, "bulk", false, "Enable Bulk Checking")
	flag.StringVar(&single_file, "single-zone", "", "Zone file to check (only use this to check individual file)")

	flag.Parse()

	pass_list := strings.Split(target_pass, ",")
	port_list := strings.Split(target_port, ",")
  
	fmt.Println("Login to Bastion Host...")
  
}
