2ndroute is an ambassador that intelligently directs service calls.  Its primary use is to make connections between components
in server side applications more dynamic and address the birdnest that often results as applicaitons evolve over time.  This
helps keep applicaitons portable between environments by always binding to a localhost:port pattern.

It is a small command line tool that stores routes in consul, and uses iptables to masquerade traffic.

## Usage

    Usage:
    2ndroute [OPTIONS] <command>

    commands:
      list - show current routes
      add  - create new route
      rm   - remove route

    options
      -d   run in deamon mode watching for route changes

## Roadmap

### proxy routes
+ manages a single instance of a HAProxy, updating / restarting when backends change
+ implement a proxy ourselves replacing HAProxy, which adds hystrix functionality (circuit breaker / bulkhead protections)
+ statistics
