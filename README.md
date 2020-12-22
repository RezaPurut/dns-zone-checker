# dns-zone-checker
Simple tool to read DNS forward zone file, compare it with DNS config file and ssh into the server. It is designed to read the file from current host/DNS server. Then, it will use Bastion server as a jump host to ssh into target server.

```
                       ----------         -------------        -------------------
                       -        -         -           -        -                 -
                       -  Host  -   -->   -  Bastion  -  -->   -  Target Server  -
                       -        -         -           -        -                 -
                       ----------         -------------        -------------------
```
## Usage
As binary:
```
./checker [flags]
```
OR

As go program:
```
go run main.go [flags]
```
### Options
```
--bulk boolean             If this option is enabled, it will compare all files in zones directory 
                           with DNS configuration file
--dns-file string          DNS configuration file to read and compare. Specifies the DNS config file 
                           (eg. named.conf, named.conf.default-zones, zones.conf, etc.)
--zone-dir string          Zone files directory. Specifies the directory that contains the zone files 
                           such as example.com.zone, db.example.com)
--single-zone string       Use this only when you want to check one zone file. This is used when 
                           bulk=false or 'bulk' is not provided
--bastion-addr string      Address or hostname of the bastion server
--bastion-key string       SSH private key path for bastion server
--bastion-user string      Username to connect to bastion server
--bastion-pass string      Password for the user to connect to bastion server
--bastion-port string      SSH port to connect to bastion server
--target-user string       Username to connect to target server
--target-pass string       Password for the target-user to connect to target server. Can provide multiple values
                           separated by comma (eg. "password,pass,abc"). It will try each of the provided password
                           to connect to the target server
--target-port string       SSH port for target server. Can provide multiple values separated by comma (eg. "22,2222")
--target-key string        SSH private key path for target server 
```
### Example 1
Bulk check
```
./checker -bulk=true -dns-file /etc/named/zones.conf -zone-dir /etc/named/zones/ -bastion-addr jumphost.example.com \
-bastion-user bastionUser -bastion-port 22 -bastion-key /home/bastionUser/.ssh/id_rsa -target-user targetUser \
-target-pass="pass,pass123" -target-port="22,2222" -target-key /home/targetUser/.ssh/id_rsa
```
