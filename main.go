package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"

	elastic "github.com/olivere/elastic/v7"
	"github.com/tjandrayana/ee/engine"
	secure "github.com/tjandrayana/searchable-encryption/transform"
)

func main() {

	content, err := ioutil.ReadFile("./dataset/users.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var users []engine.User

	json.Unmarshal(content, &users)

	// fmt.Printf("Data : %v\n", users)

	var newUsers []engine.User

	for _, d := range users {
		d.SecureAddress = secure.Transform(d.Address)
		newUsers = append(newUsers, d)
	}

	DefaultURL := "http://localhost:9200"
	client, err := elastic.NewClient(elastic.SetURL(DefaultURL))

	ctx := context.Background()

	for i, d := range newUsers {
		data, _ := json.Marshal(d)

		put1, err := client.Index().
			Index("test-encryption").
			Type("_doc").
			Id(fmt.Sprintf("%d", (i+1)*-1)).
			BodyJson(string(data)).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}

	search := "jalan"

	fmt.Printf("search : %s\n", secure.Transform(search))

	// webhooks.Run()

}
