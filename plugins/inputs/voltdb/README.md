# VoltDB Input Plugin

This plugin gathers the following statistics from a VoltDB instance

* CPU_host_id - The utilization % on each host
* LATENCY_host_id - Metrics on latencies of stored procedure calls on each host
* MEMORY_host_id - The utilization % statistics on each host
* QUEUE_host_id_site_id - Metrics on the stored procedure call queues at each site on each host
* IDLETIME_host_id_site_id - The ideletime of each site on each host

### Configuration

```toml
[[inputs.voltdb]]
  ## Specify hostnames of VoltDB nodes (physical/virtual/container) using a comma-separated string of hostnames
  hosts="voltdb-1,voltdb-2,voltdb-3"
  
  ## Specify the port that VoltDB listens on for calls using its wire protocol
  port=21212
  
  ## Specify the procedure to be called
  proc=\@Statistics (only procedure supported now)
  
  ## Specify the delta for the accumulation of the statistics on VoltDB
  ## Check https://docs.voltdb.com/UsingVoltDB/sysprocstatistics.php
  delta=0
```

### Metrics
Refer to VoltDB's documentation at https://docs.voltdb.com/UsingVoltDB/sysprocstatistics.php for the metrics returned
  