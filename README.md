## Go Json RPC

### Client
Curl
```bash
$ curl -X POST http://localhost:8080/rpc -H 'cache-control: no-cache' -H 'content-type: application/json' -d '{"params":{"Name":"go-me","Age":33},"id":"12","jsonrpc":"2.0","method":"HelloService.Hello"}'

$ curl -X POST http://localhost:8080/rpc -H 'cache-control: no-cache' -H 'content-type: application/json' -d '{"params":{},"id":"12","jsonrpc":"2.0","method":"TimeService.GetTime"}'

$ curl -X POST http://localhost:8080/rpc -H 'cache-control: no-cache' -H 'content-type: application/json' -d '{"id":"12","jsonrpc":"2.0","method":"TimeService.GetTime"}'
```

```powershell
$data=@{
    params = @{
        Name = "go-go"
        Age = 33
    }
    
    id = "1"
    jsonrpc = "2.0"
    method = "HelloService.Hello"
}

$data | ConvertTo-Json -Compress | curl.exe -X POST http://localhost:8080/rpc -H 'cache-control: no-cache' -H 'content-type: application/json' -d "@-"
```