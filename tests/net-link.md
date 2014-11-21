    {
      "timeout": "10s",
      "sut":    { "command": "test-net-link await" },
      "driver": { "command": "test-net-link establish" }
    }

# `net-link` Test basic link establishment.

## System Under Test

The SUT must start an endpoint and write its keys and paths to `/shared/sut.json`

## Driver

The Driver must start an endpoint, load the keys and paths from `/shared/sut.json` and establish a link with the other endpoint.
