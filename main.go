package main
// TODO: comment function, create file check function
import (
	"fmt"
	"flag"
	"os"
	"regexp"
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
		log_file string
	)

	flag.StringVar(&zone_file, "zone-file", "", "Zone file/files")
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
	
	flag.StringVar(&log_file, "log-file", "/var/log/dns-check.txt", "Log file")
	flag.Parse()

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
			file_list := strings.Split(zone_file, ",")
			for file := range(file_list) {
				target_addr := readFile(zone_dir, file_list[file])
	
				fmt.Printf("SSHing %s\n", target_addr)
	
				attemptConnect(bastionConn, port_list, pass_list, target_user, 
					target_key, target_addr)
			}
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
