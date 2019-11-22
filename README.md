# MCC

MCC is a High-Level Net "Blaster" for the Model-Checking Contest (MCC). This
tool can be used to unfold a PNML file (such as used for colored models in the
MCC) into a Petri net. We support the generation of Petri nets in both the TINA
(.net) and LOLA input formats.

We use the term *blaster*, instead of *unfolding*, to underline the fact that
our translation is very crude. A comparable tool is *andl*, distributed with
[Marcie](http://www-dssz.informatik.tu-cottbus.de/DSSZ/Software/Marcie). In
the future, we plan to use the transformation to compute interesting properties
of the models, like symetries and/or set of places that can be clustered
together.

## Building from source

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

## Release notes

* We recognize the case of high-level nets with a single variable used with a
  circular symmetry (basically, this is a "scalar set"). This allows us to
  manage very big instances of the *Philosophers* model.

* We have improved the performances on colored models with a large number of
  variables. This means that we can now unfold model BART in a few seconds.

* We support the generation of Petri nets with more structured place names
  (option ```--sliced```) and with labels that add traceability information to
  the colored model (with option ```--verbose```). With these options, the output
  of the tool becomes more deterministic.

* We now support the declaration of ```finiteintrange``` types in PNML. as well
  as the declaration of ```Partition``` and ```PartitionElement```. This means
  that we can now unfold model VehicularWifi (surprise model in 2019)

## License

This software is distributed under the [CECILL-B](http://www.cecill.info) license.
A copy of the license agreement is found in file LICENSE.

## Authors

* **Silvano DAL ZILIO** -  [LAAS/CNRS](https://www.laas.fr/)
