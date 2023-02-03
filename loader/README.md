# Challenge App Loader 

_Commanline interface application to load zone or green- and yellow trip data into the database ._

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

## Compilation of the loader

After successful instalation of the Go compiler change on the commandline into the "loader" directory and run the "go tidy" to download 
golang libraries and after completion of the download rund the "go build" command:

```bash
cd ./loader
go mod tidy
go: downloading modernc.org/sqlite v1.20.2
...
```` 

```bash
go build
```

```bash
ls
```

There should now appear a go binary with the name "loader[.exe]"  
On linux systems make sure the binary is executable:

```bash
chmod +x loader
```

Executing the "loader" command without arguments should show the required parameters:

```bash
pi@raspberrypi:~/go/challenge/loader $ ./loader

2 arguments expected: argument1 = zone|green|yellow, argument 2 = filename.csv
Examples:
 loader[.exe] zone   ../data/taxi_zone_lookup.csv
 loader[.exe] green  ../data/green_tripdata_2018-01_01-15.csv
 loader[.exe] yellow ../data/yellow_tripdata_2018-01_01-15.csv

pi@raspberrypi:~/go/challenge/loader $
```` 

## Loading Zone and Trip data 

Now download 3 csv file from this download link:

https://easyupload.io/m/xg063v

Unzip them into the "data" folder of your project, there should be now these 3 csv file: 

1. Zones list:   "zones.csv"
2. Green taxis:  "green_tripdata_2018-01_01-15.csv" 
3. Yellow taxis: "yellow_tripdata_2018-01_01-15.csv"


You can load them one after the other appending the filetype (zone or green or yellow) and filename parameter:

```bash
pi@raspberrypi:~/go/challenge/loader $ ./loader zone ../data/taxi_zone_lookup.csv
2023/01/22 16:32:37 done: 266 Records of type 'zone' loaded in 19.9776ms.

pi@raspberrypi:~/go/challenge/loader $ ./loader green  ../data/green_tripdata_2018-01_01-15.csv
2023/01/22 16:32:28 done: 367477 Records of type 'green' loaded in 5.1098583s.

pi@raspberrypi:~/go/challenge/loader $ ./loader yellow ../data/yellow_tripdata_2018-01_01-15.csv
2023/01/22 16:31:57 done: 3964004 Records of type 'yellow' loaded in 32.0927146s.
```

Repeated loading of each dataset is possible, in this case the previous data gets deleted. 

Have Fun !
