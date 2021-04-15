# beacon-plus

this is the custom beacon server based on official Drand impl.

## demand
1. run a process as beacon server(http server) for our all fullnodes from different region to gain the beacon
2. add specify outsea net channel to call official Drand server, espically before the actual usage(30sec)
3. store beacon locally using level datastore
4. for clients add our cusotm beacon server

## proposal
1. http server using original http package
2. main loop if 30sec time interval to call api.deand.sh/api2.drand.sh/api3.drand.sh
3. serval goroutine to call different endpoint parallal
4. using leveldb to store round data
5. must verify jwt.token of incoming calling 

数据中心:提供可靠地机房