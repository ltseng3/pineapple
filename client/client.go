package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"runtime"
	"sync"
	"time"

	"pineapple/genericsmrproto"
	"pineapple/masterproto"
	"pineapple/poisson"
	"pineapple/state"
	"pineapple/zipfian"

	"golang.org/x/sync/semaphore"
)

var masterAddr *string = flag.String("maddr", "", "Master address. Defaults to localhost")
var masterPort *int = flag.Int("mport", 7087, "Master port.")
var procs *int = flag.Int("p", 2, "GOMAXPROCS.")
var conflicts *int = flag.Int("c", 0, "Percentage of conflicts. If -1, uses Zipfian distribution.")
var forceLeader = flag.Int("l", -1, "Force client to talk to a certain replica.")
var startRange = flag.Int("sr", 0, "Key range start")
var T = flag.Int("T", 1, "Number of threads (simulated clients).")
var outstandingReqs = flag.Int64("or", 1, "Number of outstanding requests a thread can have at any given time.")
var theta = flag.Float64("theta", 0.99, "Theta zipfian parameter")
var zKeys = flag.Uint64("z", 1e9, "Number of unique keys in zipfian distribution.")
var poissonAvg = flag.Int("poisson", -1, "The average number of microseconds between requests. -1 disables Poisson.")
var percentWrites = flag.Float64("writes", 1, "A float between 0 and 1 that corresponds to the percentage of requests that should be writes. The remainder will be reads.")
var percentRMWs = flag.Float64("rmws", 0, "A float between 0 and 1 that corresponds to the percentage of writes that should be RMWs. The remainder will be regular writes.")
var blindWrites = flag.Bool("blindwrites", false, "True if writes don't need to execute before clients receive responses.")
var singleClusterTest = flag.Bool("singleClusterTest", true, "True if clients run on a VM in a single cluster")
var rampDown *int = flag.Int("rampDown", 15, "Length of the cool-down period after statistics are measured (in seconds).")
var rampUp *int = flag.Int("rampUp", 15, "Length of the warm-up period before statistics are measured (in seconds).")
var timeout *int = flag.Int("timeout", 180, "Length of the timeout used when running the client")

// Information about the latency of an operation
type response struct {
	receivedAt    time.Time
	rtt           float64 // The operation latency, in ms
	commitLatency float64 // The operation's commit latency, in ms
	isRead        bool
	replicaID     int
}

// Information pertaining to operations that have been issued but that have not
// yet received responses
type outstandingRequestInfo struct {
	sync.Mutex
	sema       *semaphore.Weighted // Controls number of outstanding operations
	startTimes map[int32]time.Time // The time at which operations were sent out
	isRead     map[int32]bool
}

// An outstandingRequestInfo per client thread
var orInfos []*outstandingRequestInfo

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	if *conflicts > 100 {
		log.Fatalf("Conflicts percentage must be between 0 and 100.\n")
	}

	orInfos = make([]*outstandingRequestInfo, *T)

	var master *rpc.Client
	var err error
	for {
		master, err = rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", *masterAddr, *masterPort))
		if err != nil {
			log.Println("Error connecting to master", err)
		} else {
			break
		}
	}

	rlReply := new(masterproto.GetReplicaListReply)
	for !rlReply.Ready {
		err := master.Call("Master.GetReplicaList", new(masterproto.GetReplicaListArgs), rlReply)
		if err != nil {
			log.Println("Error making the GetReplicaList RPC", err)
		}
	}

	leader := 0
	if *forceLeader < 0 {
		reply := new(masterproto.GetLeaderReply)
		if err = master.Call("Master.GetLeader", new(masterproto.GetLeaderArgs), reply); err != nil {
			log.Println("Error making the GetLeader RPC:", err)
		}
		leader = reply.LeaderId
		log.Printf("The leader is replica %d\n", leader)
	} else {
		leader = *forceLeader
	}

	readings := make(chan *response, 100000)

	//startTime := rand.New(rand.NewSource(time.Now().UnixNano()))
	experimentStart := time.Now()

	for i := 0; i < *T; i++ {

		// automatically allocate clients equally
		if *singleClusterTest {
			leader = i % len(rlReply.ReplicaList)
		}

		server, err := net.Dial("tcp", rlReply.ReplicaList[leader])
		if err != nil {
			log.Fatalf("Error connecting to replica %d\n", leader)
		}
		log.Print("Connected to replica ", leader, " at ", rlReply.ReplicaList[leader], "\n")

		reader := bufio.NewReader(server)
		writer := bufio.NewWriter(server)

		orInfo := &outstandingRequestInfo{
			sync.Mutex{},
			semaphore.NewWeighted(*outstandingReqs),
			make(map[int32]time.Time, *outstandingReqs),
			make(map[int32]bool, *outstandingReqs)}

		//waitTime := startTime.Intn(3)
		//time.Sleep(time.Duration(waitTime) * 100 * 1e6)
		go simulatedClientWriter(writer, orInfo)
		go simulatedClientReader(reader, orInfo, readings, leader)

		orInfos[i] = orInfo
	}
	log.Println("Here 142")
	if *singleClusterTest {
		log.Println("Here 144")
		printerMultipleFile(readings, len(rlReply.ReplicaList), experimentStart, rampDown, rampUp, timeout)
	} else {
		log.Println("Here 147")
		printer(readings)
	}
}

func simulatedClientWriter(writer *bufio.Writer, orInfo *outstandingRequestInfo) {
	args := genericsmrproto.Propose{
		CommandId: 0,
		Command:   state.Command{Op: state.PUT, K: 0, V: 0},
		Timestamp: 0,
	} // @audit autodetermine proposal type

	conflictRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	zipf := zipfian.NewZipfianGenerator(*zKeys, *theta)
	poissonGenerator := poisson.NewPoisson(*poissonAvg)
	opRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	queuedReqs := 0 // The number of poisson departures that have been missed

	for id := int32(0); ; id++ {
		log.Println("Here 165, id: ", id)
		args.CommandId = id

		// Determine key
		if *conflicts >= 0 {
			r := conflictRand.Intn(100)
			if r < *conflicts {
				args.Command.K = 42
			} else {
				//args.Command.K = state.Key(*startRange + 43 + int(id % 888))
				args.Command.K = state.Key(int32(*startRange) + 43 + id)
			}
			log.Println("Here 179, id: ", id)
		} else {
			args.Command.K = state.Key(zipf.NextNumber())
			log.Println("Here 182")
		}

		// Determine operation type
		if *percentWrites+*percentRMWs > opRand.Float64() {
			if *percentWrites > opRand.Float64() {
				if !*blindWrites {
					args.Command.Op = state.PUT // write operation
				} else {
					//args.Command.Op = state.PUT_BLIND
				}
			} else {
				args.Command.Op = state.RMW // RMW operation
			}
		} else {
			args.Command.Op = state.GET // read operation
		}
		log.Println("Here 199, type: ", args.Command.Op, " id: ", id)

		if *poissonAvg == -1 { // Poisson disabled
			orInfo.sema.Acquire(context.Background(), 1)
		} else {
			for {
				log.Println("Here 206 id: ", id)
				if orInfo.sema.TryAcquire(1) {
					if queuedReqs == 0 {
						time.Sleep(poissonGenerator.NextArrival())
					} else {
						queuedReqs -= 1
					}
					break
				}
				time.Sleep(poissonGenerator.NextArrival())
				queuedReqs += 1
			}
		}
		log.Println("Here 219 id: ", id)

		before := time.Now()
		writer.WriteByte(genericsmrproto.PROPOSE)
		args.Marshal(writer)
		writer.Flush()

		orInfo.Lock()
		if args.Command.Op == state.GET {
			orInfo.isRead[id] = true
		}
		orInfo.startTimes[id] = before
		orInfo.Unlock()
		log.Println("Here 232 id: ", id)
	}
}

func simulatedClientReader(reader *bufio.Reader, orInfo *outstandingRequestInfo, readings chan *response, leader int) {
	var reply genericsmrproto.ProposeReplyTS

	for {
		if err := reply.Unmarshal(reader); err != nil || reply.OK == 0 {
			if err != nil {
				log.Println("Error during unmarshaling:", err)
			} else if reply.OK == 0 {
				log.Println("reply.OK is 0")
			}
			log.Println(reply.OK)
			log.Println(reply.CommandId)
			break
		}

		after := time.Now()
		orInfo.sema.Release(1)

		orInfo.Lock()
		before := orInfo.startTimes[reply.CommandId]
		isRead := orInfo.isRead[reply.CommandId]
		delete(orInfo.startTimes, reply.CommandId)
		orInfo.Unlock()

		rtt := (after.Sub(before)).Seconds() * 1000
		//commitToExec := float64(reply.Timestamp) / 1e6
		commitLatency := float64(0) //rtt - commitToExec

		readings <- &response{
			after,
			rtt,
			commitLatency,
			isRead,
			leader,
		}
	}
	log.Println("Here 264")
}

func printer(readings chan *response) {
	log.Println("Here 268")
	lattputFile, err := os.Create("lattput.txt")
	if err != nil {
		log.Println("Error creating lattput file", err)
		return
	}
	//lattputFile.WriteString("# time (ns), avg lat over the past second, tput since last line, total count, totalOrs, avg commit lat over the past second\n")

	latFile, err := os.Create("latency.txt")
	if err != nil {
		log.Println("Error creating latency file", err)
		return
	}
	//latFile.WriteString("# time (ns), latency, commit latency\n")

	startTime := time.Now()

	for {
		log.Println("Here 286")
		time.Sleep(time.Second)
		count := len(readings)
		var sum float64 = 0
		var commitSum float64 = 0
		endTime := time.Now() // Set to current time in case there are no readings
		log.Println("Here 292, count: ", count)
		for i := 0; i < count; i++ {
			resp := <-readings
			log.Println("here 295, ", resp)
			// Log all to latency file
			latFile.WriteString(fmt.Sprintf("%d %f %f\n", resp.receivedAt.UnixNano(), resp.rtt, resp.commitLatency))
			log.Print("Latency: ", resp.rtt, " commit latency: ", resp.commitLatency, " isRead: ", resp.isRead, "\n")
			sum += resp.rtt
			commitSum += resp.commitLatency
			endTime = resp.receivedAt
		}
		log.Println("Here 303")
		var avg float64
		var avgCommit float64
		var tput float64
		if count > 0 {
			avg = sum / float64(count)
			avgCommit = commitSum / float64(count)
			tput = float64(count) / endTime.Sub(startTime).Seconds()
		}

		totalOrs := 0
		for i := 0; i < *T; i++ {
			orInfos[i].Lock()
			totalOrs += len(orInfos[i].startTimes)
			orInfos[i].Unlock()
		}

		// Log summary to lattput file
		lattputFile.WriteString(fmt.Sprintf("%d %f %f %d %d %f\n", endTime.UnixNano(),
			avg, tput, count, totalOrs, avgCommit))

		startTime = endTime
	}
}

func printerMultipleFile(readings chan *response, numLeader int, experimentStart time.Time, rampDown, rampUp, timeout *int) {
	lattputFile, err := os.Create("lattput.txt")
	if err != nil {
		log.Println("Error creating lattput file", err)
		return
	}

	latFileRead := make([]*os.File, numLeader)
	latFileWrite := make([]*os.File, numLeader)
	log.Println("Here 339")
	for i := 0; i < numLeader; i++ {
		fileName := fmt.Sprintf("latFileRead-%d.txt", i)
		latFileRead[i], err = os.Create(fileName)
		if err != nil {
			log.Println("Error creating latency file", err)
			return
		}
		// latFile.WriteString("# time (ns), latency, commit latency\n")

		fileName = fmt.Sprintf("latFileWrite-%d.txt", i)
		latFileWrite[i], err = os.Create(fileName)
		if err != nil {
			log.Println("Error creating latency file", err)
			return
		}
		log.Println("Here 355, i: ", i)
	}

	startTime := time.Now()
	log.Println("Here 359")
	for {
		time.Sleep(time.Second)

		count := len(readings)
		var sum float64 = 0
		var commitSum float64 = 0
		endTime := time.Now() // Set to current time in case there are no readings
		currentRuntime := time.Now().Sub(experimentStart)
		for i := 0; i < count; i++ {
			log.Println("Here 369, i: ", i)
			resp := <-readings

			// Log all to latency file if they are not within the ramp up or ramp down period.
			if *rampUp < int(currentRuntime.Seconds()) && int(currentRuntime.Seconds()) < *timeout-*rampDown {
				if resp.isRead {
					latFileRead[resp.replicaID].WriteString(fmt.Sprintf("%d %f %f\n", resp.receivedAt.UnixNano(), resp.rtt, resp.commitLatency))
				} else {
					latFileWrite[resp.replicaID].WriteString(fmt.Sprintf("%d %f %f\n", resp.receivedAt.UnixNano(), resp.rtt, resp.commitLatency))
				}
				sum += resp.rtt
				commitSum += resp.commitLatency
				endTime = resp.receivedAt
			}
		}

		var avg float64
		var avgCommit float64
		var tput float64
		if count > 0 {
			avg = sum / float64(count)
			avgCommit = commitSum / float64(count)
			tput = float64(count) / endTime.Sub(startTime).Seconds()
		}
		log.Println("Here 393, tput: ", tput)

		totalOrs := 0
		for i := 0; i < *T; i++ {
			orInfos[i].Lock()
			totalOrs += len(orInfos[i].startTimes)
			orInfos[i].Unlock()
		}

		// Log summary to lattput file
		//lattputFile.WriteString(fmt.Sprintf("%d %f %f %d %d %f\n", endTime.UnixNano(), avg, tput, count, totalOrs, avgCommit))
		// Log all to latency file if they are not within the ramp up or ramp down period.
		if *rampUp < int(currentRuntime.Seconds()) && int(currentRuntime.Seconds()) < *timeout-*rampDown {
			lattputFile.WriteString(fmt.Sprintf("%d %f %f %d %d %f\n", endTime.UnixNano(), avg, tput, count, totalOrs, avgCommit))
			log.Println("Here 407")
		}
		startTime = endTime
	}
}
