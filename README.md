
<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/dalzilio/mcc">
    <img src="docs/mcc.png" alt="Logo" width="240" height="80">
  </a>

  <p align="center">
    A tool to remove colors from your High-Level Petri nets !
    <br />
    <a href="https://github.com/dalzilio/mcc#features"><strong>see what's new »</strong></a>
    <br />
    <!-- <a href="https://github.com/dalzilio/mcc">View Demo</a> -->
  </p>
</p>

## About

MCC is a tool designed for a very specific and narrow task: to transform the
models of High-Level Petri nets, given in [PNML](http://www.pnml.org/), into
equivalent Place/Transition nets.

[![Go Report Card](https://goreportcard.com/badge/github.com/dalzilio/mcc)](https://goreportcard.com/report/github.com/dalzilio/mcc)
[![GoDoc](https://godoc.org/github.com/dalzilio/mcc?status.svg)](https://godoc.org/github.com/dalzilio/mcc)
[![Release](https://img.shields.io/github/v/release/dalzilio/mcc)](https://github.com/dalzilio/mcc/releases)
![Go](https://github.com/dalzilio/mcc/workflows/Go/badge.svg)

## Overview

The name of the tool derives from the annual [Model-Checking
Contest](https://mcc.lip6.fr/), a competition of model-checking tools that
provides a large and diverse collection of PNML models. Our choice in naming
serves to underline the main focus of the tool, which is to provide a simple,
open and extensible solution to lower the access cost for developers that would
like to engage in this competition.

MCC supports the generation of Petri nets in both the TINA (.net) and LOLA input
formats. We have also recently added a new subcommand to output the result as a
P/T net in PNML format.

We have made many improvements on the tool along the years, and it is now a very
efficient solution, on par with (and on many instances more efficient than)
other tools used for the same purpose.

## Reference

* Silvano Dal Zilio. [MCC: a Tool for Unfolding Colored Petri Nets in PNML
  Format](https://hal.laas.fr/hal-02511881). _41st International Conference on
  Application and Theory of Petri Nets and Concurrency_, Jun 2020, Paris,
  France. ⟨hal-02511881⟩

## Installing MCC or building it from source

MCC is a classic Command-Line Interface tool. You can directly install it by
copying the right binary file in your system. You can find the executable for
the latest releases on [GitHub's release page for this
project](https://github.com/dalzilio/mcc/releases). We provide binary files for
[Windows](https://github.com/dalzilio/mcc/releases/download/v1.5.0/mcc.exe),
[Linux](https://github.com/dalzilio/mcc/releases/download/v1.5.0/mcc-linux) and
[MacOS (Darwin)](https://github.com/dalzilio/mcc/releases/download/v1.5.0/mcc-darwin).

You also have the option to install the tool directly from source. For this, you
need first to install a recent Go distribution (available at
<https://golang.org/doc/install>). Then you can install the software using the
`go get` command.

```bash
$> go get github.com/dalzilio/mcc
```

You can browse the documentation for this tool on the [GoDoc page for the MCC
project](https://godoc.org/github.com/dalzilio/mcc).

## Running the program

The `mcc pnml` command accepts PNML files for high-level nets provided by the
Model-Checking Contest (also tagged as COL) and generates a P/T net equivalent
file. These files generally have the name `model.pnml`. You can invoke the
`pnml` command on this file as follows.

```text
$> mcc pnml -i model.pnml
```

You can obtain info on all available options by using the `help` command:

```text
$> mcc help
mcc transforms High-Level Petri nets in PNML format into equivalent Place/Transition nets

Usage:
  mcc [command]

Available Commands:
  help        Help about any command
  info        Print statistics or generate textual version for use with NetDraw (nd)
  lola        Generate a P/T net using the LoLa format
  pnml        Generate a P/T net using the PNML format
  skeleton    Generate the skeleton of a colored net in .net format
  smpt        Generate a P/T net file for use with SMPT
  tina        Generate a P/T net file using Tina's .net format
  version     Print the version number

Flags:
  -h, --help      help for mcc
  -v, --version   version for mcc

Use "mcc [command] --help" for more information about a command.
```

## Features and recent modifications

* We provide a new command, `skeleton`, that generates a P/T net by coalescing
  all the colors in a place into a single one, and forgetting the guards on
  transitions. This is exactly the construction defined by S. Wallner and K.
  Wolf in [Skeleton Abstraction for Universal Temporal
  Properties](https://doi.org/10.1007/978-3-030-76983-3_10) (2021).

* We support the declaration of `finiteintrange` types in PNML. as well
  as the declaration of `Partition` and `PartitionElement`. This means
  that we can now unfold model VehicularWifi (surprise model in 2019)

* We deprecated command `hlnet` since version 2.0 and replaced it with `mcc
  info`, which provides similar functionalities. Option `--stats` can be used to
  print statistics about the unfolding of a colored net, such as the computation
  time, or the number of places and transitions in the resulting P/ net. (We do
  not output a result when this option is used.)

  ```text
  $> mcc info -i model.pnml --stats
  100 place(s), 429 transition(s), 858 arc(s), 0.027s
  ```
  
* Option `--debug` in command `mcc info` generates a P/T net "documentation" for
  a colored net with information that can be displayed with the tool NetDraw
  (`nd`), which is part of the [TINA
  toolbox](http://projects.laas.fr/tina/home.php). Wedisplay information about
  types, variables and the expressions associated with arcs inside comments. We
  also add a copy of this information using the support for (sticky) notes that
  is built inside TINA's net format. You can see an example of the result
  obtained on the TableDistance-COL model below.
  
  ![TableDistance-COL model in nd](./docs/nd.png)

* You can use parameter `-` with option `-o`, in most commands, to output the
  result of the unfolding on the standard output. This way it is possible to
  pipe the result of `mcc` to another program, for instance another conversion
  tool, such as `ndrio`, which is part of the [TINA
  toolbox](http://projects.laas.fr/tina/home.php). The following example shows
  how to output a result using the [Roméo](http://romeo.rts-software.org/)
  format.

  ```text
  $> mcc tina -i model.pnml -o - | ndrio -romeo -
  ```

* Command `mcc smpt` generates a P/T net in .net format, like with `mcc tina`,
  but with simpler identifiers. It also includes traceability information
  between places and transitions in the colored net, and their P/T equivalent.

* Since version 2.0, command `mcc pnml` always output results using "structured"
  place names. We also ensure that the XML attributes for place names and id are
  equal.

## Dependencies

The code repository includes instances of PNML models from the [MCC Petri Nets
Repository](https://pnrepository.lip6.fr/) located inside the `./benchmarks`
folder. We provide a selection of instances from all the PNML colored models
used in the Model-Checking Contest. These files are included in the
repository to be used for benchmarking and continuous testing.

## License

This software is distributed under the [CECILL-B](http://www.cecill.info)
license. A copy of the license agreement is found in the [LICENSE](./LICENSE)
file.

## Authors

* **Silvano DAL ZILIO** -  [LAAS/CNRS](https://www.laas.fr/)
