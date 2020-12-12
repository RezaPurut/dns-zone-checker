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
