
package eth

import (
	"context"
    "fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
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

