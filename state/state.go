package state

import (
	"github.com/alexjipark/datastreet/types"
	bctypes "github.com/tendermint/basecoin/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	"fmt"
)


// ?? CONTRACT : State should be quick to copy.
// See CacheWrap()

type State struct {
	chainID	string
	store 	bctypes.KVStore
	cache	*bctypes.KVCache		// optional
}

func NewState (store bctypes.KVStore) *State {
	return &State {
		chainID: "",
		store: store,
	}
}

func (s* State) SetChainID(chainID string) {
	s.chainID = chainID
}

func (s* State) GetChainID() string {
	if s.chainID == "" {
		PanicSanity("Expected to have set SetChainID")
	}
	return s.chainID
}

func (s *State) Get(key []byte) (value []byte) {
	return s.store.Get(key)
}

func (s *State) Set(key []byte, value []byte) {
	s.store.Set(key, value)
}

func (s *State) GetAccount(addr []byte) *types.Account {
	return GetAccount(s.store, addr)
}

func (s *State) SetAccount(addr []byte, acc *types.Account) {
	SetAccount(s.store, addr, acc)
}

func (s *State) CacheWrap() *State {
	cache := bctypes.NewKVCache(s.store)
	return &State {
		chainID: s.chainID,
		store: cache,
		cache: cache,
	}
}

// NOTE : errors if s is not from CacheWrap() ????
func (s *State) CacheSync() {
	s.cache.Sync()
}

//--------------------------
func AccountKey(addr []byte) []byte {
	return append([]byte("base/a/"), addr...)
}

//
func GetAccount (store bctypes.KVStore, addr []byte) *types.Account {
	//data := store.Get(AccountKey(addr))
	data := store.Get(addr)
	if len(data) == 0 {
		return nil
	}
	var acc *types.Account
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic (Fmt("Error Reading Account %X error: %v",
				data, err.Error()))
	}
	return acc
}

func SetAccount (store bctypes.KVStore, addr []byte, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	//store.Set (AccountKey(addr), accBytes)
	//===== Test
	fmt.Printf("\nSetAccount : Key[%x] Value[%X]", addr, accBytes)
	store.Set (addr, accBytes)
}