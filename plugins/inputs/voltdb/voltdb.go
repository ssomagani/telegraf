package voltdb

import (
    
    "fmt"
    "database/sql/driver"
    "strconv"
    "time"
    
    "github.com/VoltDB/voltdb-client-go/voltdbclient"
    
    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type VoltDB struct {
    Host string `toml:"host"`
    Proc string `toml:"proc"`
    Delta int64 `toml:"delta"`
    
    CPU CPUStats
    Queue QueueStats
}

type CPUStats struct {
    HOSTNAME string `toml:"hostname"`
    PERCENT_USED int64 `toml:"percent_used"`
}

type QueueStats struct {
    TIMESTAMP int64
    HOST_ID int32
    HOSTNAME string
    SITE_ID int32
    CURRENT_DEPTH int32
    POLL_COUNT int64
    AVG_WAIT int64
    MAX_WAIT int64
}

func (s *VoltDB) Description() string {
    return "a demo plugin"
}

func (s *VoltDB) SampleConfig() string {
    return `
  ## Ask for a procedure to be run
  proc = SystemInformation
`
}

func (s *VoltDB) Init() error {
	return nil
}

func (s *VoltDB) accQueueStats(acc telegraf.Accumulator, conn *voltdbclient.Conn) error {
    result, err := conn.Query(s.Proc, []driver.Value{"QUEUE", s.Delta})
	if err != nil {
        return err
    }
    
    voltRows := *result.(*voltdbclient.VoltRows)
	for voltRows.AdvanceRow() {
		timestamp, err := voltRows.GetBigIntByName("TIMESTAMP")
		if err != nil {
            return err
		}
		hostId, err := voltRows.GetIntegerByName("HOST_ID")
		if err != nil {
            return err
		}

        siteId, err := voltRows.GetIntegerByName("SITE_ID")
		if err != nil {
            return err
		}
        currentDepth, err := voltRows.GetIntegerByName("CURRENT_DEPTH")
		if err != nil {
            return err
		}
        pollCount, err := voltRows.GetBigIntByName("POLL_COUNT")
		if err != nil {
            return err
		}
        avgWait, err := voltRows.GetBigIntByName("AVG_WAIT")
		if err != nil {
            return err
		}
        maxWait, err := voltRows.GetBigIntByName("MAX_WAIT")
		if err != nil {
            return err
		}
        
        fields := make(map[string]interface{})
        fields["CURRENT_DEPTH"] = fmt.Sprint(currentDepth.(int32))
        fields["POLL_COUNT"] = strconv.FormatInt(pollCount.(int64), 10)
        fields["AVG_WAIT"] = strconv.FormatInt(avgWait.(int64), 10)
        fields["MAX_WAIT"] = strconv.FormatInt(maxWait.(int64), 10)
        t := timestamp.(int64) * 1000000
        acc.AddFields("QUEUE_" + fmt.Sprint(hostId.(int32)) + "_" + fmt.Sprint(siteId.(int32)), fields, nil, time.Unix(0, t).UTC())
	}
    return nil
}

func (s *VoltDB) accCPUStats (acc telegraf.Accumulator, conn *voltdbclient.Conn) error {
    result, err := conn.Query(s.Proc, []driver.Value{"CPU", s.Delta})
	if err != nil {
        return err
    }
    
    voltRows := *result.(*voltdbclient.VoltRows)
	for voltRows.AdvanceRow() {
        timestamp, err := voltRows.GetBigIntByName("TIMESTAMP")
		if err != nil {
            return err
		}
        
		hostId, err := voltRows.GetIntegerByName("HOST_ID")
		if err != nil {
            return err
		}
		percentUsed, err := voltRows.GetBigIntByName("PERCENT_USED")
		if err != nil {
            return err
		}
        t := timestamp.(int64) * 1000000
        acc.AddFields("CPU_" + fmt.Sprint(hostId.(int32)), map[string]interface{}{"PERCENT_USED": strconv.FormatInt(percentUsed.(int64), 10)}, nil, time.Unix(0, t).UTC()) 
	}
    return nil
}

func (s *VoltDB) accLatencyStats (acc telegraf.Accumulator, conn *voltdbclient.Conn) error {
    result, err := conn.Query(s.Proc, []driver.Value{"LATENCY", s.Delta})
	if err != nil {
        return err
    }
    
    voltRows := *result.(*voltdbclient.VoltRows)
	for voltRows.AdvanceRow() {
        timestamp, err := voltRows.GetBigIntByName("TIMESTAMP")
		if err != nil {
            return err
		}
        
		hostId, err := voltRows.GetIntegerByName("HOST_ID")
		if err != nil {
            return err
		}
        
		interval, err := voltRows.GetIntegerByName("INTERVAL")
		if err != nil {
            return err
		}
        
        count, err := voltRows.GetIntegerByName("COUNT")
		if err != nil {
            return err
		}
        
        tps, err := voltRows.GetIntegerByName("TPS")
		if err != nil {
            return err
		}
        
        p50, err := voltRows.GetBigIntByName("P50")
		if err != nil {
            return err
		}
        
        p95, err := voltRows.GetBigIntByName("P95")
		if err != nil {
            return err
		}
        
        p99, err := voltRows.GetBigIntByName("P99")
		if err != nil {
            return err
		}
        
        p99_9, err := voltRows.GetBigIntByName("P99.9")
		if err != nil {
            return err
		}
        
        p99_99, err := voltRows.GetBigIntByName("P99.99")
		if err != nil {
            return err
		}
        
        p99_999, err := voltRows.GetBigIntByName("P99.999")
		if err != nil {
            return err
		}
        
        max, err := voltRows.GetBigIntByName("MAX")
		if err != nil {
            return err
		}
        
        fields := make(map[string]interface{})
        fields["INTERVAL"] = fmt.Sprint(interval.(int32))
        fields["COUNT"] = fmt.Sprint(count.(int32))
        fields["TPS"] = fmt.Sprint(tps.(int32))
        fields["P50"] = strconv.FormatInt(p50.(int64), 10)
        fields["P95"] = strconv.FormatInt(p95.(int64), 10)
        fields["P99"] = strconv.FormatInt(p99.(int64), 10)
        fields["P99.9"] = strconv.FormatInt(p99_9.(int64), 10)
        fields["P99.99"] = strconv.FormatInt(p99_99.(int64), 10)
        fields["P99.999"] = strconv.FormatInt(p99_999.(int64), 10)
        fields["MAX"] = strconv.FormatInt(max.(int64), 10)
        t := timestamp.(int64) * 1000000
        acc.AddFields("LATENCY_" + fmt.Sprint(hostId.(int32)), fields, nil, time.Unix(0, t).UTC()) 
	}
    return nil
}

func (s *VoltDB) accMemoryStats (acc telegraf.Accumulator, conn *voltdbclient.Conn) error {
    
    result, err := conn.Query(s.Proc, []driver.Value{"MEMORY", s.Delta})
	if err != nil {
        return err
    }

    voltRows := *result.(*voltdbclient.VoltRows)
	for voltRows.AdvanceRow() {
        
        timestamp, err := voltRows.GetBigIntByName("TIMESTAMP")
		if err != nil {
            return err
		}
        
		hostId, err := voltRows.GetIntegerByName("HOST_ID")
		if err != nil {
            return err
		}

		rss, err := voltRows.GetIntegerByName("RSS")
		if err != nil {
            return err
		}
        
        javaUsed, err := voltRows.GetIntegerByName("JAVAUSED")
		if err != nil {
            return err
		}
                
        javaUnused, err := voltRows.GetIntegerByName("JAVAUNUSED")
		if err != nil {
            return err
		}
                
        tupleData, err := voltRows.GetBigIntByName("TUPLEDATA")
		if err != nil {
            return err
		}
        
        tupleAllocated, err := voltRows.GetBigIntByName("TUPLEALLOCATED")
		if err != nil {
            return err
		}
    
        indexMemory, err := voltRows.GetBigIntByName("INDEXMEMORY")
		if err != nil {
            return err
		}
    
        stringMemory, err := voltRows.GetBigIntByName("STRINGMEMORY")
		if err != nil {
            return err
		}
                
        tupleCount, err := voltRows.GetBigIntByName("TUPLECOUNT")
		if err != nil {
            return err
		}
              
        pooledMemory, err := voltRows.GetBigIntByName("POOLEDMEMORY")
		if err != nil {
            return err
		}
           
        physicalMemory, err := voltRows.GetBigIntByName("PHYSICALMEMORY")
		if err != nil {
            return err
		}

        javaMaxHeap, err := voltRows.GetIntegerByName("JAVAMAXHEAP")
		if err != nil {
            return err
		}
        
        fields := make(map[string]interface{})
        fields["RSS"] = fmt.Sprint(rss.(int32))
        fields["JAVAUSED"] = fmt.Sprint(javaUsed.(int32))
        fields["JAVAUNUSED"] = fmt.Sprint(javaUnused.(int32))
        fields["TUPLEDATA"] = strconv.FormatInt(tupleData.(int64), 10)
        fields["TUPLEALLOCATED"] = strconv.FormatInt(tupleAllocated.(int64), 10)
        fields["INDEXMEMORY"] = strconv.FormatInt(indexMemory.(int64), 10)
        fields["STRINGMEMORY"] = strconv.FormatInt(stringMemory.(int64), 10)
        fields["TUPLECOUNT"] = strconv.FormatInt(tupleCount.(int64), 10)
        fields["POOLEDMEMORY"] = strconv.FormatInt(pooledMemory.(int64), 10)
        fields["PHYSICALMEMORY"] = strconv.FormatInt(physicalMemory.(int64), 10)
        fields["JAVAMAXHEAP"] = fmt.Sprint(javaMaxHeap.(int32))
        t := timestamp.(int64) * 1000000
        acc.AddFields("MEMORY_" + fmt.Sprint(hostId.(int32)), fields, nil, time.Unix(0, t).UTC())
	}
    return nil
}

func (s *VoltDB) accIdleTimeStats (acc telegraf.Accumulator, conn *voltdbclient.Conn) error {
    result, err := conn.Query(s.Proc, []driver.Value{"IDLETIME", s.Delta})
	if err != nil {
        return err
    }
    
    voltRows := *result.(*voltdbclient.VoltRows)
	for voltRows.AdvanceRow() {
        timestamp, err := voltRows.GetBigIntByName("TIMESTAMP")
		if err != nil {
            return err
		}
        
		hostId, err := voltRows.GetIntegerByName("HOST_ID")
		if err != nil {
            return err
		}
        
        siteId, err := voltRows.GetIntegerByName("SITE_ID")
		if err != nil {
            return err
		}
        
		count, err := voltRows.GetBigIntByName("COUNT")
		if err != nil {
            return err
		}
        
        percent, err := voltRows.GetFloatByName("PERCENT")
		if err != nil {
            return err
		}
        
        avg, err := voltRows.GetFloatByName("AVG")
		if err != nil {
            return err
		}
        
        min, err := voltRows.GetFloatByName("MIN")
		if err != nil {
            return err
		}
        
        max, err := voltRows.GetFloatByName("MAX")
		if err != nil {
            return err
		}
        
        stdDev, err := voltRows.GetFloatByName("STDDEV")
		if err != nil {
            return err
		}
        
        fields := make(map[string]interface{})
        fields["COUNT"] = strconv.FormatInt(count.(int64), 10)
        fields["PERCENT"] = fmt.Sprintf("%f", percent)
        fields["AVG"] = fmt.Sprintf("%f", avg)
        fields["MIN"] = fmt.Sprintf("%f", min)
        fields["MAX"] = fmt.Sprintf("%f", max)
        fields["STDDEV"] = fmt.Sprintf("%f", stdDev)
        
        t := timestamp.(int64) * 1000000
        acc.AddFields("IDLETIME_" + fmt.Sprint(hostId.(int32)) + "_" + fmt.Sprint(siteId.(int32)), fields, nil, time.Unix(0, t).UTC())
	}
    return nil
}

func (s *VoltDB) Gather(acc telegraf.Accumulator) error {
    conn, err := voltdbclient.OpenConn("voltdb://"+s.Host+":21212")
	if err != nil {
		return err
	}
	defer conn.Close()
 
    s.accCPUStats(acc, conn)
    s.accQueueStats(acc, conn)
    s.accLatencyStats(acc, conn)
    s.accMemoryStats(acc, conn)
    s.accIdleTimeStats(acc, conn)
    return nil
}

func init() {
    inputs.Add("voltdb", func() telegraf.Input { return &VoltDB{} })
}
