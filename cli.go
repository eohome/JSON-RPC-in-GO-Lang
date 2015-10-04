package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"./ds"
	"os"
)

type Args struct {
	Ssp string
	Budget float64
}

type Args2 struct {
	Id int
}


func main() {

	client, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	var reply ds.Buyresponse
	var reply2 ds.Buyresponse
	var stockSymbolAndPercentage string
	var budget float64
	var ch int
	var tID int
	for;;{
		fmt.Println("MENU:\n1)Buy Stock\n2)Check Portfolio\n3)Exit")
		fmt.Println("Choose wisely:")
		fmt.Scanf("%d",&ch)
		switch(ch){

			case 1:{
					fmt.Println("Enter Stock Symbol and Percentage:")
					fmt.Scanf("%s",&stockSymbolAndPercentage)
					fmt.Println("Enter budget")
					fmt.Scanf("%f",&budget)

					args:= &Args{stockSymbolAndPercentage,budget}
	//fmt.Println(args.Ssp,args.Budget)
	
					c := jsonrpc.NewClient(client)

					err = c.Call("Stock.Buy", args, &reply)
					if err != nil {
						log.Fatal("RPC error:", err)
					}

					fmt.Println("ID              : ",reply.TradeId)
					fmt.Println("Stocks          : ",reply.Stocks)
					fmt.Println("Unvested Amount : ",reply.UnvestedAmount)
					break;
					}
			case 2:{
					fmt.Println("ID: ")
					fmt.Scanf("%d",&tID)
					args2:= &Args2{tID}
					c := jsonrpc.NewClient(client)

					err = c.Call("Stock.CheckPortfolio", args2, &reply2)
					if err != nil {
						log.Fatal("RPC error:", err)
					}
					fmt.Println("ID              : ",reply2.TradeId)
					fmt.Println("Stocks          : ",reply2.Stocks)
					fmt.Println("CMV             : ",reply2.CurrentMarketValue)
					fmt.Println("Unvested Amount : ",reply2.UnvestedAmount)
					break;
					}
			case 3:{
					os.Exit(0)
					break;
					}
	}
	}
}
