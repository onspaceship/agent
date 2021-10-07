_Note: This is the development version of the Spaceship Agent, so there may be code here that is not running in your cluster._

<p align="center">
  <img src="https://static.onspaceship.com/FullColor.svg" width="150">
</p>

<h3 align="center">
  Agent
</h3>

<p align="center">
  The Spaceship Kubernetes Cluster Agent
</p>

---

The Spaceship Agent is a small program that runs inside a [Kubernetes](https://kubernetes.io/) cluster to deliver software updates from [the Spaceship platform](https://spaceship.run/). 

Its design is intentionally minimal, so as to not interfere with other operational tasks in the cluster. At the moment, that is solely running Jobs, updating the `image:` field on Deployments, and watching those Deployments for their rollout status. 

The Agent can also update itself as instructed by Spaceship, which allows us to keep the Agent up to date and safely roll out new versions. If there is ever any material change in functionality for the Agent, we will inform you ahead of rolling out those changes.

## Installation

We build Spaceship using Spaceship, so container images of the Agent are always available here:

```
registry.onspaceship.com/spaceship/agent:master
```

That will be the latest version with the code you see here. You will also find tags available for every commit ID (SHA hash) on the repo as well, if you want a specific version. 

## Development

The Agent requires [go 1.16 or higher](https://golang.org/) to build. 

We use [Cobra](https://github.com/spf13/cobra) for creating a CLI. The default command is to connect to the Spaceship platform and start watching for Deployment changes. We have a few sub-commands to access specific functionality during development. You can see those with `go run . help`.

## Contributing

1. [Fork the repo](https://github.com/onspaceship/agent/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request here on GitHub
