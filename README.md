# olleth-palinda-3

## Task 1 Answers

1. What happens if you remove the go-command from the Seek call in the main function? 
> If you remove the go-command from the Seek call in the main function, the program will execute Seek in the main goroutine instead of creating a new goroutine for each Seek call. This will result in a sequential execution of Seek calls.

2. What happens if you switch the declaration wg := new(sync.WaitGroup) to var wg sync.WaitGroup and the parameter wg *sync.WaitGroup to wg sync.WaitGroup?
> If you switch the declaration wg := new(sync.WaitGroup) to var wg sync.WaitGroup and the parameter wg *sync.WaitGroup to wg sync.WaitGroup, the program will not compile because the WaitGroup needs to be passed by reference.

3. What happens if you remove the buffer on the channel match?
> If you remove the buffer on the channel match, the program will deadlock because the Send operation in Seek will block indefinitely until a receiver is ready to receive from the channel.

4. What happens if you remove the default-case from the case-statement in the main function? 
> If you remove the default-case from the case-statement in the main function, the program will not print anything if no message is received on the match channel. It will just exit without any notification.

## Task 3 Runtime performance table

| Variant | Runtime (ms) |
| ------- | ------------ |
| singleworker | 7.52 |
| mapreduce | 5.23 |
