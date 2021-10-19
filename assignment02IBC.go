package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

//-------------------------------------------------------------------------------------------------------
//								FUNCTIONS TO BUILD BLOCKCHAIN
//-------------------------------------------------------------------------------------------------------
func CalculateBalance(userName string, chainHead *Block) int {
	tempchainhead := chainHead
	balance := 0
	for (tempchainhead != nil){
		for i := range tempchainhead.Data{
			if tempchainhead.Data[i].Sender == userName{
				balance -= tempchainhead.Data[i].Amount
			}
			
			if tempchainhead.Data[i].Receiver == userName{
				balance += tempchainhead.Data[i].Amount
			}
		}
		tempchainhead = tempchainhead.PrevPointer
	}
	return balance
}

func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
		storecurrbal := CalculateBalance(transaction.Sender, chainHead)  
		if ( storecurrbal < transaction.Amount){
			fmt.Println("INVALID TRANSACTION -- " + transaction.Sender + " has insufficient balance! [" + fmt.Sprint(transaction.Amount) + " needed - " + fmt.Sprint(storecurrbal) + " available]")
			return false
		}
		return true
}
func VerifyTransactionss(transactions [] BlockData, chainhead *Block) bool{
	for i := range transactions{
		if (transactions[i].Title == "Coinbase"){
			return true
		}
		if !(VerifyTransaction(&transactions[i], chainhead)){
			return false
		}
	}
	return true
}
func PremineChain(chainHead *Block, numBlocks int) *Block {
	premine := []BlockData{{Title: "Coinbase", Sender: "System", Receiver: rootUser, Amount: miningReward*numBlocks}}
	chainHead = InsertBlock(premine, chainHead)
	return chainHead
}

func checkwithintransactions(chainHead *Block) bool{
	var tempchain *Block = chainHead
	for {
		if(tempchain == nil){
			break
		}
		for i := range tempchain.Data{
			if tempchain.Data[i].Sender != "System"{
				if CalculateBalance( tempchain.Data[i].Sender,tempchain) < 0{
					//fmt.Println("****ATTENTION INSUFFICIENT FUNDS --> BLOCK INVALID****")
					return false
				}
			}
		}
		tempchain = tempchain.PrevPointer
	}
	return true
}

//Function to insert a new block in the block chain
func InsertBlock(dataToInsert []BlockData, chainHead *Block) *Block {

	if !(VerifyTransactionss(dataToInsert, chainHead)){
		return chainHead;
	}

	if chainHead == nil {		//in case first block being inserted
		var tempblock Block
		chainHead = &tempblock
		chainHead.Data = dataToInsert
		chainHead.PrevPointer = nil
		chainHead.CurrentHash = CalculateHash(chainHead)
		chainHead.PrevHash = ""
	} else {
		miningrewardtoRootUser := BlockData{Title: "MiningReward", Sender: "System", Receiver: rootUser, Amount: miningReward}
		dataToInsert = append(dataToInsert, miningrewardtoRootUser)
		var tempblock Block
		tempblock.Data = dataToInsert
		tempblock.PrevPointer = chainHead
		tempblock.PrevHash = chainHead.CurrentHash
		tempblock.CurrentHash = CalculateHash(&tempblock)
		if !checkwithintransactions(&tempblock){
			return chainHead
		}
		chainHead = &tempblock
	}

	return chainHead
}

func CalculateHash(inputBlock *Block) string {
	return fmt.Sprintf("%x",sha256.Sum256([]byte(fmt.Sprintf("%v",*inputBlock))))
}


func ListBlocks(chainHead *Block) {
	tempchainhead := chainHead
	fmt.Println("LIST OF BLOCKS IN CHAIN:")
	for {
		if (tempchainhead == nil){
			break
		}
		fmt.Print("HASH: ")
		fmt.Println(tempchainhead.CurrentHash)
		fmt.Print("TRANSACTIONS: \n")
		for i := range tempchainhead.Data{
			fmt.Println("Title: " + tempchainhead.Data[i].Title + "- Sender: " + tempchainhead.Data[i].Sender + " - Reciever: " + tempchainhead.Data[i].Receiver + " - Amount: " + fmt.Sprint(tempchainhead.Data[i].Amount))
		}
		fmt.Print("PREV HASH: ")
		fmt.Printf(tempchainhead.PrevHash)
		fmt.Print("\n----------------------------------------\n")
		tempchainhead = tempchainhead.PrevPointer
	}	
	print("\n\n")
}	

func VerifyChain(chainHead *Block) {
	tempchainhead := chainHead
	for {
		if ( tempchainhead == nil || tempchainhead.PrevPointer == nil){
			break
		}
		prevhashstored := tempchainhead.PrevHash
		tempchainhead = tempchainhead.PrevPointer
		currhash := tempchainhead.CurrentHash
		if(prevhashstored != currhash){
			fmt.Println("****ATTENTION HASHES DO NOT MATCH --> BLOCK INVALID****")
			return
		}
		for i := range tempchainhead.Data{
			if CalculateBalance( tempchainhead.Data[i].Sender,chainHead) < 0{
				fmt.Println("****ATTENTION HASHES DO NOT MATCH --> BLOCK INVALID****")
				return
			}
		}
	}
	fmt.Println("****BLOCK IS VALID! YAY :D ****")	
}
