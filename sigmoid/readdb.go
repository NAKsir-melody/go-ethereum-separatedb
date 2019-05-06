package main

import (
	"fmt"
    "github.com/ethereum/go-ethereum/core/rawdb"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/core/state"
    "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	_ "github.com/ethereum/go-ethereum/common"
    "encoding/binary"
    _ "github.com/davecgh/go-spew/spew"
)

var (
    // databaseVerisionKey tracks the current database version.
    databaseVerisionKey = []byte("DatabaseVersion")

    // headHeaderKey tracks the latest know header's hash.
    headHeaderKey = []byte("LastHeader")

    // headBlockKey tracks the latest know full block's hash.
    headBlockKey = []byte("LastBlock")

    // headFastBlockKey tracks the latest known incomplete block's hash during fast sync.
    headFastBlockKey = []byte("LastFast")

    // fastTrieProgressKey tracks the number of trie entries imported during fast sync.
    fastTrieProgressKey = []byte("TrieSync")

    // Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
    headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
    headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
    headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
    headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

    blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
    blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

    txLookupPrefix  = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata
    bloomBitsPrefix = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits

    preimagePrefix = []byte("secure-key-")      // preimagePrefix + hash -> preimage
    configPrefix   = []byte("ethereum-config-") // config prefix for the db

    // Chain index prefixes (use `i` + single byte to avoid mixing data types).
    BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress
)

func main() {

    //open db
	ldb, err := rawdb.NewLevelDBDatabase("mainnet/geth/chaindata", 0, 0, "");
	if err != nil {
		fmt.Printf("open error\n" + err.Error())
		return
	}

    //get header chain (latest blocks header)
    latest_header_hash, err := ldb.Get(headHeaderKey)
	fmt.Printf("%x\n", latest_header_hash)

    //get last block number
    block_no_key := append(headerNumberPrefix,latest_header_hash...)
    block_no, _ := ldb.Get(block_no_key)
    block_num := binary.BigEndian.Uint64(block_no)
	fmt.Printf("%x\n", block_no)
	fmt.Printf("%d\n", block_num)


	var b types.Header
	var i uint64
	trieDB := trie.NewDatabase(ldb)
    //do loop until parents' hash reached genesis
	//for  i = 0 ; i < 1; i++ {
		fmt.Printf("%d ================================\n", i)
		//new_block_num := make([]byte, 8)
		//binary.BigEndian.PutUint64(new_block_num, data-i)
		//fmt.Printf("%x\n", new_block_num)

		//get block header
		//temp := append(new_block_num, latest_header_hash...);
		temp := append(block_no, latest_header_hash...);
		block_header_key := append(headerPrefix[:], temp...);
		block_header_rlp, _ := ldb.Get(block_header_key)
		rlp.DecodeBytes(block_header_rlp, &b)
		fmt.Printf("%x\n", b)
		fmt.Printf("%x\n", b.Root)

		state_trie, err := trie.New(b.Root, trieDB) //state
		if err != nil {
			fmt.Printf(err.Error() + "\n")
	//		continue
            return
		}
        Value, _ := ldb.Get(b.Root.Bytes())
		fmt.Printf("%x\n", Value)

	it := trie.NewIterator(state_trie.NodeIterator(b.Root.Bytes()))
    for it.Next() {
        //state_trie.getKey(it.Key)
        var data state.Account
        if err := rlp.DecodeBytes(it.Value, &data); err != nil {
            panic(err)
        }
        fmt.Printf("%x\n", data)
    }

//		latest_header_hash = b.ParentHash.Bytes();
	//}

}

