# help

```
oc login https://console.abasky.net --token=<token>
docker login -u <user> -p $(oc whoami -t) https://registry.app.abasky.net
docker build . -t micromdm:testing
docker tag micromdm:testing registry.app.abasky.net/abaclock-mdm-dev/micromdm
docker push registry.app.abasky.net/abaclock-mdm-dev/micromdm
```
