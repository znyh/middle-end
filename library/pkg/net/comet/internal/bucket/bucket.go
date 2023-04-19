package bucket

import (
    "sync"

    "github.com/znyh/middle-end/library/pkg/net/comet/internal/channel"
)

// Bucket is a channel holder.
type Bucket struct {
    cLock  sync.RWMutex                // protect the channels for chs
    chs    map[string]*channel.Channel // map sub key to a channel
    ipCnts map[string]int32
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(cns int) (b *Bucket) {
    b = new(Bucket)
    b.chs = make(map[string]*channel.Channel, cns)
    b.ipCnts = make(map[string]int32)
    return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
    return len(b.chs)
}

// Put put a channel according with sub key.
func (b *Bucket) Put(ch *channel.Channel) {
    b.cLock.Lock()
    // close old channel
    if dch := b.chs[ch.Key]; dch != nil {
        dch.Close()
    }
    b.chs[ch.Key] = ch
    b.ipCnts[ch.IP]++
    b.cLock.Unlock()
    return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(dch *channel.Channel) {
    var (
        ok bool
        ch *channel.Channel
    )
    b.cLock.Lock()
    if ch, ok = b.chs[dch.Key]; ok {
        if ch == dch {
            delete(b.chs, ch.Key)
        }
        // ip counter
        if b.ipCnts[ch.IP] > 1 {
            b.ipCnts[ch.IP]--
        } else {
            delete(b.ipCnts, ch.IP)
        }
    }
    b.cLock.Unlock()
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *channel.Channel) {
    b.cLock.RLock()
    ch = b.chs[key]
    b.cLock.RUnlock()
    return
}

// IPCount get ip count.
func (b *Bucket) IPCount() (res map[string]struct{}) {
    var (
        ip string
    )
    b.cLock.RLock()
    res = make(map[string]struct{}, len(b.ipCnts))
    for ip = range b.ipCnts {
        res[ip] = struct{}{}
    }
    b.cLock.RUnlock()
    return
}
