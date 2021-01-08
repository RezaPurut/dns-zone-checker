package main
// TODO: comment function
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
	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		zone_file string
		// dns_file string
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
		// single_file string
		log_file string
		// bulk_check bool
	)

	flag.StringVar(&zone_file, "zone-file", "", "Zone file/files")
	// flag.StringVar(&dns_file, "dns-file", "", "DNS configuration file")
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
	
	// flag.BoolVar(&bulk_check, "bulk", false, "Enable Bulk Checking")
	// flag.StringVar(&single_file, "single-zone", "", "Zone file to check (only use this " + 
	// "to check individual file)")
	flag.StringVar(&log_file, "log-file", "/var/log/dns-check.txt", "Log file")
	flag.Parse()

	file_list := strings.Split(zone_file, ",")
	pass_list := strings.Split(target_pass, ",")
	port_list := strings.Split(target_port, ",")

	initializeLogging(log_file)

	fmt.Println("Login to Bastion Host...")
	log.Info("Login to Bastion Host")

	bastionConn, err := sshConnect(bast_addr, bast_user, bast_pass, 
		bast_key, bast_port)
	
	if err != nil {
		fmt.Println("Please check your bastion parameter", err)
		log.Error("Could not login to Bastion: " + err.Error())
	} else {

		if zone_dir != "" {

			for file := range(getFileName(zone_dir)) {
				target_addr := readFile(zone_dir, file)

				fmt.Printf("Checking file %s\n", file)
				fmt.Printf("SSHing %s\n", target_addr)
	
				attemptConnect(bastionConn, port_list, pass_list, target_user, 
					target_key, target_addr)
			}
			
		} else {
			for file := range(file_list) {
				target_addr := readFile(zone_dir, file_list[file])
	
				fmt.Printf("SSHing %s\n", target_addr)
	
				attemptConnect(bastionConn, port_list, pass_list, target_user, 
					target_key, target_addr)
			}
		}
		
		/* if bulk_check {
			target_list := difference(getFileName(zone_dir), 
			readFnameInConfig(zone_dir, dns_file))

			fmt.Println(target_list)
			for i := range target_list {
				target_addr := readFile(zone_dir, target_list[i])
				fmt.Printf("Checking file %s\n", target_list[i])
				log.WithField("file", target_list[i]).Info("Checking file")
				
				fmt.Printf("SSHing %s\n", target_addr)

				attemptConnect(bastionConn, port_list, pass_list, target_user, 
					target_key, target_addr)
			}
		} else {
			fmt.Println("Reading File...")
			log.WithField("file", single_file).Info("Read file")
			
			target_addr := readFile(zone_dir, single_file)
			fmt.Printf("SSHing %s\n", target_addr)

			attemptConnect(bastionConn, port_list, pass_list, target_user, 
				target_key, target_addr)
		} */
	}
}

// function to check error	
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initializeLogging(log_file string) {
	file, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Could not open log file: " + err.Error())
	}

	log.SetOutput(file)
}

func attemptConnect(bastionConn *ssh.Client, port_list, pass_list []string, target_user, 
	target_key, target_addr string) {
	
	for j := range port_list {
		fmt.Printf("Dialing target on port %s...\n", port_list[j])

		for k := range pass_list {
			log.WithFields(log.Fields{
				"target": target_addr,
				"port": port_list[j],
			}).Info("Dialing target from Bastion")
			// make new bastionConn every loop (if not, will get error handshake failed)
			conn, err := bastionConn.Dial("tcp", target_addr + ":" + port_list[j])

			if err != nil {
				fmt.Println("Dialing from Bastion Failed: ", err)
				log.Error("Unable to dial target: " + err.Error())
				
				break //prevent dialing on same port
			} else {
				
				config := sshConfig(target_user, pass_list[k], target_key)

				fmt.Printf("Trying password %s\n", pass_list[k])
				log.WithFields(log.Fields{
					"user": target_user,
					"password": pass_list[k],
					"port": port_list[j],
					"key": target_key,
				}).Info("New Client Connection started")
				ncc, chans, reqs, err := ssh.NewClientConn(conn, target_addr + ":" + 
				port_list[j], config)
					
				if err != nil {
					fmt.Println("newClientConn error: ", err)
					log.Error("newClientConn error: " + err.Error())
				} else {
					destClient := ssh.NewClient(ncc, chans, reqs)
					log.Info("New client successfully created")
					err := sshSession(destClient)

					if err != nil {
						fmt.Println("Cannot execute command", err)
					} else {
						log.WithFields(log.Fields{
							"user": target_user,
							"password": pass_list[k],
							"port": port_list[j],
							"key": target_key,
						}).Info("This is the correct param. Succeed")
					}
				}
			}
		}

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

func sshSession(conn *ssh.Client) error {
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
	
	return err
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
	
func getFileName(zone_dir string) <-chan string {
	ch := make(chan string)
	f, err := ioutil.ReadDir(zone_dir)
	check(err)

	go func() {
		for _, file := range f {
			ch <- file.Name()
		}
		close(ch)
	}()
	return ch
}

func difference(a<-chan string, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, i := range b {
		mb[i] = struct{}{}
	}

	var diff []string
	for i := range a {
		if _, found := mb[i]; !found {
			diff = append(diff, i)
		}
	}

	return diff
}

func readFile(zone_dir, fn string) string {
	f, err := os.Open(zone_dir + fn)

	check(err)

	defer f.Close()

	buf := make([]byte, 451)
	n, err := f.Read(buf)

	check(err)

	ip_string := strings.Join(parseString(string(buf[:n])), "")

	return ip_string

}

// This function reads ip address in the zones/db file
// use regex to read it
func parseString(buf string) []string {
	r, _ := regexp.Compile("\\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.)" +
	"{3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\b")

	ip_string := r.FindStringSubmatch(buf)

	return ip_string
}
