package main

import (
	"fmt"
	"log"

	"github.com/micromdm/dep"
)

func main() {
	// Auth configuration
	config := &dep.Config{
		ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
		ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
		AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
		AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
	}

	// create an http client
	client, err := dep.NewClient(config, dep.ServerURL("http://localhost:9000"))
	if err != nil {
		log.Fatal(err)
	}

	// get account information
	account, err := client.Account()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.OrgAddress)

	// fetch devices
	fetchDevices, err := client.FetchDevices(dep.Limit(100))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fetchDevices)
}
