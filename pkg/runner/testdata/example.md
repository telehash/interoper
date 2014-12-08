    {
      "timeout": "1m",
      "worker": { "command": "test-net-link await" },
      "driver": { "command": "test-net-link establish" }
    }

# `net-link` Test basic link establishment.

# Worker

The Worker must start an endpoint and write its keys and paths to `/shared/endpoint-a.json`

The Driver must start an endpoint, load the keys and paths from `/shared/endpoint-a.json` and establish a link with the other endpoint.
