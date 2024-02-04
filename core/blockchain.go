package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"sync"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/go-kit/log"
)

type Blockchain struct {
	logger log.Logger
	store  Storage
	// TODO: double check this!
	lock         sync.RWMutex
	headers      []*types.Header
	blocks       []*types.Block
	txStore      map[common.Hash]*types.Transaction
	blockStore   map[common.Hash]*types.Block
	accountState *AccountState
	stateLock    sync.RWMutex
	validator    Validator
	// TODO: make this an interface.
	contractState *State
	db            *barreldb.BarrelDatabase
}

func NewBlockchain(l log.Logger, privateKey *crypto.PrivateKey, genesis *types.Block) (*Blockchain, error) {
	// We should create all states inside the scope of the newblockchain.
	// TODO: read this from disk later on
	accountState := NewAccountState()

	//coinbase := privateKey.PublicKey()
	//accountState.CreateAccount(coinbase.Address())

	db, _ := barreldb.New()

	/*
		bc 객체가 없는 영역에서 db 활용 Sample

		_ = db.CreateTable("block", barreldb.BlockPrefix)
		_ = db.CreateBlock("kim", "youngmin")
		data, _ := db.GetBlock("kim")
	*/

	err := db.CreateTable(barreldb.HashBlockTableName, barreldb.HashBlockPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.HeightBlockTableName, barreldb.HeightBlockPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastBlockTableName, barreldb.LastBlockPrefix)
	if err != nil {
		return nil, err
	}

	err = db.CreateTable(barreldb.HashTxTableName, barreldb.HashTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.NumberTxTableName, barreldb.NumberTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastTxTableName, barreldb.LastTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastTxNumberTableName, barreldb.LastTxNumberPrefix)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		contractState: NewState(),
		headers:       []*types.Header{},
		store:         NewMemorystore(),
		logger:        l,
		accountState:  accountState,
		blockStore:    make(map[common.Hash]*types.Block),
		txStore:       make(map[common.Hash]*types.Transaction),
		db:            db,
	}
	bc.validator = NewBlockValidator(bc)

	if genesis != nil {
		err = bc.addBlockWithoutValidation(genesis)
	}

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *types.Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) handleNativeTransfer(tx *types.Transaction) error {
	bc.logger.Log(
		"msg", "handle native token transfer",
		"from", tx.From,
		"to", tx.To,
		"value", tx.Value)

	return bc.accountState.Transfer(tx.From.Address(), tx.To.Address(), tx.Value)
}

func (bc *Blockchain) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	block, ok := bc.blockStore[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash (%s) not found", hash)
	}

	return block, nil
}

func (bc *Blockchain) GetBlock(height uint32) (*types.Block, error) {
	if int32(height) > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}
	fmt.Println("height: ", height, "bc.height: ", bc.Height())
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) GetHeader(height uint32) (*types.Header, error) {
	if int32(height) > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()
	fmt.Println("heheheee: ", height)
	return bc.headers[height], nil
}

func (bc *Blockchain) GetTxByHash(hash common.Hash) (*types.Transaction, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	tx, ok := bc.txStore[hash]
	if !ok {
		return nil, fmt.Errorf("could not find tx with hash (%s)", hash)
	}

	return tx, nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return int32(height) <= bc.Height()
}

// [0, 1, 2 ,3] => 4 len
// [0, 1, 2 ,3] => 3 height
func (bc *Blockchain) Height() int32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return int32(len(bc.headers) - 1)
}

func (bc *Blockchain) handleTransaction(tx *types.Transaction) error {
	fmt.Println("txtx: ", tx)
	// If we have data inside execute that data on the VM.
	if len(tx.Data) > 0 {
		bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.GetHash(&types.TxHasher{}))

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	if err := bc.handleNativeTransfer(tx); err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *types.Block) error {
	bc.stateLock.Lock()
	for i := 0; i < len(b.Transactions); i++ {
		if err := bc.handleTransaction(b.Transactions[i]); err != nil {
			bc.logger.Log("error", err.Error())

			b.Transactions[i] = b.Transactions[len(b.Transactions)-1]
			b.Transactions = b.Transactions[:len(b.Transactions)-1]

			continue
		}
	}
	bc.stateLock.Unlock()

	// fmt.Println("========ACCOUNT STATE==============")
	// fmt.Printf("%+v\n", bc.accountState.accounts)
	// fmt.Println("========ACCOUNT STATE==============")

	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)

	if err := bc.WriteBlockWithHash(b.GetHash(types.BlockHasher{}), b); err != nil {
		return err
	}

	if err := bc.WriteBlockWithHeight(b.Height, b); err != nil {
		return err
	}

	if err := bc.WriteLastBlock(b); err != nil {
		return err
	}

	//data, _ := bc.ReadBlockByHash(b.GetHash(types.BlockHasher{}))
	//val := hexutil.Encode(data.Hash.ToSlice())
	//fmt.Println("hashblock::: ", val)
	//fmt.Println("fefefefefk::: ", data.Hash.String())

	//data, _ = bc.ReadBlockByHeight(b.Height)
	//fmt.Println("heightblock::: ", data)

	//data, _ = bc.ReadLastBlock()
	//fmt.Println("Lastblock::: ", data)

	bc.blockStore[b.GetHash(types.BlockHasher{})] = b

	for _, tx := range b.Transactions {
		bc.txStore[tx.GetHash(types.TxHasher{})] = tx

		nextTxNumber := uint32(0)
		number, err := bc.ReadLastTxNumber()
		if err != nil {

			//return err
		}

		if number != nil {
			nextTxNumber = *number + 1
		}

		if err := bc.WriteTxWithHash(tx.GetHash(types.TxHasher{}), tx); err != nil {
			return err
		}

		if err := bc.WriteTxWithNumber(nextTxNumber, tx); err != nil {
			return err
		}

		if err := bc.WriteLastTx(tx); err != nil {
			return err
		}

		if err := bc.WriteLastTxNumber(nextTxNumber); err != nil {
			return err
		}
	}
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "🔗 add new block",
		"hash", b.GetHash(types.BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
