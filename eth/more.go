
package eth

import (
	"context"
    "fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
//  "github.com/ethereum/go-ethereum/trie"
//	"github.com/ethereum/go-ethereum/log"
)

type PublicQueryAPI struct {
	b *Ethereum
}

func NewPublicQueryAPI(b *Ethereum) *PublicQueryAPI {
	return &PublicQueryAPI{b: b}
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicQueryAPI) AccountProof(ctx context.Context, bhash hexutil.Bytes, key hexutil.Bytes) ([]hexutil.Bytes, error) {
    db := api.b.ChainDb()
    hash := common.BytesToHash(bhash)
    header := core.GetHeader(db, hash, core.GetBlockNumber(db, hash))
    if header == nil {
        return nil, fmt.Errorf("Did not find the header")
    }
    statedb, err := api.b.BlockChain().State()
    if err != nil {
        return nil, err
    }
    trie, err := statedb.Database().OpenTrie(header.Root)
    if err != nil {
        return nil, err
    }
    
    var proof light.NodeList

    trie.Prove(key, 0, &proof)
    
    res := make([]hexutil.Bytes, len(proof))
    for i, v := range proof {
        res[i] = hexutil.Bytes(v)
    }
    return res, nil
}

func (api *PublicQueryAPI) AccountRLP(ctx context.Context, bhash hexutil.Bytes, key hexutil.Bytes) (hexutil.Bytes, error) {
    db := api.b.ChainDb()
    hash := common.BytesToHash(bhash)
    header := core.GetHeader(db, hash, core.GetBlockNumber(db, hash))
    if header == nil {
        return nil, fmt.Errorf("Did not find the header")
    }
    statedb, err := api.b.BlockChain().State()
    if err != nil {
        return nil, err
    }
    trie, err := statedb.Database().OpenTrie(header.Root)
    if err != nil {
        return nil, err
    }
    
    blob, err := trie.TryGet(key)
	if err != nil {
		return nil, fmt.Errorf("try get failure %s", err)
	}
    
    return blob, nil
}

func (api *PublicQueryAPI) GetBlockHeader(ctx context.Context, bhash hexutil.Bytes) (hexutil.Bytes, error) {
    db := api.b.ChainDb()
    hash := common.BytesToHash(bhash)
    header := core.GetHeader(db, hash, core.GetBlockNumber(db, hash))
    if header == nil {
        return nil, fmt.Errorf("Did not find the header")
    }
    res, err := rlp.EncodeToBytes(header)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (api *PublicQueryAPI) getAccount(statedb *state.StateDB, root common.Hash, key hexutil.Bytes) (state.Account, error) {
    trie, err := statedb.Database().OpenTrie(root)
    if err != nil {
        return state.Account{}, err
    }
    
    blob, err := trie.TryGet(key)
	if err != nil {
		return state.Account{}, fmt.Errorf("try get failure %s", err)
	}
	var account state.Account
	if err = rlp.DecodeBytes(blob, &account); err != nil {
        // log.Warn("Got RLP", "blob", blob)
		return state.Account{}, fmt.Errorf("rlp failure %s", err)
	}
	return account, nil
}

func (api *PublicQueryAPI) StorageProof(ctx context.Context, bhash hexutil.Bytes, addr hexutil.Bytes, ptr hexutil.Bytes) ([]hexutil.Bytes, error) {
    db := api.b.ChainDb()
    hash := common.BytesToHash(bhash)
    header := core.GetHeader(db, hash, core.GetBlockNumber(db, hash))
    if header == nil {
        return nil, fmt.Errorf("Did not find the header")
    }
    statedb, err := api.b.BlockChain().State()
    if err != nil {
        return nil, err
    }

//    account, err := api.getAccount(statedb, header.Root, common.BytesToHash(addr))
    account, err := api.getAccount(statedb, header.Root, addr)
	if err != nil {
        return nil, fmt.Errorf("Cannot get account %s", err)
	}
    trie, err := statedb.Database().OpenStorageTrie(common.BytesToHash(addr), account.Root)
    if err != nil {
        return nil, fmt.Errorf("Cannot open db %s", err)
    }
    
    var proof light.NodeList

    trie.Prove(ptr, 0, &proof)
    
    res := make([]hexutil.Bytes, len(proof))
    for i, v := range proof {
        res[i] = hexutil.Bytes(v)
    }
    return res, nil

}

