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
--bastion-addr string      Address or hostname of the bastion server
--bastion-key string       SSH private key path for bastion server
--bastion-user string      Username to connect to bastion server
--bastion-pass string      Password for the user to connect to bastion server
--bastion-port string      Port to connect to bastion server
```
