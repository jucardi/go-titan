package mongo

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"os"
	"sync/atomic"
	"time"
)

/*
This file contains utilities that were available in the old mgo.v2 package and are currently
unavailable in the new mongodb Golang client library
*/

var (
	machineId       = readMachineId()
	processId       = os.Getpid()
	objectIdCounter = readRandomUint32()
)

// readRandomUint32 returns a random objectIdCounter.
func readRandomUint32() uint32 {
	// We've found systems hanging in this function due to lack of entropy.
	// The randomness of these bytes is just preventing nearby clashes, so
	// just look at the time.
	return uint32(time.Now().UnixNano())
}

// NewObjectId returns a new unique ObjectId.
func NewObjectId() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().UTC().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	b[7] = byte(processId >> 8)
	b[8] = byte(processId)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return hex.EncodeToString(b[:])
}

// readMachineId generates and returns a machine id.
// If this function fails to get the hostname it will cause a runtime error.
func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		n := uint32(time.Now().UnixNano())
		sum[0] = byte(n >> 0)
		sum[1] = byte(n >> 8)
		sum[2] = byte(n >> 16)
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}
