package main

import (
	"log"
	"net"
	"net/rpc"
	"fmt"
	"net/rpc/jsonrpc"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
	"strconv"
	"strings"
	"./ds"
)


type Args struct {
	Ssp string
	Budget float64
}

type clientPortfolio struct{
	tradeId int
	compName []string
	percentStock []float64
	amountStock []float64
	investedAmount []float64
	purchaseAmount []float64
	numberOfStocks []int
	budget float64
	unvestedAmount float64
}


type yqlresponse struct{
	
Query struct{
	Count int `json:"count"`
	Created string `json:"created"`
	Lang string `json:"lang"`
	Results struct{
	Quote []struct{
	Ask string `json:"Ask"`
	Symbol string `json:"symbol"`
	}`json:"quote"`
	}`json:"results`
	}`json:"query"`
}

type Buying struct{
	name string
	askingamount float64
}


var client []clientPortfolio

type Stock struct{}
var id int
func initS(){
id=0
}

func parseSsp(args *Args) ([]string, []float64,[]float64) {
	ssp:=strings.Split(args.Ssp,",")
	compName :=make([]string,len(ssp))
	percentStock:=make([]float64,len(ssp))
	amountStock:=make([]float64,len(ssp))
	
	for i:=range ssp{
		details:=strings.Split(ssp[i],":")
		compName[i] = details[0]
		percentStock[i] , _=strconv.ParseFloat(strings.TrimSuffix(details[1],"%"),64)
		amountStock[i] = args.Budget*percentStock[i]/100
	}
	return compName,percentStock,amountStock
}


func (t *Stock) Buy(args *Args, reply *ds.Buyresponse) error {
	id++;

	ns:=0
	totalunvestedAmount:=0.00
	var uv float64

	compName,percentStock,amountStock := parseSsp(args);
	purchase:=getYahooAPI(compName);
	purchaseAmount:=make([]float64,len(amountStock))
	stockstring :=make([]string,len(amountStock))
	investedAmount:=make([]float64,len(amountStock))
	numberOfStocks:=make([]int,len(amountStock))
	for i := range amountStock{

		ns=0
		uv=0.00
		investedAmount[i]=0.0
		uv=amountStock[i]
		fmt.Println(purchase[i].askingamount)
		for uv>purchase[i].askingamount{
			uv=(uv-purchase[i].askingamount)
			ns++
		}
		investedAmount[i]=amountStock[i]-uv
		totalunvestedAmount=totalunvestedAmount+uv
		numberOfStocks[i]=ns
		purchaseAmount[i]=purchase[i].askingamount
		stockstring[i] = compName[i]+":"+strconv.Itoa(numberOfStocks[i])+":$"+strconv.FormatFloat(purchase[i].askingamount, 'f', 6, 64)


	}

	client = append(client,clientPortfolio{id,compName,percentStock,amountStock,investedAmount,purchaseAmount,numberOfStocks,args.Budget,totalunvestedAmount})


	
//GENERATING REPLY !!
	reply.TradeId =id

	reply.Stocks=strings.Join(stockstring,",")

	reply.UnvestedAmount=totalunvestedAmount

	
	return nil
}


func getYahooAPI(compName []string) []Buying {
	
	var yqlres yqlresponse;
	//var searchstring string
	searchstring:=strings.Join(compName,"%22%2C%22")
	resp, err := http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20symbol%2CAsk%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"+searchstring+"%22)&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env&callback=")
	 if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }else {

        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }

        result:=string(contents)


        err2 := json.Unmarshal([]byte(result), &yqlres)
    if err2 == nil {
        //fmt.Printf("%+v\n", yqlres.Query.Results.Quote[0].Ask)
        // fmt.Println("DONE1")
        var purchase = make([]Buying,yqlres.Query.Count)

        for i:=range yqlres.Query.Results.Quote{
        	purchase[i].name=yqlres.Query.Results.Quote[i].Symbol
        	purchase[i].askingamount, _=strconv.ParseFloat(yqlres.Query.Results.Quote[i].Ask,64)
        }
        return purchase;
    } else {
        fmt.Println(err)
        //fmt.Printf("%+v\n", yqlres.Query.Results.Quote[0].Ask)
    }
    }
	
	return nil
}

func getCurrentVal(id int) []Buying {
	
	var yqlres yqlresponse;
	//var searchstring string
	searchstring:=strings.Join(client[id-1].compName,"%22%2C%22")
	resp, err := http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20symbol%2CAsk%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"+searchstring+"%22)&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env&callback=")
	 if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }else {

        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }

        result:=string(contents)


        err2 := json.Unmarshal([]byte(result), &yqlres)
    if err2 == nil {
        //fmt.Printf("%+v\n", yqlres.Query.Results.Quote[0].Ask)
        // fmt.Println("DONE1")
        var purchase = make([]Buying,yqlres.Query.Count)

        for i:=range yqlres.Query.Results.Quote{
        	purchase[i].name=yqlres.Query.Results.Quote[i].Symbol
        	purchase[i].askingamount, _=strconv.ParseFloat(yqlres.Query.Results.Quote[i].Ask,64)
        }
        return purchase;
    } else {
        fmt.Println(err)
        //fmt.Printf("%+v\n", yqlres.Query.Results.Quote[0].Ask)
    }
    }
	
	return nil
}

type Args2 struct{
	Id int
}
func (t *Stock) CheckPortfolio(args *Args2,reply *ds.Buyresponse)error{
	fmt.Println(args.Id)
	CurrentMarketValue :=0.00
	//searchstring:=strings.Join(client[args.Id].compName,"%22%2C%22")
	purchase:=getCurrentVal(args.Id);
	fmt.Println(purchase)
	stockstring :=make([]string,len(purchase))
	 for i :=range purchase{
	 	fmt.Println(client[args.Id-1].purchaseAmount[i])
	 	if purchase[i].askingamount>client[args.Id-1].purchaseAmount[i]{
	 	stockstring[i] = client[args.Id-1].compName[i]+":"+strconv.Itoa(client[args.Id-1].numberOfStocks[i])+":+$"+strconv.FormatFloat(purchase[i].askingamount, 'f', 6, 64)
	 }else if purchase[i].askingamount<client[args.Id-1].purchaseAmount[i]{
	 	stockstring[i] = client[args.Id-1].compName[i]+":"+strconv.Itoa(client[args.Id-1].numberOfStocks[i])+":-$"+strconv.FormatFloat(purchase[i].askingamount, 'f', 6, 64)
	 }else{
	 	stockstring[i] = client[args.Id-1].compName[i]+":"+strconv.Itoa(client[args.Id-1].numberOfStocks[i])+":$"+strconv.FormatFloat(purchase[i].askingamount, 'f', 6, 64)
	 }
	 CurrentMarketValue=CurrentMarketValue+float64(client[args.Id-1].numberOfStocks[i])*purchase[i].askingamount
	 }

	reply.TradeId =args.Id
	reply.Stocks=strings.Join(stockstring,",")
	reply.CurrentMarketValue=CurrentMarketValue
	reply.UnvestedAmount=client[args.Id-1].unvestedAmount

	return nil
}

func main() {
	initS()
	stock := new(Stock)
	server := rpc.NewServer()
	server.Register(stock)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}

