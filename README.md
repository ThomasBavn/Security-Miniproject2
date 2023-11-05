# Security miniproject 2

## How to run

Open 4 terminals 

```bash
# Terminal 1, this is the hospital
go run . 0

# Terminal 2, this is Alice
go run . 1

# Terminal 3, this is Bob
go run . 2

# Terminal 4, this is Charlie
go run . 3
```

## What you will see

After running the program, each of the non-hospital nodes will send their own port as their secret
to the hospital. You will be able to see in the hospital terminal that the hospital 
gets the final value of `15006` which is `5001 + 5002 + 5003` which is the sum of all the ports.

