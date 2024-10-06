package gocache

import (
	"fmt"
	"log"
	"sync"

	pb "zjkung.github.com/g-cach-e/gocachepb"
	"zjkung.github.com/g-cach-e/singleflight"
)

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ReadOnlyByte, err error) {
	// each key is only fetched once (either locally or remotely)
	// regardless of the number of concurrent callers.
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GoCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return view.(ReadOnlyByte), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ReadOnlyByte, error) {
	req := &pb.Request{Group: g.name, Key: key}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ReadOnlyByte{}, err
	}
	return ReadOnlyByte{bytes: res.Value}, nil
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: *NewCache(cacheBytes),
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ReadOnlyByte, error) {
	if key == "" {
		return ReadOnlyByte{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[GoCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) getLocally(key string) (ReadOnlyByte, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ReadOnlyByte{}, err

	}
	value := ReadOnlyByte{bytes: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ReadOnlyByte) {
	g.mainCache.Add(key, value)
}
