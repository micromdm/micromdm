// Package dep is a client library for Apple's Device Enrollment Program
/*
Configure and create an http client passing the Oauth credentials from the server token.
	config := dep.Config{
		ConsumerKey:    "CK_3a419c0b",
		ConsumerSecret: "CS_3fb23281",
		AccessToken:    "AT_O8473841",
		AccessSecret:   "AS_9d141598",
	}
	client := dep.NewClient(config)

Use the new DEP client:
	account, err := client.Account()
	if err != nil {
		// handle err
	}

*/
package dep
