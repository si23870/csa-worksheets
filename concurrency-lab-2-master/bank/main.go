package main

import (
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var debug *bool
var channywanny = make(chan transaction)
var woohoo = make(chan bool, 3000)

// executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, done chan<- bool) {
	for {
		t := <-channywanny

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		//mutex.Lock()
		//bank.lockAccount(t.from, strconv.Itoa(executorId))
		//bank.lockAccount(t.to, strconv.Itoa(executorId))

		fmt.Println("Executor\t", executorId, "attempting transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		// add locking semaphore
		fmt.Println("Executor\t", executorId, "locked account", from)
		fmt.Println("Executor\t", executorId, "locked account", to)
		//mutex.Unlock()

		bank.execute(t, executorId)

		bank.unlockAccount(t.from, "0")
		fmt.Println("Executor\t", executorId, "unlocked account", from)
		bank.unlockAccount(t.to, "0")
		fmt.Println("Executor\t", executorId, "unlocked account", to)

		bank.removeCompleted(e, executorId) // Removing this line will break visualisations.
		fmt.Println("OMG  ", executorId)
		done <- true
		woohoo <- true
		fmt.Println("Hello?? ", executorId)
	}
}

func meowhead(bank *bank, transactionQueue chan transaction) {
	var waitingRoom = make(chan transaction, 2000)
	var tNum int
	for t := range transactionQueue {
		if bank.accounts[t.from].locked || bank.accounts[t.to].locked {
			fmt.Println("-----------sent to WAITING ROOM")
			transactionQueue <- t
			fmt.Println("waiting room size: ", len(transactionQueue))

		} else {
			fmt.Println("------SCHEDULED")
			bank.lockAccount(t.from, "0")
			bank.lockAccount(t.to, "0")
			fmt.Println("AHH")
			channywanny <- t
			fmt.Println("AHHHHHHHHHHHHHHHHHHHHHH")
			tNum++
			if tNum >= 3 {
				for tNum > 0 {
					<-woohoo
					tNum--
				}
			}
		}

	}
	fmt.Println("FINISHED MAIN CHAN")
	for t := range waitingRoom {
		if bank.accounts[t.from].locked || bank.accounts[t.to].locked {
			waitingRoom <- t
		} else {
			channywanny <- t
			bank.lockAccount(t.from, "0")
			bank.lockAccount(t.to, "0")
		}
	}

}

func toChar(i int) rune {
	return rune('A' + i)
}

// main creates a bank and executors that will be handling the incoming transactions.
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	debug = flag.Bool("debug", false, "generate DOT graphs of the state of the bank")
	flag.Parse()

	bankSize := 6 // Must be even for correct visualisation.
	transactions := 1000

	accounts := make([]*account, bankSize)
	for i := range accounts {
		accounts[i] = &account{name: string(toChar(i)), balance: 1000}
	}

	bank := bank{
		accounts:               accounts,
		transactionsInProgress: list.New(),
		gen:                    newGenerator(),
	}

	startSum := bank.sum()

	transactionQueue := make(chan transaction, transactions)
	expectedMoneyTransferred := 0
	for i := 0; i < transactions; i++ {
		t := bank.getTransaction()
		expectedMoneyTransferred += t.amount
		transactionQueue <- t
	}

	done := make(chan bool)

	go meowhead(&bank, transactionQueue)
	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, done)
	}

	for total := 0; total < transactions; total++ {
		fmt.Println("Completed transactions\t", total)
		<-done
	}

	fmt.Println()
	fmt.Println("Expected transferred", expectedMoneyTransferred)
	fmt.Println("Actual transferred", bank.moneyTransferred)
	fmt.Println("Expected sum", startSum)
	fmt.Println("Actual sum", bank.sum())
	if bank.sum() != startSum {
		panic("sum of the account balances does not much the starting sum")
	} else if len(transactionQueue) > 0 {
		panic("not all transactions have been executed")
	} else if bank.moneyTransferred != expectedMoneyTransferred {
		panic("incorrect amount of money was transferred")
	} else {
		fmt.Println("The bank works!")
	}
}
