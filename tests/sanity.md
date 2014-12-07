    {
      "timeout": "10s"
    }

# `sanity` Test the basic interoper features

## Scenario

The SUT must signal that it is ready by writing `{"cmd":"ready"}` to STDOUT. Then it must read the `{"cmd":"done"}` command from STDIN.

The Driver must also signal that it is ready by writing `{"cmd":"ready"}` to STDOUT. Then it is expected to sleep 5 seconds after which it must signal that it is done by writing the `{"cmd":"done"}` command to STDOUT.

## Failure conditions

* The SUT fails to start.
* The Driver fails to start.
* The SUT fails signal that it is ready.
* The Driver fails signal that it is ready.
* The Driver fails signal that it is done.
* The SUT fails wait for the `done` signal.

## Success conditions

* The SUT must exit after the Driver.
* Both the Driver and the SUT must exit before the timeout.
