# MCC

High-Level Nets Blaster for the Model-Checking Contest.

## Getting Started

Before building the tool you need to install a Go distribution at <https://golang.org/doc/install>.

### Installing

You can simply install the software using the *go vet* command then *go install*.

We use [Cobra]("https://github.com/spf13/cobra") to generate the CLI, hence the strange syntax.

```bash
go get github.com/dalzilio/mcc
go install github.com/dalzilio/mcc/cpn/cmd/mcc
```

## Running the program

Download a *model.pnml* file (high-level net, also tagged COL) from the MCC repository. The *hlnet* command can be invoked as follows.

```bash
mcc hlnet -i model.pnml
```

To find help on a command simply type

```bash
mcc help hlnet
```

## License

This software is distributed under the [CECILL-B]("http://www.cecill.info") license.
A copy of the license agreement is found in file LICENSE.

## Authors

* **Silvano DAL ZILIO** -  [LAAS/CNRS]("https://www.laas.fr/")
