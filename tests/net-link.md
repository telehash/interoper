    {
      "timeout": "3m"
    }

# `net-link` Test basic link establishment.

## Scenario

The SUT must start an endpoint and write its keys and paths to `/shared/sut.json`.

The Driver must start an endpoint, load the keys and paths from `/shared/sut.json` and establish a link with the SUT. The driver must close the link after 2.5 minutes.

## Failure conditions

* The link fails to open (after 1 minute)
* The link breaks before it is closed.

## Success conditions

* The remains open for at least 2.5 minutes (It keeps the exchange open)
* The link close cleanly
