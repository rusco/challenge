# Challenge App Server 

_Webserver to serve the 3 HTTP Endpoints /top-zones, /zone-trips and /list-yellow ._

This app is writen in Go, to compile you have to install the [Go compiler](https://go.dev/dl/) on your server.
Please follow the instructions on the Go compiler download page.

This command:
```bash
go version
```
should result in an output like this:

```bash
pi@raspberrypi:~ $ go version
go version go1.19.5 linux/arm
pi@raspberrypi:~ $
```

## Compilation of the server

After successful instalation of the Go compiler change on the commandline into the "loader" directory and run the "go tidy" to download 
golang libraries and after completion of the download rund the "go build" command:

```bash
cd ./loader
go mod tidy
go: downloading golang.org/x/sys v0.4.0
...
```` 

```bash
go build
```

```bash
ls
```

There should now appear a go binary with the name "server" or "server.exe" (windows only). 
On linux systems make sure the binary is executable:

```bash
chmod +x server
```
The server package comes with a default "server.config.json" file where you can change the following 3 parameters:

1. logsql:              **true** (default) or false
2. port:                **8081** (default) or another port you might want to choose
3. openlocalbrowser:    **false** (default) or true. Use this setting for example on a windows server to open a browser

Executing the "server" command starts the server:

```bash
http://localhost:8081       
```` 

## Accessing the http endpoints with curl 

The file index.js in the folder "public" has listed some url tests (endpoints object) which can be tested with the included "index.html" via a Browser and also via curl, here follow some examples:

curl -X GET -H "Content-type: application/json" -H "Accept:  application/json" "http://localhost:8081/api/v1/dbversion"

curl -X POST -H "Content-type: application/json" -H "Accept:  application/json" "http://localhost:8081/api/v1/dbversion"

The last url paramter (starting with /api) should be changed accordingly, examples :

- '/api/v1/about' 
- '/api/v1/dbversion' 
- '/api/v1/top-zones/pickups' 
- '/api/v1/top-zones/dropoffs' 
- '/api/v1/errorexample' 
- '/api/v1/top-zones/errorexamle' 
- '/api/v1/zone-trips/7/2018-01-06' 
- '/api/v1/zone-trips/' 
- '/api/v1/zone-trips/errortest' 
- '/api/v1/list-yellow/limit=100' 
- '/api/v1/list-yellow/limit=10/offset=10' 
- '/api/v1/list-yellow' 
- '/api/v1/list-yellow/sort=pu_datetime.asc/sort=pu_locationid.asc' 
- '/api/v1/list-yellow/sort=pu_datetime.asc/sort=pu_locationid.desc' 
- '/api/v1/list-yellow/sort=pu_datetime.asc/limit=10/filter=pu_datetime:gt:2018-01-20/sort=pu_locationid.asc/filter=do_locationid:gte:100/offset=30' 
- '/api/v1/list-yellow/limit=10/offset=10/sort=pu_datetime.desc/sort=pu_locationid.desc' 


If you encounter any issues please contact the hotline: 

Jochen Rebhan, j.rebhan@gmail.com
