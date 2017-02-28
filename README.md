## Why?

for my own personal study of the mysql mariadb and how it handles transactions.  Can be used to generate transaction/certification conflicts

## usage

```
Usage of mysql-punch:
  -c int
    	total number of connections to create (default 2)
  -d string
    	database to use for testing (default "test")
  -i int
    	interval in miliseconds to sleep between sql statements (default 200)
  -p string
    	password defaults to blank
  -port string
    	mysql server port (default "3306")
  -q string
    	file containing list of queries to run per connection
  -s string
    	comma delimited list of servers to use (default "localhost")
  -u string
    	user default to root (default "root")
```

## Example success

```
cat queries.txt
update tblzone_time set time_in_mins = "22 mins" where key_value = "key value"
update tblzone_time set time_in_mins = "11 mins" where key_value = "key value"

Daniels-MBP-2:mysql-punch danl$ mysql-punch -c 100 -d test -i 200 -p password -q queries.txt -s 10.244.7.2
Starting connections...
Finished starting connections...
Listening for SEGTERM events...
Running...
```

## Example transaction conflict

```
cat queries.txt
update tblzone_time set time_in_mins = "22 mins" where key_value = "key value"
update tblzone_time set time_in_mins = "11 mins" where key_value = "key value"


Daniels-MBP-2:mysql-punch danl$ mysql-punch -c 100 -d test -i 200 -p password -q queries.txt -s 10.244.7.2,10.244.8.2,10.244.9.2
Starting connections...
Finished starting connections...
Listening for SEGTERM events...
Running...
Error 1317: Query execution was interrupted
Error 1317: Query execution was interrupted
Error 1317: Query execution was interrupted
```

### mysql logs

```
*** Priority TRANSACTION:
TRANSACTION 841845, ACTIVE 0 sec starting index read
mysql tables in use 1, locked 1
1 lock struct(s), heap size 360, 0 row lock(s)
MySQL thread id 1, OS thread handle 0x7f0b0b6c8700, query id 1569620 Update_rows_log_event::find_row(179653)

*** Victim TRANSACTION:
TRANSACTION 841844, ACTIVE 0 sec
2 lock struct(s), heap size 360, 2 row lock(s), undo log entries 1
MySQL thread id 5524, OS thread handle 0x7f0a40ef3700, query id 1569619 192.168.50.1 root closing tables
update tblzone_time set time_in_mins = "11 mins" where key_value = "key value"
*** WAITING FOR THIS LOCK TO BE GRANTED:
RECORD LOCKS space id 5 page no 3 n bits 72 index `GEN_CLUST_INDEX` of table `test`.`tblzone_time` trx table locks 1 total table locks 2  trx id 841844 lock_mode X lock hold time 0 wait time before grant 0
170228 20:05:18 [Note] WSREP: cluster conflict due to high priority abort for threads:
170228 20:05:18 [Note] WSREP: Winning thread:
   THD: 1, mode: applier, state: executing, conflict: no conflict, seqno: 179653
   SQL: (null)
170228 20:05:18 [Note] WSREP: Victim thread:
   THD: 5524, mode: local, state: idle, conflict: no conflict, seqno: -1
   SQL: (null)
170228 20:05:19 [Note] WSREP: cluster conflict due to certification failure for threads:
170228 20:05:19 [Note] WSREP: Victim thread:
   THD: 5529, mode: local, state: executing, conflict: cert failure, seqno: 179664
   SQL: COMMIT
```
