    {
      "timeout": "10s"
    }

# `sanity` Test the basic interoper features

## Scenario

The worker must signal that it is ready by writing `{"ty":"ready"}` to STDOUT.

The Driver must also signal that it is ready by writing `{"ty":"ready"}` to STDOUT. Then it is expected to sleep 5 seconds after which it must signal that it is done by writing the `{"ty":"done"}` command to STDOUT.

## Failure conditions

* The worker fails to start.
* The Driver fails to start.
* The worker fails signal that it is ready.
* The Driver fails signal that it is ready.
* The Driver fails signal that it is done.
* The worker fails wait for the `done` signal.

## Success conditions

* The worker must exit after the Driver.
* Both the Driver and the worker must exit before the timeout.
