package service

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// {
// 	baseFeePerGas: 7,
//   blobGasUsed: 0,
//   difficulty: 0,
//   excessBlobGas: 0,
//   extraData: "0xd883010f0a846765746888676f312e32342e32856c696e7578",
//   gasLimit: 30000000,
//   gasUsed: 123005,
//   hash: "0x30591ec14026fc94781f067c14b42c2d6fe33dd89471dc1f40265c3a65e28236",
//   logsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
//   miner: "0x8943545177806ed17b9f23f0a21ee5948ecaa776",
//   mixHash: "0x76ec60209b272256265b3893354b5e07b2a6915a30a6bc8c2d060866c06e53a6",
//   nonce: "0x0000000000000000",
//   number: 1094,
//   parentBeaconBlockRoot: "0xdf24898ba4eb46d21d99ebe47cce1caa26d0d89cc2a6d6a21ace435400b7cfcb",
//   parentHash: "0xa4bc1d0221f34e5034bbd6fae25f2d3d06110d1c6ae0854ca6812d7a1d31b974",
//   receiptsRoot: "0x7452e34081d0d764d5b9be4535a5217e640d12cc301b051a00e4ebff8c4d5051",
//   sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
//   size: 1059,
//   stateRoot: "0xa8e0331f21c9c1f06a4b005e515267a6f80aff204416dc39542971588589abc4",
//   timestamp: 1745674074,
//   transactions: ["0x92110e18d6c16b95be696fa76cae40b4824e7be4229d46e3c060ad8ebb606057"],
//   transactionsRoot: "0x9ff042ab9c3b7fb5927367d12714e7839a422e6b4e457d8e3bc706f93275acb2",
//   uncles: [],
//   withdrawals: [],
//   withdrawalsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
// }

// Block 表示以太坊区块的数据结构
type Block struct {
	BaseFeePerGas         uint64   `json:"baseFeePerGas"`
	BlobGasUsed           uint64   `json:"blobGasUsed"`
	Difficulty            uint64   `json:"difficulty"`
	ExcessBlobGas         uint64   `json:"excessBlobGas"`
	ExtraData             string   `json:"extraData"`
	GasLimit              uint64   `json:"gasLimit"`
	GasUsed               uint64   `json:"gasUsed"`
	Hash                  string   `json:"hash"`
	LogsBloom             string   `json:"logsBloom"`
	Miner                 string   `json:"miner"`
	MixHash               string   `json:"mixHash"`
	Nonce                 string   `json:"nonce"`
	Number                uint64   `json:"number"`
	ParentBeaconBlockRoot string   `json:"parentBeaconBlockRoot"`
	ParentHash            string   `json:"parentHash"`
	ReceiptsRoot          string   `json:"receiptsRoot"`
	Sha3Uncles            string   `json:"sha3Uncles"`
	Size                  uint64   `json:"size"`
	StateRoot             string   `json:"stateRoot"`
	Timestamp             uint64   `json:"timestamp"`
	Transactions          []string `json:"transactions"`
	TransactionsRoot      string   `json:"transactionsRoot"`
	Uncles                []string `json:"uncles"`
	Withdrawals           []string `json:"withdrawals"`
	WithdrawalsRoot       string   `json:"withdrawalsRoot"`
}

// FromTypesBlock 将 types.Block 转换为自定义的 Block 结构体
func FromTypesBlock(b *types.Block) *Block {
	if b == nil {
		return nil
	}

	// 获取所有交易哈希
	txHashes := make([]string, len(b.Transactions()))
	for i, tx := range b.Transactions() {
		txHashes[i] = tx.Hash().Hex()
	}

	// 获取所有叔块哈希
	uncleHashes := make([]string, len(b.Uncles()))
	for i, uncle := range b.Uncles() {
		uncleHashes[i] = uncle.Hash().Hex()
	}

	// 获取所有提款记录哈希
	withdrawalHashes := make([]string, 0)
	if withdrawals := b.Withdrawals(); withdrawals != nil {
		withdrawalHashes = make([]string, len(withdrawals))
		for i, withdrawal := range withdrawals {
			withdrawalHashes[i] = withdrawal.Address.Hex()
		}
	}

	header := b.Header()
	baseFee := b.BaseFee()
	var baseFeePerGas uint64
	if baseFee != nil {
		baseFeePerGas = baseFee.Uint64()
	}

	return &Block{
		BaseFeePerGas:         baseFeePerGas,
		BlobGasUsed:           *header.BlobGasUsed,
		Difficulty:            header.Difficulty.Uint64(),
		ExcessBlobGas:         *header.ExcessBlobGas,
		ExtraData:             "0x" + hex(header.Extra),
		GasLimit:              header.GasLimit,
		GasUsed:               header.GasUsed,
		Hash:                  b.Hash().Hex(),
		LogsBloom:             "0x" + hex(header.Bloom.Bytes()),
		Miner:                 header.Coinbase.Hex(),
		MixHash:               header.MixDigest.Hex(),
		Nonce:                 "0x" + hex(header.Nonce[:]),
		Number:                header.Number.Uint64(),
		ParentBeaconBlockRoot: header.ParentHash.Hex(),
		ParentHash:            header.ParentHash.Hex(),
		ReceiptsRoot:          header.ReceiptHash.Hex(),
		Sha3Uncles:            header.UncleHash.Hex(),
		Size:                  uint64(b.Size()),
		StateRoot:             header.Root.Hex(),
		Timestamp:             header.Time,
		Transactions:          txHashes,
		TransactionsRoot:      header.TxHash.Hex(),
		Uncles:                uncleHashes,
		Withdrawals:           withdrawalHashes,
		WithdrawalsRoot:       header.WithdrawalsHash.Hex(),
	}
}

// hex 将字节切片转换为十六进制字符串（不带0x前缀）
func hex(b []byte) string {
	const hexDigits = "0123456789abcdef"
	res := make([]byte, len(b)*2)
	for i, v := range b {
		res[i*2] = hexDigits[v>>4]
		res[i*2+1] = hexDigits[v&0x0f]
	}
	return string(res)
}
