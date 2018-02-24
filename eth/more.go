
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
	"github.com/ethereum/go-ethereum/log"
)

type PublicQueryAPI struct {
	b *Ethereum
}

func NewPublicQueryAPI(b *Ethereum) *PublicQueryAPI {
	return &PublicQueryAPI{b: b}
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicQueryAPI) AccountProof(ctx context.Context, bhash hexutil.Bytes, key hexutil.Bytes) (hexutil.Bytes, error) {
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
    
    /*
    test how to get stuff from trie
    bts, err := trie.TryGet(key)
    
    if err != nil {
       return nil, err
    }
    
    if bts != nil {
       return bts, nil
    }
    */
    
    var proof light.NodeList
    var proofs [][]rlp.RawValue

    trie.Prove(key, 0, &proof)

	proofs = append(proofs, proof)
    rlp, err := rlp.EncodeToBytes(proofs)
    if err != nil {
        return nil, err
    }
    return rlp, nil
}

func (api *PublicQueryAPI) getAccount(statedb *state.StateDB, root common.Hash, /* hash common.Hash */ key hexutil.Bytes) (state.Account, error) {
    trie, err := statedb.Database().OpenTrie(root)
    if err != nil {
        return state.Account{}, err
    }
    
    blob, err := trie.TryGet(key)
/*    trie, err := trie.New(root, statedb.Database().TrieDB())
	if err != nil {
        return state.Account{}, fmt.Errorf("trie failure %s", err)
	}
	blob, err := trie.TryGet(hash[:]) */
	if err != nil {
		return state.Account{}, fmt.Errorf("try get failure %s", err)
	}
	var account state.Account
	if err = rlp.DecodeBytes(blob, &account); err != nil {
        log.Warn("Got RLP", "blob", blob)
		return state.Account{}, fmt.Errorf("rlp failure %s", err)
	}
	return account, nil
}

func (api *PublicQueryAPI) StorageProof(ctx context.Context, bhash hexutil.Bytes, addr hexutil.Bytes, ptr hexutil.Bytes) (hexutil.Bytes, error) {
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
    var proofs [][]rlp.RawValue

    trie.Prove(ptr, 0, &proof)

	proofs = append(proofs, proof)
    rlp, err := rlp.EncodeToBytes(proofs)
    if err != nil {
        return nil, err
    }
    return rlp, nil
}

