

# Backend Challenge  

_Loader and Webserver for the Backend challenge according the specification in the docs folder._


This package contains 2 deliveries, one each in the folder "loader" and "server".
Please refer to the README.md files in these folders.

A raspberry-pi hosted version is (currently) available at: http://rebhan.duckdns.org (slow hardware!)


![image](https://user-images.githubusercontent.com/91183/213931562-1fb60caf-f930-4beb-a122-a7e196825fb1.png)

## Roadmap

For the next version 1.1 significant performance enhancements regarding the sql statements "top-zones" and "zone-trips" are planned. 

The version 1.1 database schema will include   integrate pu_zone and do_zone in a new "alltrips" table.

### Version 1.1 Database Changes:

```sql
-- migration scripts from current schema to new schema: 
 CREATE TABLE IF NOT EXISTS alltrips (
    pu_datetime     TEXT,
    do_datetime     TEXT,
    pu_locationid   INTEGER,
    do_locationid   INTEGER,
    color           TEXT,
    pu_locationzone TEXT,
    do_locationzone TEXT
);

CREATE INDEX pu_location_idx ON alltrips (
    pu_locationid
);
CREATE INDEX do_location_idx ON alltrips (
    do_locationid
);

INSERT INTO alltrips 
    SELECT t.pu_datetime,
        t.do_datetime,
        t.pu_locationid,
        t.do_locationid,
        t.color,
        z.zone pu_locationzone,
        z2.zone do_locationzone
    FROM trips t
        LEFT JOIN
        zones z ON z.locationid = t.pu_locationid
        LEFT JOIN
        zones z2 ON z2.locationid = t.do_locationid;

--afterwards: 
DROP TABLE trips; 
DROP TABLE zones; 




--new top-zones sqlt statement v1.1 : 
 WITH pucount AS (
    SELECT pu_locationzone,
           pu_locationid,
           count(pu_locationid) pu_total
      FROM alltrips t
     GROUP BY pu_locationid
     ORDER BY count(pu_locationid) DESC
),
docount AS (
    SELECT do_locationzone,
           do_locationid,
           count(do_locationid) do_total
      FROM alltrips t
     GROUP BY do_locationid
     ORDER BY count(do_locationid) DESC
)
SELECT *
  FROM pucount
       LEFT JOIN
       docount ON pucount.pu_locationid = docount.do_locationid-- 
 --ORDER BY pu_total DESC
 ORDER BY do_total DESC
  
 LIMIT 5;
 
 --new zone-trips sql statement v1.1:
WITH pus AS (
    SELECT a.pu_locationzone zone,
           date('2018-01-12') date,
           count(a.pu_locationid) pu
      FROM alltrips a
     WHERE a.pu_locationid = 36 AND 
           date(a.pu_datetime) = date('2018-01-12') 
),
dos AS (
    SELECT count(a.do_locationid) do_
      FROM alltrips a
     WHERE a.do_locationid = 36 AND 
           date(a.do_datetime) = date('2018-01-12') 
)
SELECT *
  FROM pus,
       dos;
	   
 --list-yellow statement v1.1:
 SELECT pu_datetime,
       do_datetime,
       pu_locationid,
       do_locationid
  FROM alltrips t
 WHERE COLOR = 'yellow'
 --- ...

```

Have Fun !



