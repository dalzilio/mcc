# MCC

MCC is a High-Level Net "Blaster" for the Model-Checking Contest (MCC). This
tool can be used to unfold a PNML file (such as used for colored models in the
[Model-Checking Contest](https://mcc.lip6.fr/)) into a Petri net. We support the
generation of Petri nets in both the TINA (.net) and LOLA input formats.

We use the term *blaster*, instead of *unfolding*, to underline the fact that
our translation is, in its principle, very crude. Nonetheless, we have made many
improvements on the tool along the years and ```mcc``` is now a very efficient
solution, on par with (and on many instances better) than comparable tools used
for the same purpose.

In the future, we plan to use the transformation to compute interesting
properties of the models, like symetries and/or set of places that can be
clustered together.

## Installing MCC or building it from source

MCC is a classic Command-Line Interface tool. You can directly install it by
copying the right binary file in your system. You can find the executable for
the latest releases on [GitHub's release page for this
project](https://github.com/dalzilio/mcc/releases). We provide binary files for
[Windows](https://github.com/dalzilio/mcc/releases/download/v1.0.0/mcc.exe),
[Linux](https://github.com/dalzilio/mcc/releases/download/v1.0.0/mcc-linux) and
[MacOS
(Darwin)](https://github.com/dalzilio/mcc/releases/download/v1.0.0/mcc-darwin).

You also have the option to install the tool directly from source. For this, you
need first to install a recent Go distribution (available at
<https://golang.org/doc/install>). Then you can install the software using the
*go get* command.

```bash
$> go get github.com/dalzilio/mcc
```

You can browse the documentation for this tool on the [GoDoc page for the mcc
project](https://godoc.org/github.com/dalzilio/mcc).

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

## Features

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

* We have added a new option, ```--stats```, to print statistics about the
  computation, such as computation time, or the number of places and transitions
  in the resulting Place/Transition net. (We do not output a result when this
  option is used.)

* We have modified option ```--debug``` in order to add more visual information
  when displaying the resulting model using tool ```nd``` for the TINA toolbox.
  We still display information about types, variables and the expressions
  associated with arcs inside comments. We also add a copy of this information
  using the support for (sticky) notes nodes that is built inside TINA's net
  format. You can see an example of the result obtained on the TableDistance
  model below.
  
![tekst alternatywny](./docs/nd.png)

## Sources

The code repository includes instances of PNML models from the [MCC Petri Nets
Repository](https://pnrepository.lip6.fr/) inside the ```./benchmarks``` folder.
We provide instances for all the PNML colored models used in the 2019
Model-Checking Contest. These files are included in repository to be used for
benchmarking and continuous testing.

## License

This software is distributed under the [CECILL-B](http://www.cecill.info)
license. A copy of the license agreement is found in the [LICENSE](./LICENSE)
file.

## Authors

* **Silvano DAL ZILIO** -  [LAAS/CNRS](https://www.laas.fr/)
