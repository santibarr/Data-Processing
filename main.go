package main

import (
	"errors"
	"fmt"
)

type DataBase struct {
	actualData       map[string]int
	transactionState bool
	transactionData  map[string]int
}

// Return a pointer to Database struct -> needs to create the maps
func createDB() *DataBase {
	return &DataBase{
		actualData:       make(map[string]int),
		transactionData:  make(map[string]int),
		transactionState: false,
	}

}

// function for beginning the database
func (db *DataBase) begin_transaction() error {
	if db.transactionState {
		return errors.New("ERR: transaction already in progress")
	}

	//set the state to to true
	db.transactionState = true
	db.transactionData = make(map[string]int)
	db.actualData = make(map[string]int)
	return nil
}

// put the key into the data base
func (db *DataBase) put(key string, value int) error {

	if !db.transactionState {
		return errors.New("ERR: transaction is not in progress")
	}
	db.transactionData[key] = value
	return nil
}

func (db *DataBase) get(key string) int {

	if val, ok := db.actualData[key]; ok {
		return val
	} else {
		return 0
	}
}

func (db *DataBase) commit() error {

	//if there is a transaction going on
	if db.transactionState {
		for k, v := range db.transactionData {

			db.actualData[k] = v
		}

		db.transactionState = false
		db.transactionData = make(map[string]int) //reset

		return nil
	}

	return errors.New("There is no open transaction")
}

func (db *DataBase) rollback() error {

	//reset everything
	if db.transactionState {

		db.transactionState = false
		db.transactionData = make(map[string]int)
		return nil
	}

	return errors.New("There is no open transaction")
}

func main() {

	//created the DB
	db := createDB()

	// should return null, because A doesn’t exist in the DB yet
	val := db.get("A")

	if val == 0 {
		fmt.Println("Could not get value")
	}

	// should throw an error because a transaction is not in progress
	err := db.put("A", 5)

	if err != nil {
		fmt.Println("Transaction is not in progress")
	}

	// starts a new transaction
	err = db.begin_transaction()

	if err != nil {
		fmt.Println("Transaction is already in progress")
	}

	// set’s value of A to 5, but its not committed yet
	err = db.put("A", 5)

	if err != nil {
		fmt.Println("Transaction is not in progress")
	}

	// should return null, because updates to A are not committed yet
	val = db.get("A")

	if val == 0 {
		fmt.Println("Value has not been committed")
	} else {
		fmt.Println("This is the value: ", val)
	}

	// update A’s value to 6 within the transaction
	err = db.put("A", 6)

	if err != nil {
		fmt.Println("Transaction is not in progress")
	} 

	// commits the open transaction
	db.commit()

	// should return 6, that was the last value of A to be committed
	val = db.get("A")

	if val == 0 {
		fmt.Println("Value has not been committed ")
	} else {
		fmt.Println("This is the value: ", val)
	}

	// throws an error, because there is no open transaction
	err = db.commit()

	if err != nil {
		fmt.Println("There is no open transaction")
	}

	// throws an error because there is no ongoing transaction
	err = db.rollback()

	if err != nil {
		fmt.Println("There is no open transaction")
	}

	// should return null because B does not exist in the database
	val = db.get("B")

	if err != nil {
		fmt.Println("Does not exist in DataBase")
	}

	//Restart everything again
	db.begin_transaction()

	// Set key B’s value to 10 within the transaction
	err = db.put("B", 10)

	if err != nil {
		fmt.Println("Transaction is not in progress")
	}

	err = db.rollback()

	// Rollback the transaction - revert any changes made to B
	if err == nil {
		fmt.Println("Changes have been reverted")
	}

	val = db.get("B")

	// Should return null because changes to B were rolled back
	if val == 0 {
		fmt.Println("Changes have been rolled back")
	}

}
