package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"pineapple/src/masterproto"
	"pineapple/src/pineapple"
)

var masterAddr *string = flag.String("maddr", "10.10.1.1", "Master address. Defaults to 10.10.1.1.")
var masterPort *int = flag.Int("mport", 7087, "Master port.  Defaults to 7087.")
var myAddr *string = flag.String("addr", "10.10.1.1", "Server address (this machine). Defaults to 10.10.1.1.")
var portnum *int = flag.Int("port", 7070, "Port # to listen on. Defaults to 7070")
var doPineapple *bool = flag.Bool("pineapple", true, " Use Pineapple as the replication protocol. Defaults to true.")
var procs *int = flag.Int("p", 2, "GOMAXPROCS. Defaults to 2")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var exec = flag.Bool("exec", false, "Execute commands.")
var dreply = flag.Bool("dreply", true, "Reply to client only after command has been executed.")
var beacon = flag.Bool("beacon", false, "Send beacons to other replicas to compare their relative speeds.")
var durable = flag.Bool("durable", false, "Log to a stable store (i.e., a file in the current dir).")

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt)
		go catchKill(interrupt)
	}

	log.Printf("Server starting on port %d\n", *portnum)

	replicaId, nodeList := registerWithMaster(fmt.Sprintf("%s:%d", *masterAddr, *masterPort))

	if *doPineapple {
		log.Println("Starting Pineapple replica...")
		rep := pineapple.NewReplica(replicaId, nodeList, *exec, *dreply)
		rpc.Register(rep)
	}

	rpc.HandleHTTP()
	//listen for RPC on a different port (8070 by default)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *portnum+1000))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(l, nil)
}

func registerWithMaster(masterAddr string) (int, []string) {
	args := &masterproto.RegisterArgs{*myAddr, *portnum}
	var reply masterproto.RegisterReply

	for done := false; !done; {
		mcli, err := rpc.DialHTTP("tcp", masterAddr)
		if err == nil {
			err = mcli.Call("Master.Register", args, &reply)
			if err == nil && reply.Ready == true {
				done = true
				break
			}
		}
		time.Sleep(1e9)
	}

	return reply.ReplicaId, reply.NodeList
}

func catchKill(interrupt chan os.Signal) {
	signal := <-interrupt

	if signal == syscall.SIGURG {
		fmt.Println("Caught SIGURG, proceeding with run.")
		return
	}

	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	fmt.Println("Caught signal")
	os.Exit(0)
}
