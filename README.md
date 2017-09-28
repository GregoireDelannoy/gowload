# gowload
Simple stupid dropbox-like application written in Go

The goal is to replace my current {Own|Next}Cloud on a gen 1 Rapsberry Pi, which is far too heavy for the board.

Goals

* Simple standalone application to be used behind a reverse proxy
* As few features as I can get with
* Does not have to scale to a lot of users
* Does not mess with files or access ; you can still FTP in, modify files, etc...
* Works without JS, progressive enhancement

Design tradeoffs

* Authentication is handled by Apache or Nginx used as a reverse
* Users and their directory have to be created by hand

Key features

* Directory browse
* Public links
* Whole folder download as .zip
* Upload of multiple files

Known bugs/security issues

* As the authentication is entirely handled by the reverse proxy, any local application that can issue requests to 127.0.0.1 with a X-Basic-Auth header set can impersonate an user. This does include any user that can port forward with SSH... TODO: replace it with unix sockets (should be trivial)
* This program exposes the local filesystem to the web. It probably has A LOT of issues and is dangerous to use!


## How to install/use

Build binary

* Modify constants (directories) and ports to suit your needs
* Go get any external dependency (As for now: "github.com/pierrre/archivefile/zip")
* Go build!

Configure the reverse proxy (example with Nginx)
```
location /files/ {
	auth_basic             "Restricted";
    auth_basic_user_file   /etc/nginx/htpasswd;
    proxy_pass             http://127.0.0.1:3000;
	proxy_set_header	X-Basic-Auth $http_authorization;
 	proxy_redirect      off;
}
```

```
location /links/ {
    proxy_pass             http://192.168.0.15:3000;
	proxy_set_header	X-Basic-Auth "";
 	proxy_redirect      off;
}
```



