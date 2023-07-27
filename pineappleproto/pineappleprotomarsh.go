package pineappleproto

import (
	"bufio"
	"encoding/binary"
	"io"
	"pineapple/fastrpc"
	"pineapple/state"
	"sync"
)

type byteReader interface {
	io.Reader
	ReadByte() (c byte, err error)
}

func (t *GetReply) New() fastrpc.Serializable {
	return new(GetReply)
}
func (t *GetReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 38, true
}

type GetReplyCache struct {
	mu    sync.Mutex
	cache []*GetReply
}

func NewGetReplyCache() *GetReplyCache {
	c := &GetReplyCache{}
	c.cache = make([]*GetReply, 0)
	return c
}

func (p *GetReplyCache) Get() *GetReply {
	var t *GetReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &GetReply{}
	}
	return t
}
func (p *GetReplyCache) Put(t *GetReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *GetReply) Marshal(wire io.Writer) {
	var b [38]byte
	var bs []byte
	bs = b[:38]
	tmp32 := t.Instance
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	bs[4] = byte(t.OK)
	bs[5] = byte(t.Write)
	tmp64 := t.Key
	bs[6] = byte(tmp64 >> 56)
	bs[7] = byte(tmp64 >> 48)
	bs[8] = byte(tmp64 >> 40)
	bs[9] = byte(tmp64 >> 32)
	bs[10] = byte(tmp64 >> 24)
	bs[11] = byte(tmp64 >> 16)
	bs[12] = byte(tmp64 >> 8)
	bs[13] = byte(tmp64)
	tmp64 = t.Payload.Tag.Timestamp
	bs[14] = byte(tmp64 >> 56)
	bs[15] = byte(tmp64 >> 48)
	bs[16] = byte(tmp64 >> 40)
	bs[17] = byte(tmp64 >> 32)
	bs[18] = byte(tmp64 >> 24)
	bs[19] = byte(tmp64 >> 16)
	bs[20] = byte(tmp64 >> 8)
	bs[21] = byte(tmp64)
	tmp64 = t.Payload.Tag.ID
	bs[22] = byte(tmp64 >> 56)
	bs[23] = byte(tmp64 >> 48)
	bs[24] = byte(tmp64 >> 40)
	bs[25] = byte(tmp64 >> 32)
	bs[26] = byte(tmp64 >> 24)
	bs[27] = byte(tmp64 >> 16)
	bs[28] = byte(tmp64 >> 8)
	bs[29] = byte(tmp64)
	tmp64 = t.Payload.Value
	bs[30] = byte(tmp64 >> 56)
	bs[31] = byte(tmp64 >> 48)
	bs[32] = byte(tmp64 >> 40)
	bs[33] = byte(tmp64 >> 32)
	bs[34] = byte(tmp64 >> 24)
	bs[35] = byte(tmp64 >> 16)
	bs[36] = byte(tmp64 >> 8)
	bs[37] = byte(tmp64)
	wire.Write(bs)
}

func (t *GetReply) Unmarshal(wire io.Reader) error {
	var b [38]byte
	var bs []byte
	bs = b[:38]
	if _, err := io.ReadAtLeast(wire, bs, 38); err != nil {
		return err
	}
	t.Instance = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.OK = uint8(bs[4])
	t.Write = uint8(bs[5])
	t.Key = int(((uint64(bs[6]) << 56) | (uint64(bs[7]) << 48) | (uint64(bs[8]) << 40) | (uint64(bs[9]) << 32) | (uint64(bs[10]) << 24) | (uint64(bs[11]) << 16) | (uint64(bs[12]) << 8) | uint64(bs[13])))
	t.Payload.Tag.Timestamp = int(((uint64(bs[14]) << 56) | (uint64(bs[15]) << 48) | (uint64(bs[16]) << 40) | (uint64(bs[17]) << 32) | (uint64(bs[18]) << 24) | (uint64(bs[19]) << 16) | (uint64(bs[20]) << 8) | uint64(bs[21])))
	t.Payload.Tag.ID = int(((uint64(bs[22]) << 56) | (uint64(bs[23]) << 48) | (uint64(bs[24]) << 40) | (uint64(bs[25]) << 32) | (uint64(bs[26]) << 24) | (uint64(bs[27]) << 16) | (uint64(bs[28]) << 8) | uint64(bs[29])))
	t.Payload.Value = int(((uint64(bs[30]) << 56) | (uint64(bs[31]) << 48) | (uint64(bs[32]) << 40) | (uint64(bs[33]) << 32) | (uint64(bs[34]) << 24) | (uint64(bs[35]) << 16) | (uint64(bs[36]) << 8) | uint64(bs[37])))
	return nil
}

func (t *Commit) New() fastrpc.Serializable {
	return new(Commit)
}
func (t *Commit) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type CommitCache struct {
	mu    sync.Mutex
	cache []*Commit
}

func NewCommitCache() *CommitCache {
	c := &CommitCache{}
	c.cache = make([]*Commit, 0)
	return c
}

func (p *CommitCache) Get() *Commit {
	var t *Commit
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Commit{}
	}
	return t
}
func (p *CommitCache) Put(t *Commit) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Commit) Marshal(wire io.Writer) {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	tmp32 := t.LeaderId
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp32 = t.Ballot
	bs[8] = byte(tmp32 >> 24)
	bs[9] = byte(tmp32 >> 16)
	bs[10] = byte(tmp32 >> 8)
	bs[11] = byte(tmp32)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Command))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Marshal(wire)
	}
}

func (t *Commit) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [12]byte
	var bs []byte
	bs = b[:12]
	if _, err := io.ReadAtLeast(wire, bs, 12); err != nil {
		return err
	}
	t.LeaderId = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Ballot = int32(((uint32(bs[8]) << 24) | (uint32(bs[9]) << 16) | (uint32(bs[10]) << 8) | uint32(bs[11])))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Command = make([]state.Command, alen1)
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Unmarshal(wire)
	}
	return nil
}

func (t *RMWSetReply) New() fastrpc.Serializable {
	return new(RMWSetReply)
}
func (t *RMWSetReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 9, true
}

type RMWSetReplyCache struct {
	mu    sync.Mutex
	cache []*RMWSetReply
}

func NewRMWSetReplyCache() *RMWSetReplyCache {
	c := &RMWSetReplyCache{}
	c.cache = make([]*RMWSetReply, 0)
	return c
}

func (p *RMWSetReplyCache) Get() *RMWSetReply {
	var t *RMWSetReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &RMWSetReply{}
	}
	return t
}
func (p *RMWSetReplyCache) Put(t *RMWSetReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *RMWSetReply) Marshal(wire io.Writer) {
	var b [9]byte
	var bs []byte
	bs = b[:9]
	tmp32 := t.Instance
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	bs[4] = byte(t.OK)
	tmp32 = t.Ballot
	bs[5] = byte(tmp32 >> 24)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 8)
	bs[8] = byte(tmp32)
	wire.Write(bs)
}

func (t *RMWSetReply) Unmarshal(wire io.Reader) error {
	var b [9]byte
	var bs []byte
	bs = b[:9]
	if _, err := io.ReadAtLeast(wire, bs, 9); err != nil {
		return err
	}
	t.Instance = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.OK = uint8(bs[4])
	t.Ballot = int32(((uint32(bs[5]) << 24) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 8) | uint32(bs[8])))
	return nil
}

func (t *Get) New() fastrpc.Serializable {
	return new(Get)
}
func (t *Get) BinarySize() (nbytes int, sizeKnown bool) {
	return 17, true
}

type GetCache struct {
	mu    sync.Mutex
	cache []*Get
}

func NewGetCache() *GetCache {
	c := &GetCache{}
	c.cache = make([]*Get, 0)
	return c
}

func (p *GetCache) Get() *Get {
	var t *Get
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Get{}
	}
	return t
}
func (p *GetCache) Put(t *Get) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Get) Marshal(wire io.Writer) {
	var b [17]byte
	var bs []byte
	bs = b[:17]
	tmp32 := t.ReplicaID
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	bs[8] = byte(t.Write)
	tmp64 := t.Key
	bs[9] = byte(tmp64 >> 56)
	bs[10] = byte(tmp64 >> 48)
	bs[11] = byte(tmp64 >> 40)
	bs[12] = byte(tmp64 >> 32)
	bs[13] = byte(tmp64 >> 24)
	bs[14] = byte(tmp64 >> 16)
	bs[15] = byte(tmp64 >> 8)
	bs[16] = byte(tmp64)
	wire.Write(bs)
}

func (t *Get) Unmarshal(wire io.Reader) error {
	var b [17]byte
	var bs []byte
	bs = b[:17]
	if _, err := io.ReadAtLeast(wire, bs, 17); err != nil {
		return err
	}
	t.ReplicaID = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Write = uint8(bs[8])
	t.Key = int(((uint64(bs[9]) << 56) | (uint64(bs[10]) << 48) | (uint64(bs[11]) << 40) | (uint64(bs[12]) << 32) | (uint64(bs[13]) << 24) | (uint64(bs[14]) << 16) | (uint64(bs[15]) << 8) | uint64(bs[16])))
	return nil
}

func (t *SetReply) New() fastrpc.Serializable {
	return new(SetReply)
}
func (t *SetReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 4, true
}

type SetReplyCache struct {
	mu    sync.Mutex
	cache []*SetReply
}

func NewSetReplyCache() *SetReplyCache {
	c := &SetReplyCache{}
	c.cache = make([]*SetReply, 0)
	return c
}

func (p *SetReplyCache) Get() *SetReply {
	var t *SetReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &SetReply{}
	}
	return t
}
func (p *SetReplyCache) Put(t *SetReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *SetReply) Marshal(wire io.Writer) {
	var b [4]byte
	var bs []byte
	bs = b[:4]
	tmp32 := t.Instance
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	wire.Write(bs)
}

func (t *SetReply) Unmarshal(wire io.Reader) error {
	var b [4]byte
	var bs []byte
	bs = b[:4]
	if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
		return err
	}
	t.Instance = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	return nil
}

func (t *RMWGetReply) New() fastrpc.Serializable {
	return new(RMWGetReply)
}
func (t *RMWGetReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 40, true
}

type RMWGetReplyCache struct {
	mu    sync.Mutex
	cache []*RMWGetReply
}

func NewRMWGetReplyCache() *RMWGetReplyCache {
	c := &RMWGetReplyCache{}
	c.cache = make([]*RMWGetReply, 0)
	return c
}

func (p *RMWGetReplyCache) Get() *RMWGetReply {
	var t *RMWGetReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &RMWGetReply{}
	}
	return t
}
func (p *RMWGetReplyCache) Put(t *RMWGetReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *RMWGetReply) Marshal(wire io.Writer) {
	var b [40]byte
	var bs []byte
	bs = b[:40]
	tmp32 := t.Instance
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp64 := t.Key
	bs[8] = byte(tmp64 >> 56)
	bs[9] = byte(tmp64 >> 48)
	bs[10] = byte(tmp64 >> 40)
	bs[11] = byte(tmp64 >> 32)
	bs[12] = byte(tmp64 >> 24)
	bs[13] = byte(tmp64 >> 16)
	bs[14] = byte(tmp64 >> 8)
	bs[15] = byte(tmp64)
	tmp64 = t.Payload.Tag.Timestamp
	bs[16] = byte(tmp64 >> 56)
	bs[17] = byte(tmp64 >> 48)
	bs[18] = byte(tmp64 >> 40)
	bs[19] = byte(tmp64 >> 32)
	bs[20] = byte(tmp64 >> 24)
	bs[21] = byte(tmp64 >> 16)
	bs[22] = byte(tmp64 >> 8)
	bs[23] = byte(tmp64)
	tmp64 = t.Payload.Tag.ID
	bs[24] = byte(tmp64 >> 56)
	bs[25] = byte(tmp64 >> 48)
	bs[26] = byte(tmp64 >> 40)
	bs[27] = byte(tmp64 >> 32)
	bs[28] = byte(tmp64 >> 24)
	bs[29] = byte(tmp64 >> 16)
	bs[30] = byte(tmp64 >> 8)
	bs[31] = byte(tmp64)
	tmp64 = t.Payload.Value
	bs[32] = byte(tmp64 >> 56)
	bs[33] = byte(tmp64 >> 48)
	bs[34] = byte(tmp64 >> 40)
	bs[35] = byte(tmp64 >> 32)
	bs[36] = byte(tmp64 >> 24)
	bs[37] = byte(tmp64 >> 16)
	bs[38] = byte(tmp64 >> 8)
	bs[39] = byte(tmp64)
	wire.Write(bs)
}

func (t *RMWGetReply) Unmarshal(wire io.Reader) error {
	var b [40]byte
	var bs []byte
	bs = b[:40]
	if _, err := io.ReadAtLeast(wire, bs, 40); err != nil {
		return err
	}
	t.Instance = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Ballot = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Key = int(((uint64(bs[8]) << 56) | (uint64(bs[9]) << 48) | (uint64(bs[10]) << 40) | (uint64(bs[11]) << 32) | (uint64(bs[12]) << 24) | (uint64(bs[13]) << 16) | (uint64(bs[14]) << 8) | uint64(bs[15])))
	t.Payload.Tag.Timestamp = int(((uint64(bs[16]) << 56) | (uint64(bs[17]) << 48) | (uint64(bs[18]) << 40) | (uint64(bs[19]) << 32) | (uint64(bs[20]) << 24) | (uint64(bs[21]) << 16) | (uint64(bs[22]) << 8) | uint64(bs[23])))
	t.Payload.Tag.ID = int(((uint64(bs[24]) << 56) | (uint64(bs[25]) << 48) | (uint64(bs[26]) << 40) | (uint64(bs[27]) << 32) | (uint64(bs[28]) << 24) | (uint64(bs[29]) << 16) | (uint64(bs[30]) << 8) | uint64(bs[31])))
	t.Payload.Value = int(((uint64(bs[32]) << 56) | (uint64(bs[33]) << 48) | (uint64(bs[34]) << 40) | (uint64(bs[35]) << 32) | (uint64(bs[36]) << 24) | (uint64(bs[37]) << 16) | (uint64(bs[38]) << 8) | uint64(bs[39])))
	return nil
}

func (t *RMWSet) New() fastrpc.Serializable {
	return new(RMWSet)
}
func (t *RMWSet) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type RMWSetCache struct {
	mu    sync.Mutex
	cache []*RMWSet
}

func NewRMWSetCache() *RMWSetCache {
	c := &RMWSetCache{}
	c.cache = make([]*RMWSet, 0)
	return c
}

func (p *RMWSetCache) Get() *RMWSet {
	var t *RMWSet
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &RMWSet{}
	}
	return t
}
func (p *RMWSetCache) Put(t *RMWSet) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *RMWSet) Marshal(wire io.Writer) {
	var b [24]byte
	var bs []byte
	bs = b[:12]
	tmp32 := t.LeaderId
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp32 = t.Ballot
	bs[8] = byte(tmp32 >> 24)
	bs[9] = byte(tmp32 >> 16)
	bs[10] = byte(tmp32 >> 8)
	bs[11] = byte(tmp32)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Command))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Marshal(wire)
	}
	tmp64 := t.Payload.Tag.Timestamp
	bs[0] = byte(tmp64 >> 56)
	bs[1] = byte(tmp64 >> 48)
	bs[2] = byte(tmp64 >> 40)
	bs[3] = byte(tmp64 >> 32)
	bs[4] = byte(tmp64 >> 24)
	bs[5] = byte(tmp64 >> 16)
	bs[6] = byte(tmp64 >> 8)
	bs[7] = byte(tmp64)
	tmp64 = t.Payload.Tag.ID
	bs[8] = byte(tmp64 >> 56)
	bs[9] = byte(tmp64 >> 48)
	bs[10] = byte(tmp64 >> 40)
	bs[11] = byte(tmp64 >> 32)
	bs[12] = byte(tmp64 >> 24)
	bs[13] = byte(tmp64 >> 16)
	bs[14] = byte(tmp64 >> 8)
	bs[15] = byte(tmp64)
	tmp64 = t.Payload.Value
	bs[16] = byte(tmp64 >> 56)
	bs[17] = byte(tmp64 >> 48)
	bs[18] = byte(tmp64 >> 40)
	bs[19] = byte(tmp64 >> 32)
	bs[20] = byte(tmp64 >> 24)
	bs[21] = byte(tmp64 >> 16)
	bs[22] = byte(tmp64 >> 8)
	bs[23] = byte(tmp64)
	wire.Write(bs)
}

func (t *RMWSet) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [24]byte
	var bs []byte
	bs = b[:12]
	if _, err := io.ReadAtLeast(wire, bs, 12); err != nil {
		return err
	}
	t.LeaderId = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Ballot = int32(((uint32(bs[8]) << 24) | (uint32(bs[9]) << 16) | (uint32(bs[10]) << 8) | uint32(bs[11])))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Command = make([]state.Command, alen1)
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Unmarshal(wire)
	}
	bs = b[:24]
	if _, err := io.ReadAtLeast(wire, bs, 24); err != nil {
		return err
	}
	t.Payload.Tag.Timestamp = int(((uint64(bs[0]) << 56) | (uint64(bs[1]) << 48) | (uint64(bs[2]) << 40) | (uint64(bs[3]) << 32) | (uint64(bs[4]) << 24) | (uint64(bs[5]) << 16) | (uint64(bs[6]) << 8) | uint64(bs[7])))
	t.Payload.Tag.ID = int(((uint64(bs[8]) << 56) | (uint64(bs[9]) << 48) | (uint64(bs[10]) << 40) | (uint64(bs[11]) << 32) | (uint64(bs[12]) << 24) | (uint64(bs[13]) << 16) | (uint64(bs[14]) << 8) | uint64(bs[15])))
	t.Payload.Value = int(((uint64(bs[16]) << 56) | (uint64(bs[17]) << 48) | (uint64(bs[18]) << 40) | (uint64(bs[19]) << 32) | (uint64(bs[20]) << 24) | (uint64(bs[21]) << 16) | (uint64(bs[22]) << 8) | uint64(bs[23])))
	return nil
}

func (t *CommitShort) New() fastrpc.Serializable {
	return new(CommitShort)
}
func (t *CommitShort) BinarySize() (nbytes int, sizeKnown bool) {
	return 16, true
}

type CommitShortCache struct {
	mu    sync.Mutex
	cache []*CommitShort
}

func NewCommitShortCache() *CommitShortCache {
	c := &CommitShortCache{}
	c.cache = make([]*CommitShort, 0)
	return c
}

func (p *CommitShortCache) Get() *CommitShort {
	var t *CommitShort
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &CommitShort{}
	}
	return t
}
func (p *CommitShortCache) Put(t *CommitShort) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *CommitShort) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.LeaderId
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp32 = t.Count
	bs[8] = byte(tmp32 >> 24)
	bs[9] = byte(tmp32 >> 16)
	bs[10] = byte(tmp32 >> 8)
	bs[11] = byte(tmp32)
	tmp32 = t.Ballot
	bs[12] = byte(tmp32 >> 24)
	bs[13] = byte(tmp32 >> 16)
	bs[14] = byte(tmp32 >> 8)
	bs[15] = byte(tmp32)
	wire.Write(bs)
}

func (t *CommitShort) Unmarshal(wire io.Reader) error {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.LeaderId = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Count = int32(((uint32(bs[8]) << 24) | (uint32(bs[9]) << 16) | (uint32(bs[10]) << 8) | uint32(bs[11])))
	t.Ballot = int32(((uint32(bs[12]) << 24) | (uint32(bs[13]) << 16) | (uint32(bs[14]) << 8) | uint32(bs[15])))
	return nil
}

func (t *Set) New() fastrpc.Serializable {
	return new(Set)
}
func (t *Set) BinarySize() (nbytes int, sizeKnown bool) {
	return 41, true
}

type SetCache struct {
	mu    sync.Mutex
	cache []*Set
}

func NewSetCache() *SetCache {
	c := &SetCache{}
	c.cache = make([]*Set, 0)
	return c
}

func (p *SetCache) Get() *Set {
	var t *Set
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Set{}
	}
	return t
}
func (p *SetCache) Put(t *Set) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Set) Marshal(wire io.Writer) {
	var b [41]byte
	var bs []byte
	bs = b[:41]
	tmp32 := t.ReplicaID
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	bs[8] = byte(t.Write)
	tmp64 := t.Key
	bs[9] = byte(tmp64 >> 56)
	bs[10] = byte(tmp64 >> 48)
	bs[11] = byte(tmp64 >> 40)
	bs[12] = byte(tmp64 >> 32)
	bs[13] = byte(tmp64 >> 24)
	bs[14] = byte(tmp64 >> 16)
	bs[15] = byte(tmp64 >> 8)
	bs[16] = byte(tmp64)
	tmp64 = t.Payload.Tag.Timestamp
	bs[17] = byte(tmp64 >> 56)
	bs[18] = byte(tmp64 >> 48)
	bs[19] = byte(tmp64 >> 40)
	bs[20] = byte(tmp64 >> 32)
	bs[21] = byte(tmp64 >> 24)
	bs[22] = byte(tmp64 >> 16)
	bs[23] = byte(tmp64 >> 8)
	bs[24] = byte(tmp64)
	tmp64 = t.Payload.Tag.ID
	bs[25] = byte(tmp64 >> 56)
	bs[26] = byte(tmp64 >> 48)
	bs[27] = byte(tmp64 >> 40)
	bs[28] = byte(tmp64 >> 32)
	bs[29] = byte(tmp64 >> 24)
	bs[30] = byte(tmp64 >> 16)
	bs[31] = byte(tmp64 >> 8)
	bs[32] = byte(tmp64)
	tmp64 = t.Payload.Value
	bs[33] = byte(tmp64 >> 56)
	bs[34] = byte(tmp64 >> 48)
	bs[35] = byte(tmp64 >> 40)
	bs[36] = byte(tmp64 >> 32)
	bs[37] = byte(tmp64 >> 24)
	bs[38] = byte(tmp64 >> 16)
	bs[39] = byte(tmp64 >> 8)
	bs[40] = byte(tmp64)
	wire.Write(bs)
}

func (t *Set) Unmarshal(wire io.Reader) error {
	var b [41]byte
	var bs []byte
	bs = b[:41]
	if _, err := io.ReadAtLeast(wire, bs, 41); err != nil {
		return err
	}
	t.ReplicaID = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Write = uint8(bs[8])
	t.Key = int(((uint64(bs[9]) << 56) | (uint64(bs[10]) << 48) | (uint64(bs[11]) << 40) | (uint64(bs[12]) << 32) | (uint64(bs[13]) << 24) | (uint64(bs[14]) << 16) | (uint64(bs[15]) << 8) | uint64(bs[16])))
	t.Payload.Tag.Timestamp = int(((uint64(bs[17]) << 56) | (uint64(bs[18]) << 48) | (uint64(bs[19]) << 40) | (uint64(bs[20]) << 32) | (uint64(bs[21]) << 24) | (uint64(bs[22]) << 16) | (uint64(bs[23]) << 8) | uint64(bs[24])))
	t.Payload.Tag.ID = int(((uint64(bs[25]) << 56) | (uint64(bs[26]) << 48) | (uint64(bs[27]) << 40) | (uint64(bs[28]) << 32) | (uint64(bs[29]) << 24) | (uint64(bs[30]) << 16) | (uint64(bs[31]) << 8) | uint64(bs[32])))
	t.Payload.Value = int(((uint64(bs[33]) << 56) | (uint64(bs[34]) << 48) | (uint64(bs[35]) << 40) | (uint64(bs[36]) << 32) | (uint64(bs[37]) << 24) | (uint64(bs[38]) << 16) | (uint64(bs[39]) << 8) | uint64(bs[40])))
	return nil
}

func (t *PrepareReply) New() fastrpc.Serializable {
	return new(PrepareReply)
}
func (t *PrepareReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type PrepareReplyCache struct {
	mu    sync.Mutex
	cache []*PrepareReply
}

func NewPrepareReplyCache() *PrepareReplyCache {
	c := &PrepareReplyCache{}
	c.cache = make([]*PrepareReply, 0)
	return c
}

func (p *PrepareReplyCache) Get() *PrepareReply {
	var t *PrepareReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &PrepareReply{}
	}
	return t
}
func (p *PrepareReplyCache) Put(t *PrepareReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *PrepareReply) Marshal(wire io.Writer) {
	var b [10]byte
	var bs []byte
	bs = b[:9]
	tmp32 := t.Instance
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	bs[4] = byte(t.OK)
	tmp32 = t.Ballot
	bs[5] = byte(tmp32 >> 24)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 8)
	bs[8] = byte(tmp32)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Command))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Marshal(wire)
	}
}

func (t *PrepareReply) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [10]byte
	var bs []byte
	bs = b[:9]
	if _, err := io.ReadAtLeast(wire, bs, 9); err != nil {
		return err
	}
	t.Instance = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.OK = uint8(bs[4])
	t.Ballot = int32(((uint32(bs[5]) << 24) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 8) | uint32(bs[8])))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Command = make([]state.Command, alen1)
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Unmarshal(wire)
	}
	return nil
}

func (t *RMWGet) New() fastrpc.Serializable {
	return new(RMWGet)
}
func (t *RMWGet) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type RMWGetCache struct {
	mu    sync.Mutex
	cache []*RMWGet
}

func NewRMWGetCache() *RMWGetCache {
	c := &RMWGetCache{}
	c.cache = make([]*RMWGet, 0)
	return c
}

func (p *RMWGetCache) Get() *RMWGet {
	var t *RMWGet
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &RMWGet{}
	}
	return t
}
func (p *RMWGetCache) Put(t *RMWGet) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *RMWGet) Marshal(wire io.Writer) {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	tmp32 := t.LeaderId
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp32 = t.Ballot
	bs[8] = byte(tmp32 >> 24)
	bs[9] = byte(tmp32 >> 16)
	bs[10] = byte(tmp32 >> 8)
	bs[11] = byte(tmp32)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Command))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Marshal(wire)
	}
}

func (t *RMWGet) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [12]byte
	var bs []byte
	bs = b[:12]
	if _, err := io.ReadAtLeast(wire, bs, 12); err != nil {
		return err
	}
	t.LeaderId = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Ballot = int32(((uint32(bs[8]) << 24) | (uint32(bs[9]) << 16) | (uint32(bs[10]) << 8) | uint32(bs[11])))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Command = make([]state.Command, alen1)
	for i := int64(0); i < alen1; i++ {
		t.Command[i].Unmarshal(wire)
	}
	return nil
}

func (t *Tag) New() fastrpc.Serializable {
	return new(Tag)
}
func (t *Tag) BinarySize() (nbytes int, sizeKnown bool) {
	return 16, true
}

type TagCache struct {
	mu    sync.Mutex
	cache []*Tag
}

func NewTagCache() *TagCache {
	c := &TagCache{}
	c.cache = make([]*Tag, 0)
	return c
}

func (p *TagCache) Get() *Tag {
	var t *Tag
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Tag{}
	}
	return t
}
func (p *TagCache) Put(t *Tag) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Tag) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp64 := t.Timestamp
	bs[0] = byte(tmp64 >> 56)
	bs[1] = byte(tmp64 >> 48)
	bs[2] = byte(tmp64 >> 40)
	bs[3] = byte(tmp64 >> 32)
	bs[4] = byte(tmp64 >> 24)
	bs[5] = byte(tmp64 >> 16)
	bs[6] = byte(tmp64 >> 8)
	bs[7] = byte(tmp64)
	tmp64 = t.ID
	bs[8] = byte(tmp64 >> 56)
	bs[9] = byte(tmp64 >> 48)
	bs[10] = byte(tmp64 >> 40)
	bs[11] = byte(tmp64 >> 32)
	bs[12] = byte(tmp64 >> 24)
	bs[13] = byte(tmp64 >> 16)
	bs[14] = byte(tmp64 >> 8)
	bs[15] = byte(tmp64)
	wire.Write(bs)
}

func (t *Tag) Unmarshal(wire io.Reader) error {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Timestamp = int(((uint64(bs[0]) << 56) | (uint64(bs[1]) << 48) | (uint64(bs[2]) << 40) | (uint64(bs[3]) << 32) | (uint64(bs[4]) << 24) | (uint64(bs[5]) << 16) | (uint64(bs[6]) << 8) | uint64(bs[7])))
	t.ID = int(((uint64(bs[8]) << 56) | (uint64(bs[9]) << 48) | (uint64(bs[10]) << 40) | (uint64(bs[11]) << 32) | (uint64(bs[12]) << 24) | (uint64(bs[13]) << 16) | (uint64(bs[14]) << 8) | uint64(bs[15])))
	return nil
}

func (t *Payload) New() fastrpc.Serializable {
	return new(Payload)
}
func (t *Payload) BinarySize() (nbytes int, sizeKnown bool) {
	return 24, true
}

type PayloadCache struct {
	mu    sync.Mutex
	cache []*Payload
}

func NewPayloadCache() *PayloadCache {
	c := &PayloadCache{}
	c.cache = make([]*Payload, 0)
	return c
}

func (p *PayloadCache) Get() *Payload {
	var t *Payload
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Payload{}
	}
	return t
}
func (p *PayloadCache) Put(t *Payload) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Payload) Marshal(wire io.Writer) {
	var b [24]byte
	var bs []byte
	bs = b[:24]
	tmp64 := t.Tag.Timestamp
	bs[0] = byte(tmp64 >> 56)
	bs[1] = byte(tmp64 >> 48)
	bs[2] = byte(tmp64 >> 40)
	bs[3] = byte(tmp64 >> 32)
	bs[4] = byte(tmp64 >> 24)
	bs[5] = byte(tmp64 >> 16)
	bs[6] = byte(tmp64 >> 8)
	bs[7] = byte(tmp64)
	tmp64 = t.Tag.ID
	bs[8] = byte(tmp64 >> 56)
	bs[9] = byte(tmp64 >> 48)
	bs[10] = byte(tmp64 >> 40)
	bs[11] = byte(tmp64 >> 32)
	bs[12] = byte(tmp64 >> 24)
	bs[13] = byte(tmp64 >> 16)
	bs[14] = byte(tmp64 >> 8)
	bs[15] = byte(tmp64)
	tmp64 = t.Value
	bs[16] = byte(tmp64 >> 56)
	bs[17] = byte(tmp64 >> 48)
	bs[18] = byte(tmp64 >> 40)
	bs[19] = byte(tmp64 >> 32)
	bs[20] = byte(tmp64 >> 24)
	bs[21] = byte(tmp64 >> 16)
	bs[22] = byte(tmp64 >> 8)
	bs[23] = byte(tmp64)
	wire.Write(bs)
}

func (t *Payload) Unmarshal(wire io.Reader) error {
	var b [24]byte
	var bs []byte
	bs = b[:24]
	if _, err := io.ReadAtLeast(wire, bs, 24); err != nil {
		return err
	}
	t.Tag.Timestamp = int(((uint64(bs[0]) << 56) | (uint64(bs[1]) << 48) | (uint64(bs[2]) << 40) | (uint64(bs[3]) << 32) | (uint64(bs[4]) << 24) | (uint64(bs[5]) << 16) | (uint64(bs[6]) << 8) | uint64(bs[7])))
	t.Tag.ID = int(((uint64(bs[8]) << 56) | (uint64(bs[9]) << 48) | (uint64(bs[10]) << 40) | (uint64(bs[11]) << 32) | (uint64(bs[12]) << 24) | (uint64(bs[13]) << 16) | (uint64(bs[14]) << 8) | uint64(bs[15])))
	t.Value = int(((uint64(bs[16]) << 56) | (uint64(bs[17]) << 48) | (uint64(bs[18]) << 40) | (uint64(bs[19]) << 32) | (uint64(bs[20]) << 24) | (uint64(bs[21]) << 16) | (uint64(bs[22]) << 8) | uint64(bs[23])))
	return nil
}

func (t *Prepare) New() fastrpc.Serializable {
	return new(Prepare)
}
func (t *Prepare) BinarySize() (nbytes int, sizeKnown bool) {
	return 13, true
}

type PrepareCache struct {
	mu    sync.Mutex
	cache []*Prepare
}

func NewPrepareCache() *PrepareCache {
	c := &PrepareCache{}
	c.cache = make([]*Prepare, 0)
	return c
}

func (p *PrepareCache) Get() *Prepare {
	var t *Prepare
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Prepare{}
	}
	return t
}
func (p *PrepareCache) Put(t *Prepare) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Prepare) Marshal(wire io.Writer) {
	var b [13]byte
	var bs []byte
	bs = b[:13]
	tmp32 := t.LeaderId
	bs[0] = byte(tmp32 >> 24)
	bs[1] = byte(tmp32 >> 16)
	bs[2] = byte(tmp32 >> 8)
	bs[3] = byte(tmp32)
	tmp32 = t.Instance
	bs[4] = byte(tmp32 >> 24)
	bs[5] = byte(tmp32 >> 16)
	bs[6] = byte(tmp32 >> 8)
	bs[7] = byte(tmp32)
	tmp32 = t.Ballot
	bs[8] = byte(tmp32 >> 24)
	bs[9] = byte(tmp32 >> 16)
	bs[10] = byte(tmp32 >> 8)
	bs[11] = byte(tmp32)
	bs[12] = byte(t.ToInfinity)
	wire.Write(bs)
}

func (t *Prepare) Unmarshal(wire io.Reader) error {
	var b [13]byte
	var bs []byte
	bs = b[:13]
	if _, err := io.ReadAtLeast(wire, bs, 13); err != nil {
		return err
	}
	t.LeaderId = int32(((uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])))
	t.Instance = int32(((uint32(bs[4]) << 24) | (uint32(bs[5]) << 16) | (uint32(bs[6]) << 8) | uint32(bs[7])))
	t.Ballot = int32(((uint32(bs[8]) << 24) | (uint32(bs[9]) << 16) | (uint32(bs[10]) << 8) | uint32(bs[11])))
	t.ToInfinity = uint8(bs[12])
	return nil
}
