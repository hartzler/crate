
## process

A unix process that runs in a container

## container

An instance of a component, consisting of one or more processes, running in an  environment, on a machine.

## component

A binary w/ metadata describing how components connect and interact, along with storage volume details.

Run in one or more containers on one or more machines.

## service

A set of components and their connections.

Services run on one or more machines.

## product

A set of services that work together to provide or satisfy specific use cases.

Each use case must be:
* economic
* correct
* intuitive
* safe (bounded risk)
* available
* performant
* scalable
* trackable

This requires the input of several roles or perspectives:
* business (economic)
* domain expert (correct, intuitive, performant, available)
* UI/UX (correct, intuitive)
* development (economic, correct, available, performant, scalable)
* operations (economic, available, performant, scalable)

Products have dedicated machines (can be virtual for multi-tenanting / density).

## program

a set of products that work together to provide a high level function.
