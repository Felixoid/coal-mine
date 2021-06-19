# coal-mine
Carbon metrics generator for graphite-web system.

# Installation
To get the working binary run the command  
`go install github.com/Felixoid/coal-mine`

# How to use
The program accepts multiple `--const`, `--counter` and `--random` arguments as metrics generator names. They can be specified as curly brace expandable masks, for example `server{01..10}.soft{1..5}` will generate 50 metrics with two nodes. `--from` and `--until` accept the same values as graphite-web `/render` handler. `--value` and `--deviation` values affect the each next point for the metrics.

Run `coal-mine config-example` to see the full explanation of each generator type.

`coal-mine --random '1.{001..100}.3.4{22..225}' --random '1.{101..200}.3.4{22..225}' --random '1.{201..300}.3.4{22..225}' --random '1.{301..400}.3.4{22..225}' --counter '22.22.33.{10..100}' --from -2d --until 23h --step 300`

Additionally, the generators can be set through the configuration file with `-c/--config config.toml` argument. Then each generator can have custom `from/until/step/value/deviation` parameters.

## Simulate on-time metrics sending
To mock the normal metrics sending, for example, to perform the load test, the program has a special mode:  
`coal-mine online --random '1.{001..00}.3.4{22..25}' --step 3 --randomize`  
It's highly recommended to use `--randomize` in the online mode to send metrics each second.

On `Ctrl+C` it will finish the current writes and then exits.
