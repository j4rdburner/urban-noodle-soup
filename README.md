# SETUP

```
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-file.json"
```

You will also need to edit the `main.go` file and replace the `projectID` variable on line 95 with the correct value.

---

# RUNNING THE SERVER

```
go run main.go
```

---

# QUERYING THE SERVER

```
curl "http://localhost:8080/checkRequirement?address=0xYourWalletAddress&threshold=4&modifier=1"
```

This assumes your `PORT` env var is either set to `8080` or unspecified.

---

# EXAMPLE RESPONSE

```
{
    "address": "0x1234567890abcdef1234567890abcdef12345678",
    "passed": false
}

```

I am too broke to continue :(
