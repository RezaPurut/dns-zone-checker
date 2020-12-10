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

	var (
		dns_file string
		zone_dir string
		bast_addr string
		bast_user string
		bast_pass string
		bast_port string
		bast_key string
		target_user string
		target_pass string
		target_port string
		target_key string
		single_file string
		
		bulk_check bool
	)
  
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
	bastionConn, err := sshConnect(bast_addr, bast_user, bast_pass, bast_key, bast_port)
	
	if err != nil {
		fmt.Println(err)
	} else {
		if bulk_check {
			target_list := difference(getFileName(zone_dir), readFnameInConfig(zone_dir, dns_file))

			fmt.Println(target_list)
		}
	}
  
}

// function to check error	
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PublicKey(file string) ssh.AuthMethod {
	buff, err := ioutil.ReadFile(file)

	check(err)

	key, err := ssh.ParsePrivateKey(buff)

	check(err)

	return ssh.PublicKeys(key)
}

func sshConfig(username, pass, key_path string) *ssh.ClientConfig{
	var config *ssh.ClientConfig
	if key_path != "" {
		config = &ssh.ClientConfig {
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
				PublicKey(key_path),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout: 5 * time.Second,
		}
	} else {
		config = &ssh.ClientConfig {
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout: 5 * time.Second,
		}
	}
	
	return config
}

func sshConnect(server, username, pass, key_path, port string) (*ssh.Client, error) {
	config := sshConfig(username, pass, key_path)

	connection, err := ssh.Dial("tcp", server + ":" + port, config)
	
	return connection, err
}

func sshSession(conn *ssh.Client) {
	session, err := conn.NewSession()

	check(err)

	defer session.Close()

	sessStdOut, err := session.StdoutPipe()
	check(err)
	go io.Copy(os.Stdout, sessStdOut)

	sessStderr, err := session.StderrPipe()
	check(err)
	go io.Copy(os.Stderr, sessStderr)

	err = session.Run("ls -la")

}

// This function read the .conf file record line by line 
// using regex to search for zones directory keyword until 
// it finds a double quote ("). Then the search result is 
// cleaned by removing the duplicated file name.
func readFnameInConfig(zone_dir, dns_file string) []string {
	f, err := os.Open(dns_file)
	check(err)

	defer f.Close()

	keyword := strings.Split(zone_dir, "/")

	scanner := bufio.NewScanner(f)

	var m = make(map[string]bool)
	var a = []string{}

	// if user set zone_dir end with '/'
	if keyword[len(keyword)-1] == "" {
		keyword[len(keyword)-1] = keyword[len(keyword)-2]
	}

	r := regexp.MustCompile(keyword[len(keyword)-1] + "/(.*?)\"")
	for scanner.Scan() {
		fname := r.FindAllStringSubmatch(scanner.Text(), -1)
		for _, v := range fname {
			s := string(v[1])
			if m[s] != true {
				a = append(a, s)
				m[s] = true
			}
		}
	}
	
	return a
}
