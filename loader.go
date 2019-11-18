package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rqlite/gorqlite"
)

var (
	addr   = "http://rbox1:4001"
	leader string
	peers  []string
	trace  bool
)

func init() {
	flag.StringVar(&addr, "addr", addr, "rqlite leader address")
	flag.BoolVar(&trace, "trace", false, "debug trace to stderr")
}

func main() {
	flag.Parse()
	if trace {
		gorqlite.TraceOn(os.Stderr)
	}

	conn, err := gorqlite.Open(addr)
	if err != nil {
		panic(err)
	}
	leader, err = conn.Leader()
	if err != nil {
		panic(err)
	}
	fmt.Printf("OUR FEARLESS LEADER: %s\n", leader)

	peers, err = conn.Peers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("WE HAVE %d PEERS: %+v\n", len(peers), peers)
	for _, peer := range peers {
		log.Println("peer:", peer)
	}

	statements := []string{
		"INSERT INTO fakery (name) VALUES ('joe bob')",
	}
	for i := 0; i < 1000; i++ {
		results, err := conn.Write(statements)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, result := range results {
			if result.Err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("Last rowid:", result.LastInsertID)
		}
		time.Sleep(time.Second * 1)
	}
}

// these URLs are just generic database URLs, not rqlite API URLs,
// so you don't need to worry about the various rqlite paths ("/db/query"), etc.
// just supply the base url and not "db" or anything after it.

// yes, you need the http or https

// no, you cannot specify a database name in the URL (this is sqlite, after all).

/*
conn, err := gorqlite.Open("http://") // connects to localhost on 4001 without auth
conn, err := gorqlite.Open("https://") // same but with https
conn, err := gorqlite.Open("https://localhost:4001/") // same only explicitly

// with auth:
conn, err := gorqlite.Open("https://mary:secret2@localhost:4001/")
// different server, setting the rqlite consistency level
conn, err := gorqlite.Open("https://mary:secret2@server1.example.com:4001/?level=none")
// same without auth, setting the rqlite consistency level
conn, err := gorliqte.Open("https://server2.example.com:4001/?level=weak")
// different port, setting the rqlite consistency level and timeout
conn, err := gorqlite.Open("https://localhost:2265/?level=strong&timeout=30")

// change our minds
conn.SetConsistencyLevel("none")
conn.SetConsistencyLevel("weak")
conn.SetConsistencyLevel("strong")

// set the http timeout.  Note that rqlite has various internal timeouts, but this
// timeout applies to the http.Client and its work.  It is measured in seconds.
conn.SetTimeout(10)

// simulate database/sql Prepare()
statements := make ([]string,0)
pattern := "INSERT INTO secret_agents(id, hero_name, abbrev) VALUES (%d, '%s', '%3s')"
statements = append(statements,fmt.Sprintf(pattern,125718,"Speed Gibson","Speed"))
statements = append(statements,fmt.Sprintf(pattern,209166,"Clint Barlow","Clint"))
statements = append(statements,fmt.Sprintf(pattern,44107,"Barney Dunlap","Barney"))
results, err := conn.Write(statements)

// now we have an array of []WriteResult

for n, v := range WriteResult {
	fmt.Printf("for result %d, %d rows were affected\n",n,v.RowsAffected)
	if ( v.Err != nil ) {
		fmt.Printf("   we have this error: %s\n",v.Err.Error())
	}
}

// or if we have an auto_increment column
res, err := conn.WriteOne("INSERT INTO foo (name) values ('bar')")
fmt.Printf("last insert id was %d\n",res.LastInsertID)

// just like database/sql, you're required to Next() before any Scan() or Map()

// note that rqlite is only going to send JSON types - see the encoding/json docs
// which means all numbers are float64s.  gorqlite will convert to int64s for you
// because it is convenient but other formats you will have to handle yourself

var id int64
var name string
rows, err := conn.QueryOne("select id, name from secret_agents where id > 500")
fmt.Printf("query returned %d rows\n",rows.NumRows)
for rows.Next() {
	err := response.Scan(&id, &name)
	fmt.Printf("this is row number %d\n",response.RowNumber)
	fmt.Printf("there are %d rows overall%d\n",response.NumRows)
}

// just like WriteOne()/Write(), QueryOne() takes a single statement,
// while Query() takes a []string.  You'd only use Query() if you wanted
// to transactionally group a bunch of queries, and then you'd get back
// a []QueryResult

// alternatively, use Next()/Map()

for rows.Next() {
	m, err := response.Map()
	// m is now a map[column name as string]interface{}
	id := m["name"].(float64) // the only json number type
	name := m["name"].(string)
}

// get rqlite cluster information
leader, err := conn.Leader()
// err could be set if the cluster wasn't answering, etc.
fmt.Println("current leader is"leader)
peers, err := conn.Peers()
for n, p := range peers {
	fmt.Printf("cluster peer %d: %s\n",n,p)
}

// turn on debug tracing to the io.Writer of your choice.
// gorqlite will verbosely write very granular debug information.
// this is similar to perl's DBI->Trace() facility.
// note that this is done at the package level, not the connection
// level, so you can debug Open() etc. if need be.

f, err := os.OpenFile("/tmp/deep_insights.log",OS_RDWR|os.O_CREATE|os.O_APPEND,0644)
gorqlite.TraceOn(f)

// change my mind and watch the trace
gorqlite.TraceOn(os.Stderr)

// turn off
gorqlite.TraceOff()
*/
