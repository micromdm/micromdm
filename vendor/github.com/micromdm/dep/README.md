dep is a client library for working with Apple's Device Enrollment Program

See [godoc](https://godoc.org/github.com/micromdm/dep) for detailed usage.


# Usage

Configure and create an http client passing the Oauth credentials from the server token.
```
    config := dep.Config{
        ConsumerKey:    "CK_3a419c0b",
        ConsumerSecret: "CS_3fb23281",
        AccessToken:    "AT_O8473841",
        AccessSecret:   "AS_9d141598",
    }
    client, err := dep.NewClient(config)
    if err != nil {
        // handle err
    }
```

Use the new DEP client:
```
    account, err := client.Account()
    if err != nil {
        // handle err
    }
```

# Example

In the examples folder, there's an example that you can try running against the `depsim` binary.
