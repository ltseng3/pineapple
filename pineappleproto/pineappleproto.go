package pineappleproto

import "pineapple/state"

const (
	PREPARE uint8 = iota
	PREPARE_REPLY
	ACCEPT
	ACCEPT_REPLY
)

type Tag struct {
	Timestamp int
	ID        int
}

type Payload struct {
	Tag   Tag
	Value int
}

type Get struct {
	ReplicaID int32
	Instance  int32
	Write     uint8
	Key       int
	Payload   Payload
}

type GetReply struct {
	Instance int32
	OK       uint8
	Write    uint8
	Key      int
	Payload  Payload
}

type Set struct {
	ReplicaID int32
	Instance  int32
	Write     uint8
	Key       int
	Payload   Payload
}

type SetReply struct {
	Instance int32
}

type Prepare struct {
	LeaderId   int32
	Instance   int32
	Ballot     int32
	ToInfinity uint8
}

type PrepareReply struct {
	Instance int32
	OK       uint8
	Ballot   int32
	Command  []state.Command
}

type RMWGet struct {
	LeaderId int32
	Instance int32
	Ballot   int32
	Command  []state.Command
}

type RMWGetReply struct {
	Instance int32
	Ballot   int32
	Key      int
	Payload  Payload
}

type RMWSet struct {
	LeaderId int32
	Instance int32
	Ballot   int32
	Command  []state.Command
	Key      int
	Payload  Payload
}

type RMWSetReply struct {
	Instance int32
	OK       uint8
	Ballot   int32
}

type Commit struct {
	LeaderId int32
	Instance int32
	Ballot   int32
	Command  []state.Command
}

type CommitShort struct {
	LeaderId int32
	Instance int32
	Count    int32
	Ballot   int32
}
