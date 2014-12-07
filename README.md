# Interoper

**Interoper** is a simple interop testing tool for Telehash. It uses docker under the hood and makes json reports.


## Test file format

Test files are Markdown files which start with an (indent formatted) code block of some JSON followed by a description of the test.

The JSON header can have the following entries:

* `"timeout"` defines the maximum runtime for a test. (examples: `"3ms"`, `"3s"`, `"3m"`, `"3h"`)
* `"containers"` defines container settings for the individual roles. It must be an object. Each entry can have the `"command"` entry.

```markdown
    {
      "timeout": "3m"
    }

# `net-link` Test basic link establishment.

## Scenario

The SUT must start an endpoint and write its keys and paths to `/shared/id_sut.json`.

The Driver must start an endpoint, load the keys and paths from `/shared/id_sut.json` and establish a link with the SUT. The driver must close the link after 2.5 minutes.

## Failure conditions

* The link fails to open (after 1 minute)
* The link breaks before it is closed.

## Success conditions

* The remains open for at least 2.5 minutes (It keeps the exchange open)
* The link closes cleanly
```


## Test Control Protocol.

Test processes can use regular logging to save informational data. In addition to the regullar logging processes can also send (and receive) JSON formatted commands over STDIN/STDOUT/STDERR.

All process types must send a `{"cmd":"ready"}` when they are ready with there setup and can start running the test. The Driver will then run its test and when it is done it send the `{"cmd":"done"}` command. The done command is forwarded to the SUT (so that it can shut itself down).
