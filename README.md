DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.

# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**


# Snap processor plugin - Tags Filter

This plugin allows filtering metrics by tags.

It's used in the [Snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

* [golang 1.7+](https://golang.org/dl/) (needed only for building)

### Installation

#### Download Tags Filter plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-processor-tags-filter/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-processor-tags-filter

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-processor-tags-filter.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage

To use this plugin you need to specify filtering rules in processor's config section of your task manifest.

To create a rule for a tag you must add a config field with tag name followed by either `.allow` or `.deny` as a key and comma-separated values as a value.

Example plugin config:
```
"foo.allow": "foovalue,fooval",
"foo.deny":  "badvalue",
"baz.allow": "bazval"
```
Which evaluates to "filter metrics that have value `foovalue` or `fooval` of tag `foo` or have value `bazval` of tag `baz`, but don't have value `badvalue` of tag `foo`".

## Documentation

### Examples

Example filtering metrics from [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil) and publishing with [file publisher plugin](https://github.com/intelsdi-x/snap-plugin-publisher-file).

[Start Snap daemon](https://github.com/intelsdi-x/snap#running-snap) (in this case with logging set to 1 and trust disabled):
```
$ snapteld -l 1 -t 0 &
```

Download and load Snap plugins:
```
$ snaptel plugin load snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-processor-tags-filter
$ snaptel plugin load snap-plugin-publisher-file
```

Create a task (you can find example task manifests in [examples/tasks](https://github.com/intelsdi-x/snap-plugin-processor-tags-filter/tree/master/examples/tasks) folder):
```
$ snaptel task create -t psutil-tags-filter-file.yml
Using task manifest to create task
Task created
ID: 870262de-e0ad-45b7-9983-ea33d1f2c925
Name: Task-870262de-e0ad-45b7-9983-ea33d1f2c925
State: Running
```

Ensure task is running and collecting metrics:
```
$ snaptel task watch 870262de-e0ad-45b7-9983-ea33d1f2c925
```

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-processor-tags-filter/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-processor-tags-filter/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, an open telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions! 

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Klaudiusz Dembler](https://github.com/kdembler)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
