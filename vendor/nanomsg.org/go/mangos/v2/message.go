// Copyright 2016 The Mangos Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mangos

import (
	"sync"
)

// Message encapsulates the messages that we exchange back and forth.  The
// meaning of the Header and Body fields, and where the splits occur, will
// vary depending on the protocol.  Note however that any headers applied by
// transport layers (including TCP/ethernet headers, and SP protocol
// independent length headers), are *not* included in the Header.
type Message struct {
	// Header carries any protocol (SP) specific header.  Applications
	// should not modify or use this unless they are using Raw mode.
	// No user data may be placed here.
	Header []byte

	// Body carries the body of the message.  This can also be thought
	// of as the message "payload".
	Body []byte

	// Pipe may be set on message receipt, to indicate the Pipe from
	// which the Message was received.  There are no guarantees that the
	// Pipe is still active, and applications should only use this for
	// informational purposes.
	Pipe Pipe

	bbuf  []byte
	hbuf  []byte
	bsize int
	pool  *sync.Pool
}

type msgCacheInfo struct {
	maxbody int
	pool    *sync.Pool
}

func newMsg(sz int) *Message {
	m := &Message{}
	m.bbuf = make([]byte, 0, sz)
	m.hbuf = make([]byte, 0, 32)
	m.bsize = sz
	return m
}

// We can tweak these!
var messageCache = []msgCacheInfo{
	{
		maxbody: 64,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(64) },
		},
	}, {
		maxbody: 128,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(128) },
		},
	}, {
		maxbody: 256,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(256) },
		},
	}, {
		maxbody: 512,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(512) },
		},
	}, {
		maxbody: 1024,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(1024) },
		},
	}, {
		maxbody: 4096,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(4096) },
		},
	}, {
		maxbody: 8192,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(8192) },
		},
	}, {
		maxbody: 65536,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(65536) },
		},
	},
}

// Free releases the message to the pool from which it was allocated.
// While this is not strictly necessary thanks to GC, doing so allows
// for the resources to be recycled without engaging GC.  This can have
// rather substantial benefits for performance.
func (m *Message) Free() {
	for i := range messageCache {
		if m.bsize == messageCache[i].maxbody {
			messageCache[i].pool.Put(m)
			return
		}
	}
}

// Dup creates a "duplicate" message.
// Reference counting was found to be error prone, so we have elected
// to simply make a full copy of the message for now.
func (m *Message) Dup() *Message {
	dup := NewMessage(len(m.Body))
	dup.Body = append(dup.Body, m.Body...)
	dup.Header = append(dup.Header, m.Header...)
	dup.Pipe = m.Pipe
	return dup
}

// NewMessage is the supported way to obtain a new Message.  This makes
// use of a "cache" which greatly reduces the load on the garbage collector.
func NewMessage(sz int) *Message {
	var m *Message
	for i := range messageCache {
		if sz < messageCache[i].maxbody {
			m = messageCache[i].pool.Get().(*Message)
			break
		}
	}
	if m == nil {
		m = newMsg(sz)
	}

	m.Body = m.bbuf
	m.Header = m.hbuf
	return m
}
