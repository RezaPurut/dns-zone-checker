# dns-zone-checker
Simple tool to read DNS forward zone file, compare it with DNS config file and ssh into the server. It is designed to read the file from current host/DNS server. Then, ssh into Bastion host and then ssh into target server.

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
### Options
```
--bulk boolean             If this option is enabled, it will compare all files in zones directory with DNS configuration file
--dns-file string          DNS configuration file to read and compare. Specifies the DNS config file (eg. named.conf, named.conf.default-zones, zones.conf, etc.)
--zone-dir string          Zone files directory. Specifies the directory that store the zone files such as example.com.zone, db.example.com)
--bastion-addr string      Address or hostname of the bastion server
--bastion-key string       SSH private key path for bastion server
--bastion-user string      Username to connect to bastion server
--bastion-pass string      Password for the user to connect to bastion server
--bastion-port string      Port to connect to bastion server
```
