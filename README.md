# dssh
Simple console SSH client and connection manager.
## Description and usage
DSSH searchs hosts configuration in YAML files (extension *.hosts) located in folder %USERPROFILE%/.dssh.
It supports only authentication by password, no port forwarding and etc.

### Why
- Windows OpenSSH client does not support sends Backspace as ^H, but some devices support only ^H. It's borring to press ^H.
- Need to remember a lot of password. Auth by private key is not suitable for my situation.


### Configuration
Supports basic YAML format. Name of keys are case sensitive. Unknow keys are ignored.
Supports only following keys:
- Address
- UserName
- Password

If address value does not contains port value, then default value port=22 is used.

Below is example of hosts configuration:
```yaml
# file: sample.hosts
proxy_server:
  Address: 192.168.1.1
  UserName: ubuntu
  Password: pass

ftp_server:
  Address: 172.16.10.10:24
  UserName: anonymous
  Password: anonymous
```

At present moment supports only Windows platform.

DSSH provides completion for [Clink](<https://github.com/chrisant996/clink>) and can be installed by command:
```shell
$ dssh --install-completion
```

To show configuration of host:
```shell
$ dssh --show proxy_server
File: sample.config
proxy_server:
  Address: 192.168.1.1
  UserName: ubuntu
  Password: pass
```

Connect to host:
```
$ dssh proxy_server
ubuntu@ubuntu:~#
```

## TODO
- [ ] May be change to different hosts configuraion file format as YAML depends on indentation
- [ ] Store password encrypted or use KeyPass
- [ ] Add support for Linux
