# Tokenize!

Have some data you need encrypted and need to get a token back? Then do I have the service for you!

## Purpose
Sometimes you need to encrypt some data and get a token to refer to that data in a unique way. The encrypted data should only be available to certain users (RBAC)
but you may have attributes you want to know about that stored data. This tokenize service lets a user take a payload, such as sensitive data, and store and encrypt it.

## How to use

This service has some pretty simple endpoints:

### POST /token
Input your values and get a token back. You can set the TTL, metadata about the token, send the payload to be encrypted, and a type of token.

```
POST /token
{
  "data": {
    "metadata": {
      "bin": "41111111",
      "last4": "1111" 
    },
    "payload": "{'card_number': '4111111111111111','exp': '0128'}",
    "token_type": "card",
    "ttl": 7200
  }
}
```

### GET /token/{token}
This will return the token properties without the payload.

### GET /token/{token}/decrypt
This will return the token properties with the payload decrypted. 

### POST /token/{token}
Update the metadata and TTL of the token

```
POST /token/{token}
{
  "metadata": {
    "foo": "bar"
  },
  "ttl": 600
}
```

### DELETE /token/{token}
This will delete a token

## To Do:

- [ ] Update the service runner
- [ ] Dockerize and add Compose for dev
- [ ] Logging
- [ ] Test all the things
- [ ] Operationalize
- [x] Datastore
