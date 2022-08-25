[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/FFFFFaraway/farctl/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/FFFFFaraway/farctl.svg)](https://pkg.go.dev/github.com/FFFFFaraway/farctl)
[![Go](https://github.com/FFFFFaraway/farctl/actions/workflows/go.yml/badge.svg)](https://github.com/FFFFFaraway/farctl/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/FFFFFaraway/farctl)](https://goreportcard.com/report/github.com/FFFFFaraway/farctl)

# What is it

`farctl`is a simple CLI tool for machine learning engineer to deploy [MPIJob](https://github.com/FFFFFaraway/MPI-Operator) in Kubernetes cluster without Kubernetes-related knowledge or manually deployment of yaml files. Imitated from [this project](https://github.com/kubeflow/arena), I reinvented the wheel for learning purpose again.

# How to install

## Requirements

- Kubernetes cluster. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster.
- [kubectl]([Install Tools | Kubernetes](https://kubernetes.io/docs/tasks/tools/))
- [helm]([Helm | Installing Helm](https://helm.sh/docs/intro/install/))
- [Golang]([Download and install - The Go Programming Language](https://go.dev/doc/install))

## Installation

```bash
go install github.com/FFFFFaraway/farctl@latest
```

# How to use

## Submit MPIJob

1. You'll need to write deep learning code using horovod. For example [here](https://github.com/FFFFFaraway/sample-python-train)

2. You'll need to upload the code to some public available platform like [github](https://github.com), or [gitlab](https://about.gitlab.com), so that the container could pull the code down.

3. Submit the job, for example:

   ```bash
   farctl submit test -i https://github.com/FFFFFaraway/sample-python-train.git -c "python generate_data.py" -c "python main.py" --gang -n 2
   ```

   - Test is the name of the submitted MPIJob
   - -i denote the url of git clone
   - -c denote the command as entry point. we can have multiple commands by using multiple -c
   - -gang denote we'll use [gang scheduler](https://github.com/FFFFFaraway/gang-scheduler). But it's needed a extra installation.
   - -n denote the number of workers to be created
   - Other options can be found by typing `farctl submit -h`

## List MPIJob

```bash
farctl list
```

We could monitor the status of mpijobs here:

```bash
Namespace  Name  ReadyWorkers/Total  LauncherStatus  Age
farctl     test  1/2                 WaitingWorkers  1m17s
```

```bash
Namespace  Name  ReadyWorkers/Total  LauncherStatus  Age
farctl     test  2/2                 Running         2m3s
```

## Get MPIJob Log

When the `LauncherStatus` become `Running`, we can access the log of the MPIJob:

```bash
farctl log test
```

## Get applyed MPIJob Configuration

```bash
farctl get test
```

## Delete MPIJob

```bash
farctl delete test
```

