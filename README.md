# MCC

MCC is a High-Level Net "Blaster" for the Model-Checking Contest (MCC).

We use the term *blaster*, instead of *unfolding*, to underline the fact that
our translation is very crude. A comparable tool is *andl*, distributed with
[Marcie](http://www-dssz.informatik.tu-cottbus.de/DSSZ/Software/Marcie). In
the future, we plan to use the transformation to compute interesting properties
of the models, like symetries and/or set of places that can be clustered
together.

At the moment, we simply recognize the case of high-level nets with a single
variable used with a circular symmetry (basically, this is a "scalar set"). This
allows us to manage very big instances of the *Philosophers* model.

## Getting Started

Before building the tool you need to install a recent  Go distribution at
<https://golang.org/doc/install>. Then you can install the software using the
*go get* command.

```bash
go get github.com/dalzilio/mcc
```

## Running the program

The *mcc hlnet* command accepts PNML files for high-level nets provided by the MCC (also tagged as COL).
These files generally have the name *model.pnml*.
You can invoke the *hlnet* command on this file as follows.

```bash
$> mcc hlnet -i model.pnml
```

You can obtain info on the other available options by using the *help* command:

```text
$> mcc help
collection of tools for the MCC

Usage:
  mcc [command]

Available Commands:
  hlnet       generates a .net or .tpn file from a PNML file describing a high-level net
  lola        generates a net file in the LoLa format from a PNML file describing a high-level net

Use "mcc [command] --help" for more information about a command.
$> mcc help hlnet
```

## License

This software is distributed under the [CECILL-B](http://www.cecill.info) license.
A copy of the license agreement is found in file LICENSE.

## Authors

* **Silvano DAL ZILIO** -  [LAAS/CNRS](https://www.laas.fr/)
