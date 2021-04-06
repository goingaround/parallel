_parallel.Operation_
An object built to run multiple anonymous functions simultaneously with relative delay (or exactly on the same time for delay=0). 

thread safe and reuseable.


`NewOperation(funcs ...func()) (*operation, error)` - Initializes a new _Operation_

-  _funcs_ - the functions of the operation

returns error on 
 - given no functions
 - sent nil function

`(op *operation) Run(timeout, delay time.Duration) error` - start running the operation
- _timeout_ - the max duration the operation can run (timeout <= 0 is no timeout)
- _delay_ - transitive delay duration between process execution (delay <= 0 is no delay)

returns error on
 - timeout exceeded
