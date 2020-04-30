package stats

import (
         "strconv"
         "github.com/shirou/gopsutil/cpu"
         "github.com/shirou/gopsutil/mem"
         "log"
         "os"
         "encoding/csv"
         "time"
         // "github.com/jasonlvhit/gocron"
 )

type Stats struct{
	Time string
	CpuNum int
	TotalMem int
	FreeMem int
	MemPercent float32
	CpuPercents []float32

}

func New(Time string, CpuNum int, TotalMem int,FreeMem int, MemPercent float32, CpuPercents []float32) Stats {  
    s := Stats {Time, CpuNum , TotalMem ,FreeMem , MemPercent , CpuPercents}
    return s
}

func (s Stats) SendMetrics(){
 	vmStat, err := mem.VirtualMemory()
 	cpuStat, err := cpu.Info()
 	percentage, err := cpu.Percent(0, true)
 // 	for _,cpu := range percentage {
	// 	log.Println(strconv.FormatFloat(cpu, 'f', 6, 64))
	// }
 	

 	var filename ="logs/process" + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "log.csv"
 	f, err := os.OpenFile(filename,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	fi, err := f.Stat()
	csvwriter := csv.NewWriter(f)
	var row = []string{}
	//size := strconv.FormatInt(fi.Size(),10)
	// log.Println(size + "**********************")
	if fi.Size() == 0{
		row = []string{
		"Time","CPU NUM", "Total Mem(bytes)","Free Mem(bytes)", "Pecent Mem Usage", "CPU Util 0", "CPU Util 1", "CPU Util 2",
		"CPU Util 3","CPU Util 4","CPU Util 5","CPU Util 6","CPU Util 7","CPU Util 8","CPU Util 9",
		"CPU Util 10","CPU Util 11","CPU Util 12","CPU Util 13","CPU Util 14","CPU Util 15","CPU Util 16", 
		"CPU Util 17","CPU Util 18","CPU Util 19","CPU Util 20","CPU Util 21","CPU Util 22","CPU Util 23"}
		csvwriter.Write(row)
	}

	t:= time.Now()

	s.Time = t.String()
	s.CpuNum = cpuStat[0].CPU
	s.TotalMem = vmStat.Total
	s.FreeMem = vmStat.Free
	s.MemPercent = vmStat.UsedPercent
	s.CpuPercents = percentage

	row = []string{s.Time,strconv.FormatInt(int64(s.CpuNum, 10),strconv.FormatUint(s.TotalMem, 10),strconv.FormatUint(s.FreeMem, 10),strconv.FormatFloat(s.MemPercent, 'f', 2, 64)}
	for _,cpu := range percentage {
		row = append(row, strconv.FormatFloat(cpu, 'f', 2, 64))
	}

	csvwriter.Write(row)
	csvwriter.Flush()
	f.Close()
}

func ExecuteCronJob() {
 	for {
 		SendMetrics()
 		time.Sleep(2 * time.Second)
	}
    // gocron.Every(2).Second().Do(SendMetrics)
    // <- gocron.Start()
}