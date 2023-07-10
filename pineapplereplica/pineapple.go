package pineapplereplica

import (
	"encoding/binary"
	"io"
	"log"
	"sync"
	"time"

	"../dlog"
	"../fastrpc"
	"../genericsmr"
	"../genericsmrproto"
	"../pineapplereplicaproto"
	"../state"
)

const CLOCK = 1000 * 10
const CHAN_BUFFER_SIZE = 200000
const TRUE = uint8(1)
const FALSE = uint8(0)

type InstanceStatus int

const (
	PREPARING InstanceStatus = iota
	PREPARED
	ACCEPTED
	COMMITTED
)

// Node: performs ABD operations on single read write, and Paxos on multi read write and RMW
type Replica struct {
	*genericsmr.Replica // extends a generic Paxos replica

	// ABD
	mu           sync.Mutex
	getChan      chan fastrpc.Serializable
	setChan      chan fastrpc.Serializable
	getReplyChan chan fastrpc.Serializable
	setReplyChan chan fastrpc.Serializable
	getRPC       uint8
	setRPC       uint8
	getReplyRPC  uint8
	setReplyRPC  uint8

	// Paxos
	prepareChan      chan fastrpc.Serializable
	acceptChan       chan fastrpc.Serializable
	prepareReplyChan chan fastrpc.Serializable
	acceptReplyChan  chan fastrpc.Serializable
	commitChan       chan fastrpc.Serializable
	commitShortChan  chan fastrpc.Serializable
	prepareRPC       uint8
	acceptRPC        uint8
	prepareReplyRPC  uint8
	acceptReplyRPC   uint8
	commitRPC        uint8
	commitShortRPC   uint8

	IsLeader      bool // does this replica think it is the leader
	Shutdown      bool
	instanceSpace []*Instance // the space of all instances (used and not yet used)
	defaultBallot int32       // default ballot for new instances (0 until a Prepare(ballot, instance->infinity) from a leader)
	crtInstance   int32       // highest active instance number that this replica knows about

	flush         bool
	committedUpTo int32
}

type Instance struct {
	cmds   []state.Command
	ballot int32
	status InstanceStatus
	lb     *LeaderBookkeeping
	data   pineapplereplicaproto.Payload
}

type LeaderBookkeeping struct {
	clientProposals []*genericsmr.Propose
	maxRecvBallot   int32
	prepareOKs      int
	acceptOKs       int
	nacks           int
	completed       bool
}

func NewReplica(id int, peerAddrList []string, exec bool, dreply bool) *Replica {
	// extends a normal replica
	r := &Replica{
		genericsmr.NewReplica(id, peerAddrList, exec, dreply),
		sync.Mutex{},
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, 3*CHAN_BUFFER_SIZE),
		0,
		0,
		0,
		0,
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		make(chan fastrpc.Serializable, CHAN_BUFFER_SIZE),
		0,
		0,
		0,
		0,
		0,
		0,

		false,
		false,
		make([]*Instance, 20*1024*1024),
		0,
		0,

		false,
		0,
	}

	// ABD
	r.getRPC = r.RegisterRPC(new(pineapplereplicaproto.Get), r.getChan)
	r.setRPC = r.RegisterRPC(new(pineapplereplicaproto.Set), r.setChan)
	r.getReplyRPC = r.RegisterRPC(new(pineapplereplicaproto.GetReply), r.getReplyChan)
	r.setReplyRPC = r.RegisterRPC(new(pineapplereplicaproto.SetReply), r.setReplyChan)

	// Paxos
	r.prepareRPC = r.RegisterRPC(new(pineapplereplicaproto.Prepare), r.prepareChan)
	r.acceptRPC = r.RegisterRPC(new(pineapplereplicaproto.Accept), r.acceptChan)
	r.prepareReplyRPC = r.RegisterRPC(new(pineapplereplicaproto.PrepareReply), r.prepareReplyChan)
	r.acceptReplyRPC = r.RegisterRPC(new(pineapplereplicaproto.AcceptReply), r.acceptReplyChan)
	r.commitRPC = r.RegisterRPC(new(pineapplereplicaproto.Commit), r.commitChan)
	r.commitShortRPC = r.RegisterRPC(new(pineapplereplicaproto.CommitShort), r.commitShortChan)

	go r.Run()

	return r
}

func (r *Replica) replyPrepare(replicaId int32, reply *pineapplereplicaproto.PrepareReply) {
	r.SendMsg(replicaId, r.prepareReplyRPC, reply)
}

func (r *Replica) replyAccept(replicaId int32, reply *pineapplereplicaproto.AcceptReply) {
	r.SendMsg(replicaId, r.acceptReplyRPC, reply)
}

func (r *Replica) replyGet(replicaId int32, reply *pineapplereplicaproto.GetReply) {
	r.SendMsg(replicaId, r.getReplyRPC, reply)
}

func (r *Replica) replySet(replicaId int32, reply *pineapplereplicaproto.SetReply) {
	r.SendMsg(replicaId, r.setReplyRPC, reply)
}

// Get Phase
func (r *Replica) bcastGet(instance int32, write bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Prepare bcast failed:", err)
		}
	}()
	wr := FALSE
	if write {
		wr = TRUE
	}
	args := &pineapplereplicaproto.Get{r.Id, instance, wr}

	n := r.N - 1
	q := r.Id

	// Send to each connected replica
	for sent := 0; sent < n; {
		q = (q + 1) % int32(r.N)
		if q == r.Id {
			break
		}
		if !r.Alive[q] {
			continue
		}
		sent++
		r.SendMsg(q, r.getRPC, args)
	}
}

// ABD Get phase, returns vt pair currently held by node
func (r *Replica) handleGet(get *pineapplereplicaproto.Get) {
	var getReply *pineapplereplicaproto.GetReply
	var command state.Command
	ok := TRUE

	// If init or payload is empty, simply return empty payload
	if r.instanceSpace[r.crtInstance] == nil {
		getReply = &pineapplereplicaproto.GetReply{get.Instance, ok, get.Write, pineapplereplicaproto.Payload{}} // Empty payload
		r.replyGet(get.ReplicaID, getReply)
		return
	}

	// Return the most recent data held by storage node only if READ, since payload would be overwritten in write
	if get.Write == 1 {
		getReply = &pineapplereplicaproto.GetReply{get.Instance, ok, get.Write, r.instanceSpace[r.crtInstance].data}
		command.Op = 1
	} else {
		getReply = &pineapplereplicaproto.GetReply{get.Instance, ok, get.Write, pineapplereplicaproto.Payload{}} // Empty payload
	}

	/*
		cmds := make([]state.Command, 1)

			if getReply.OK == TRUE {
				r.recordCommands(cmds)
				r.sync()
			}
	*/

	r.replyGet(get.ReplicaID, getReply)
}

// ABD Set phase, sets current vt pair according to write message to all storage nodes
func (r *Replica) handleSet(set *pineapplereplicaproto.Set) {
	var setReply *pineapplereplicaproto.SetReply
	ok := TRUE

	// If init or payload is empty, simply return empty payload
	if r.instanceSpace[r.crtInstance] == nil {
		setReply = &pineapplereplicaproto.SetReply{set.Instance, ok, set.Write, pineapplereplicaproto.Payload{}} // Empty payload
		r.replySet(set.ReplicaID, setReply)
		return
	}

	if set.Payload.Tag.Timestamp > r.instanceSpace[r.crtInstance].data.Tag.Timestamp {
		r.instanceSpace[r.crtInstance].data = set.Payload
		// ok = TRUE
	}

	// OK FALSE indicates that no changes were made to the current data held.
	setReply = &pineapplereplicaproto.SetReply{set.Instance, ok, set.Write, r.instanceSpace[r.crtInstance].data}

	//r.sync()
	r.replySet(set.ReplicaID, setReply)
}

// Chooses the most recent vt pair after waiting for majority acknowlegements (or increment timestamp if write)
func (r *Replica) handleGetReply(getReply *pineapplereplicaproto.GetReply) {
	inst := r.instanceSpace[getReply.Instance]

	// Find the largest received timestamp
	if getReply.Payload.Tag.Timestamp > r.instanceSpace[r.crtInstance].data.Tag.Timestamp {
		r.instanceSpace[r.crtInstance].data = getReply.Payload
	}

	log.Print("Instance: ", getReply.Instance)
	log.Print("Currently held replica data: ", r.instanceSpace[r.crtInstance].data, " by replica ", r.Id)

	// Send the new vt pair to all nodes after getting majority
	if getReply.OK == TRUE {
		inst.lb.prepareOKs++

		if inst.lb.prepareOKs+1 > r.N>>1 {
			write := false
			inst.status = PREPARED
			inst.lb.nacks = 0
			// If write, choose a higher unique timestamp (by adjoining replica ID with Timestamp++)
			if getReply.Write == 1 {
				write = true
				r.mu.Lock()
				r.instanceSpace[r.crtInstance].data.Tag.Timestamp++
				r.instanceSpace[r.crtInstance].data.Tag.ID = int(r.Id)
				r.mu.Unlock()
				log.Print("Write received, current data: ", r.instanceSpace[r.crtInstance].data)
			}
			r.sync()
			r.bcastSet(getReply.Instance, getReply.Payload, write)
		}
	}
}

// Set Phase
func (r *Replica) bcastSet(instance int32, payload pineapplereplicaproto.Payload, write bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Prepare bcast failed:", err)
		}
	}()
	wr := FALSE
	if write {
		wr = TRUE
	}
	log.Println("Preparing args for Set Phase")
	args := &pineapplereplicaproto.Set{r.Id, instance, wr, payload}
	log.Print("Replica ID: ", args.ReplicaID, " Instance: ", args.Instance, " Write: ", args.Write)

	n := r.N - 1
	q := r.Id

	// Send to each connected replica
	for sent := 0; sent < n; {
		q = (q + 1) % int32(r.N)
		if q == r.Id {
			break
		}
		if !r.Alive[q] {
			continue
		}
		sent++
		r.SendMsg(q, r.setRPC, args)
	}
}

// Response handler for Set request on nodes
func (r *Replica) handleSetReply(setReply *pineapplereplicaproto.SetReply) {
	inst := r.instanceSpace[setReply.Instance]
	log.Print("Set reply status ", setReply.OK, " for instance ", setReply.Instance)

	if setReply.OK == TRUE {
		inst.lb.acceptOKs++

		// Wait for majority of acknowledgements
		if inst.lb.acceptOKs+1 > r.N>>1 {
			if inst.lb.clientProposals != nil && r.Dreply && !inst.lb.completed {
				// give client the all clear
				log.Print(inst.lb.clientProposals)
				propreply := &genericsmrproto.ProposeReplyTS{
					TRUE,
					inst.lb.clientProposals[0].CommandId,
					state.NIL,
					inst.lb.clientProposals[0].Timestamp}
				log.Print("Replying to client")
				r.ReplyProposeTS(propreply, inst.lb.clientProposals[0].Reply)
				inst.lb.completed = true
			}

			r.sync() //is this necessary?
		}
	}

	log.Print("Set reply received")
}

func (r *Replica) bcastPrepare(instance int32, ballot int32, toInfinity bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Prepare bcast failed:", err)
		}
	}()
	ti := FALSE
	if toInfinity {
		ti = TRUE
	}
	args := &pineapplereplicaproto.Prepare{r.Id, instance, ballot, ti}

	n := r.N - 1
	q := r.Id

	for sent := 0; sent < n; {
		q = (q + 1) % int32(r.N)
		if q == r.Id {
			break
		}
		if !r.Alive[q] {
			continue
		}
		sent++
		r.SendMsg(q, r.prepareRPC, args)
	}
}

var pa pineapplereplicaproto.Accept

func (r *Replica) bcastAccept(instance int32, ballot int32, command []state.Command) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Accept bcast failed:", err)
		}
	}()
	pa.LeaderId = r.Id
	pa.Instance = instance
	pa.Ballot = ballot
	pa.Command = command
	args := &pa
	//args := &paxosproto.Accept{r.Id, instance, ballot, command}

	n := r.N - 1
	q := r.Id

	for sent := 0; sent < n; {
		q = (q + 1) % int32(r.N)
		if q == r.Id {
			break
		}
		if !r.Alive[q] {
			continue
		}
		sent++
		r.SendMsg(q, r.acceptRPC, args)
	}
}

var pc pineapplereplicaproto.Commit
var pcs pineapplereplicaproto.CommitShort

func (r *Replica) bcastCommit(instance int32, ballot int32, command []state.Command) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Commit bcast failed:", err)
		}
	}()
	pc.LeaderId = r.Id
	pc.Instance = instance
	pc.Ballot = ballot
	pc.Command = command
	args := &pc
	pcs.LeaderId = r.Id
	pcs.Instance = instance
	pcs.Ballot = ballot
	pcs.Count = int32(len(command))
	argsShort := &pcs

	//args := &paxosproto.Commit{r.Id, instance, command}

	n := r.N - 1
	q := r.Id
	sent := 0

	for sent < n {
		q = (q + 1) % int32(r.N)
		if q == r.Id {
			break
		}
		if !r.Alive[q] {
			continue
		}
		sent++
		r.SendMsg(q, r.commitShortRPC, argsShort)
	}
	if q != r.Id {
		for sent < r.N-1 {
			q = (q + 1) % int32(r.N)
			if q == r.Id {
				break
			}
			if !r.Alive[q] {
				continue
			}
			sent++
			r.SendMsg(q, r.commitRPC, args)
		}
	}
}

func (r *Replica) handlePropose(propose *genericsmr.Propose) {
	/*
		if !r.IsLeader {
			preply := &genericsmrproto.ProposeReplyTS{TRUE, -1, state.NIL, 0}
			r.ReplyProposeTS(preply, propose.Reply)
			return
		}
	*/

	for r.instanceSpace[r.crtInstance] != nil {
		r.crtInstance++
	}

	instNo := r.crtInstance

	cmds := make([]state.Command, 1)
	proposals := make([]*genericsmr.Propose, 1)
	cmds[0] = propose.Command
	proposals[0] = propose

	// ABD
	r.instanceSpace[instNo] = &Instance{
		cmds,
		r.makeUniqueBallot(0),
		PREPARING,
		&LeaderBookkeeping{proposals, 0, 0, 0, 0, false},
		pineapplereplicaproto.Payload{Tag: pineapplereplicaproto.Tag{Timestamp: int(propose.Timestamp), ID: int(r.Id)}, Value: int(propose.Command.V)},
	}

	r.recordInstanceMetadata(r.instanceSpace[instNo])
	r.recordCommands(cmds)
	r.sync()

	// Construct the pineapple payload from proposal data
	r.bcastGet(instNo, false) // @audit autodetermine proposal type here
	log.Printf("ABD for instance %d\n", instNo)

	// Use Paxos if operation is not Read / Write
	if propose.Command.Op != state.PUT || propose.Command.Op != state.GET {
		if r.defaultBallot == -1 {
			r.instanceSpace[instNo] = &Instance{
				cmds,
				r.makeUniqueBallot(0),
				PREPARING,
				&LeaderBookkeeping{proposals, 0, 0, 0, 0, false},
				pineapplereplicaproto.Payload{},
			}
			r.bcastPrepare(instNo, r.makeUniqueBallot(0), true)
			dlog.Printf("Classic round for instance %d\n", instNo)
		} else {
			r.instanceSpace[instNo] = &Instance{
				cmds,
				r.defaultBallot,
				PREPARED,
				&LeaderBookkeeping{proposals, 0, 0, 0, 0, false},
				pineapplereplicaproto.Payload{},
			}

			r.recordInstanceMetadata(r.instanceSpace[instNo])
			r.recordCommands(cmds)
			r.sync()

			r.bcastAccept(instNo, r.defaultBallot, cmds)
			dlog.Printf("Fast round for instance %d\n", instNo)
		}
	}
}

func (r *Replica) handlePrepare(prepare *pineapplereplicaproto.Prepare) {
	inst := r.instanceSpace[prepare.Instance]
	var preply *pineapplereplicaproto.PrepareReply

	if inst == nil {
		ok := TRUE
		if r.defaultBallot > prepare.Ballot {
			ok = FALSE
		}
		preply = &pineapplereplicaproto.PrepareReply{prepare.Instance, ok, r.defaultBallot, make([]state.Command, 0)}
	} else {
		ok := TRUE
		if prepare.Ballot < inst.ballot {
			ok = FALSE
		}
		preply = &pineapplereplicaproto.PrepareReply{prepare.Instance, ok, inst.ballot, inst.cmds}
	}

	r.replyPrepare(prepare.LeaderId, preply)

	if prepare.ToInfinity == TRUE && prepare.Ballot > r.defaultBallot {
		r.defaultBallot = prepare.Ballot
	}
}

func (r *Replica) handleAccept(accept *pineapplereplicaproto.Accept) {
	inst := r.instanceSpace[accept.Instance]
	var areply *pineapplereplicaproto.AcceptReply

	if inst == nil {
		if accept.Ballot < r.defaultBallot {
			areply = &pineapplereplicaproto.AcceptReply{accept.Instance, FALSE, r.defaultBallot}
		} else {
			r.instanceSpace[accept.Instance] = &Instance{
				accept.Command,
				accept.Ballot,
				ACCEPTED,
				nil,
				pineapplereplicaproto.Payload{},
			}
			areply = &pineapplereplicaproto.AcceptReply{accept.Instance, TRUE, r.defaultBallot}
		}
	} else if inst.ballot > accept.Ballot {
		areply = &pineapplereplicaproto.AcceptReply{accept.Instance, FALSE, inst.ballot}
	} else if inst.ballot < accept.Ballot {
		inst.cmds = accept.Command
		inst.ballot = accept.Ballot
		inst.status = ACCEPTED
		areply = &pineapplereplicaproto.AcceptReply{accept.Instance, TRUE, inst.ballot}
		if inst.lb != nil && inst.lb.clientProposals != nil {
			//TODO: is this correct?
			// try the proposal in a different instance
			for i := 0; i < len(inst.lb.clientProposals); i++ {
				r.ProposeChan <- inst.lb.clientProposals[i]
			}
			inst.lb.clientProposals = nil
		}
	} else {
		// reordered ACCEPT
		r.instanceSpace[accept.Instance].cmds = accept.Command
		if r.instanceSpace[accept.Instance].status != COMMITTED {
			r.instanceSpace[accept.Instance].status = ACCEPTED
		}
		areply = &pineapplereplicaproto.AcceptReply{accept.Instance, TRUE, r.defaultBallot}
	}

	if areply.OK == TRUE {
		r.recordInstanceMetadata(r.instanceSpace[accept.Instance])
		r.recordCommands(accept.Command)
		r.sync()
	}

	r.replyAccept(accept.LeaderId, areply)
}

func (r *Replica) handleCommit(commit *pineapplereplicaproto.Commit) {
	inst := r.instanceSpace[commit.Instance]

	dlog.Printf("Committing instance %d\n", commit.Instance)

	if inst == nil {
		r.instanceSpace[commit.Instance] = &Instance{
			commit.Command,
			commit.Ballot,
			COMMITTED,
			nil,
			pineapplereplicaproto.Payload{},
		}
	} else {
		r.instanceSpace[commit.Instance].cmds = commit.Command
		r.instanceSpace[commit.Instance].status = COMMITTED
		r.instanceSpace[commit.Instance].ballot = commit.Ballot
		if inst.lb != nil && inst.lb.clientProposals != nil {
			for i := 0; i < len(inst.lb.clientProposals); i++ {
				r.ProposeChan <- inst.lb.clientProposals[i]
			}
			inst.lb.clientProposals = nil
		}
	}

	r.updateCommittedUpTo()

	r.recordInstanceMetadata(r.instanceSpace[commit.Instance])
	r.recordCommands(commit.Command)
}

func (r *Replica) handleCommitShort(commit *pineapplereplicaproto.CommitShort) {
	inst := r.instanceSpace[commit.Instance]

	dlog.Printf("Committing instance %d\n", commit.Instance)

	if inst == nil {
		r.instanceSpace[commit.Instance] = &Instance{nil,
			commit.Ballot,
			COMMITTED,
			nil,
			pineapplereplicaproto.Payload{},
		}
	} else {
		r.instanceSpace[commit.Instance].status = COMMITTED
		r.instanceSpace[commit.Instance].ballot = commit.Ballot
		if inst.lb != nil && inst.lb.clientProposals != nil {
			for i := 0; i < len(inst.lb.clientProposals); i++ {
				r.ProposeChan <- inst.lb.clientProposals[i]
			}
			inst.lb.clientProposals = nil
		}
	}

	r.updateCommittedUpTo()

	r.recordInstanceMetadata(r.instanceSpace[commit.Instance])
}

func (r *Replica) handlePrepareReply(preply *pineapplereplicaproto.PrepareReply) {
	inst := r.instanceSpace[preply.Instance]

	if inst.status != PREPARING {
		// TODO: should replies for non-current ballots be ignored?
		// we've moved on -- these are delayed replies, so just ignore
		return
	}

	if preply.OK == TRUE {
		inst.lb.prepareOKs++

		if preply.Ballot > inst.lb.maxRecvBallot {
			inst.cmds = preply.Command
			inst.lb.maxRecvBallot = preply.Ballot
			if inst.lb.clientProposals != nil {
				// there is already a competing command for this instance,
				// so we put the client proposal back in the queue so that
				// we know to try it in another instance
				for i := 0; i < len(inst.lb.clientProposals); i++ {
					r.ProposeChan <- inst.lb.clientProposals[i]
				}
				inst.lb.clientProposals = nil
			}
		}

		if inst.lb.prepareOKs+1 > r.N>>1 {
			inst.status = PREPARED
			inst.lb.nacks = 0
			if inst.ballot > r.defaultBallot {
				r.defaultBallot = inst.ballot
			}
			r.recordInstanceMetadata(r.instanceSpace[preply.Instance])
			r.sync()
			r.bcastAccept(preply.Instance, inst.ballot, inst.cmds)
		}
	} else {
		// TODO: there is probably another active leader
		inst.lb.nacks++
		if preply.Ballot > inst.lb.maxRecvBallot {
			inst.lb.maxRecvBallot = preply.Ballot
		}
		if inst.lb.nacks >= r.N>>1 {
			if inst.lb.clientProposals != nil {
				// try the proposals in another instance
				for i := 0; i < len(inst.lb.clientProposals); i++ {
					r.ProposeChan <- inst.lb.clientProposals[i]
				}
				inst.lb.clientProposals = nil
			}
		}
	}
}

func (r *Replica) handleAcceptReply(areply *pineapplereplicaproto.AcceptReply) {
	inst := r.instanceSpace[areply.Instance]

	if inst.status != PREPARED && inst.status != ACCEPTED {
		// we've move on, these are delayed replies, so just ignore
		return
	}

	if areply.OK == TRUE {
		inst.lb.acceptOKs++
		if inst.lb.acceptOKs+1 > r.N>>1 {
			inst = r.instanceSpace[areply.Instance]
			inst.status = COMMITTED
			if inst.lb.clientProposals != nil && !r.Dreply {
				// give client the all clear
				for i := 0; i < len(inst.cmds); i++ {
					propreply := &genericsmrproto.ProposeReplyTS{
						TRUE,
						inst.lb.clientProposals[i].CommandId,
						state.NIL,
						inst.lb.clientProposals[i].Timestamp}
					r.ReplyProposeTS(propreply, inst.lb.clientProposals[i].Reply)
				}
			}

			r.recordInstanceMetadata(r.instanceSpace[areply.Instance])
			r.sync() //is this necessary?

			r.updateCommittedUpTo()

			r.bcastCommit(areply.Instance, inst.ballot, inst.cmds)
		}
	} else {
		// TODO: there is probably another active leader
		inst.lb.nacks++
		if areply.Ballot > inst.lb.maxRecvBallot {
			inst.lb.maxRecvBallot = areply.Ballot
		}
		if inst.lb.nacks >= r.N>>1 {
			// TODO
		}
	}
}

func (r *Replica) executeCommands() {
	i := int32(0)
	for !r.Shutdown {
		executed := false

		for i <= r.committedUpTo {
			if r.instanceSpace[i].cmds != nil {
				inst := r.instanceSpace[i]
				for j := 0; j < len(inst.cmds); j++ {
					val := inst.cmds[j].Execute(r.State)
					if r.Dreply && inst.lb != nil && inst.lb.clientProposals != nil {
						propreply := &genericsmrproto.ProposeReplyTS{
							TRUE,
							inst.lb.clientProposals[j].CommandId,
							val,
							inst.lb.clientProposals[j].Timestamp}
						r.ReplyProposeTS(propreply, inst.lb.clientProposals[j].Reply)
					}
				}
				i++
				executed = true
			} else {
				break
			}
		}

		if !executed {
			time.Sleep(CLOCK)
		}
	}

}

var clockChan chan bool

func (r *Replica) makeUniqueBallot(ballot int32) int32 {
	return (ballot << 4) | r.Id
}

func (r *Replica) updateCommittedUpTo() {
	for r.instanceSpace[r.committedUpTo+1] != nil &&
		r.instanceSpace[r.committedUpTo+1].status == COMMITTED {
		r.committedUpTo++
	}
}

// append a log entry to stable storage
func (r *Replica) recordInstanceMetadata(inst *Instance) {
	if !r.Durable {
		return
	}

	var b [5]byte
	binary.LittleEndian.PutUint32(b[0:4], uint32(inst.ballot))
	b[4] = byte(inst.status)
	r.StableStore.Write(b[:])
}

// write a sequence of commands to stable storage
func (r *Replica) recordCommands(cmds []state.Command) {
	if !r.Durable {
		return
	}

	if cmds == nil {
		return
	}
	for i := 0; i < len(cmds); i++ {
		cmds[i].Marshal(io.Writer(r.StableStore))
	}
}

// sync with the stable store
func (r *Replica) sync() {
	if !r.Durable {
		return
	}

	r.StableStore.Sync()
}

func (r *Replica) clock() {
	for !r.Shutdown {
		time.Sleep(CLOCK)
		clockChan <- true
	}
}

// Main processing loop
func (r *Replica) Run() {
	r.ConnectToPeers()

	log.Println("Waiting for client connections")

	go r.WaitForClientConnections()

	clockChan = make(chan bool, 1)
	go r.clock()

	// We don't directly access r.ProposeChan, because we want to do pipelining periodically,
	// so we introduces a channel pointer: onOffProposChan:
	onOffProposeChan := r.ProposeChan

	for !r.Shutdown {

		select {

		case <-clockChan:
			// activate the new proposals channel
			onOffProposeChan = r.ProposeChan
			break
		case setS := <-r.setChan:
			set := setS.(*pineapplereplicaproto.Set)
			//got a Write message
			log.Printf("Received Write from replica %d, for instance %d\n", set.Payload.Tag.ID, set.Instance)
			r.handleSet(set)
			break
		case getS := <-r.getChan:
			get := getS.(*pineapplereplicaproto.Get)
			//got a Read message
			log.Printf("Received Read from replica %d, for instance %d\n", get.ReplicaID, get.Instance)
			r.handleGet(get)
			break
		case setReplyS := <-r.setReplyChan:
			setReply := setReplyS.(*pineapplereplicaproto.SetReply)
			//got a Write reply
			log.Printf("Received WriteReply for instance %d\n", setReply.Instance)
			r.handleSetReply(setReply)
			break
		case getReplyS := <-r.getReplyChan:
			getReply := getReplyS.(*pineapplereplicaproto.GetReply)
			//got a Read reply
			log.Printf("Received ReadReply for instance %d\n", getReply.Instance)
			r.handleGetReply(getReply)
			break
		case propose := <-onOffProposeChan:
			//got a Propose from a client
			log.Printf("Proposal with op %d\n", propose.Command.Op)
			// Handle proposal: single read-write object goes to ABD, multi read/write or RMW goes to Paxos
			r.handlePropose(propose)
			// deactivate the new proposals channel to prioritize the handling of protocol messages
			onOffProposeChan = nil
			break
		case prepareS := <-r.prepareChan:
			prepare := prepareS.(*pineapplereplicaproto.Prepare)
			//got a Prepare message
			dlog.Printf("Received Prepare from replica %d, for instance %d\n", prepare.LeaderId, prepare.Instance)
			r.handlePrepare(prepare)
			break
		case acceptS := <-r.acceptChan:
			accept := acceptS.(*pineapplereplicaproto.Accept)
			//got an Accept message
			dlog.Printf("Received Accept from replica %d, for instance %d\n", accept.LeaderId, accept.Instance)
			r.handleAccept(accept)
			break
		case prepareReplyS := <-r.prepareReplyChan:
			prepareReply := prepareReplyS.(*pineapplereplicaproto.PrepareReply)
			//got a Prepare reply
			dlog.Printf("Received PrepareReply for instance %d\n", prepareReply.Instance)
			r.handlePrepareReply(prepareReply)
			break
		case acceptReplyS := <-r.acceptReplyChan:
			acceptReply := acceptReplyS.(*pineapplereplicaproto.AcceptReply)
			//got an Accept reply
			dlog.Printf("Received AcceptReply for instance %d\n", acceptReply.Instance)
			r.handleAcceptReply(acceptReply)
			break
		}
	}
}

/* RPC to be called by master */

func (r *Replica) BeTheLeader(args *genericsmrproto.BeTheLeaderArgs, reply *genericsmrproto.BeTheLeaderReply) error {
	r.IsLeader = true
	return nil
}
