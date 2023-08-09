package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"encoding/json"

	elastic "github.com/olivere/elastic/v7"
	"github.com/tjandrayana/ee/engine"
	"github.com/tjandrayana/ee/secure"
	"github.com/tjandrayana/searchable-encryption/transform"
)

var es = true

func main() {
	key := []byte("0123456789abcdef0123456789abcdef") // 256-bit key

	if es {

		content, err := ioutil.ReadFile("./dataset/users.json")
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}

		var users []engine.User

		json.Unmarshal(content, &users)

		var newUsers []engine.SecureUser

		for _, d := range users {

			name := strings.ToLower(d.Name)
			email := strings.ToLower(d.Email)
			address := strings.ToLower(d.Address)
			phone := strings.ToLower(d.PhoneNumber)

			cipherName, _ := secure.EncryptAES([]byte(name), key)
			if string(cipherName) != "" {
				name = string(cipherName)
			}

			cipherEmail, _ := secure.EncryptAES([]byte(email), key)
			if string(cipherEmail) != "" {
				email = string(cipherEmail)
			}

			cipherAddress, _ := secure.EncryptAES([]byte(address), key)
			if string(cipherAddress) != "" {
				address = string(cipherAddress)
			}

			cipherPhoneNumber, _ := secure.EncryptAES([]byte(phone), key)
			if string(cipherPhoneNumber) != "" {
				phone = string(cipherPhoneNumber)
			}

			newData := engine.SecureUser{
				Name:              name,
				Email:             email,
				Address:           address,
				PhoneNumber:       phone,
				SecureName:        strings.ToLower(transform.Transform(d.Name)),
				SecureEmail:       strings.ToLower(transform.Transform(d.Email)),
				SecureAddress:     strings.ToLower(transform.Transform(d.Address)),
				SecurePhoneNumber: strings.ToLower(transform.Transform(d.PhoneNumber)),
			}

			newUsers = append(newUsers, newData)
		}

		DefaultURL := "http://10.21.32.5:9200"
		client, _ := elastic.NewClient(elastic.SetURL(DefaultURL))

		ctx := context.Background()

		for i, d := range newUsers {
			data, _ := json.Marshal(d)

			put1, err := client.Index().
				Index("test-encryption-v1").
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
	} else {

		// fmt.Printf("search : %s\n", secure.CustomEncrypt("Budi Santoso", customKey))
		// fmt.Printf("search : %s\n", secure.CustomEncrypt("Gajah Mada", customKey))
		// fmt.Printf("search : %s\n", secure.CustomEncrypt("Gajah Moda", customKey))

		fmt.Printf("search : %s\n", strings.ToLower(transform.Transform("Budi")))
		fmt.Printf("search : %s\n", strings.ToLower(transform.Transform("Gajah Mada")))
		fmt.Printf("search : %s\n", strings.ToLower(transform.Transform("Gajah")))
		fmt.Printf("search : %s\n", strings.ToLower(transform.Transform("Santo")))

	}

	// webhooks.Run()

}
