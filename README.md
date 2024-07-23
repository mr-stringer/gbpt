# GBPT

**This project is in it's early stages.**

GBPT is a project that aims to offer a 'Pricing-as-Code' solution for VM and
disk infrastructure for Azure.

## Why there is a requirement for Pricing-as-Code

Getting an idea of how much Azure infrastructure might cost can be tricky.
The most popular option is to use the
[Azure Calculator](https://azure.microsoft.com/en-gb/pricing/calculator/), which
is a web-based tool involving lots of clicking and drop down menus. Azure
Calculator allows you to quickly get pricing for any Azure service, but
it difficult to manage at scale.

Using Azure Calculator for very large solutions becomes problematic. As more
systems are added, the web interface slows down. When working with hundreds of
lines, the web interface is unusable. Moreover, when changes are required, such
as deploying systems to a different region, the user is required to make lots
of mouse-based manual changes. On large configurations, this can take it many
hours.

The advantages of Pricing-as-Code are as follows:

* Pricing as code plain text (in this case YAML) files
* Files can be version controlled in git and provide a full audit trail
* Changes to SKUs can be done manually with a graphical text editor or can
  use find/replace techniques in an editor or on the command line.
  
## Usage

gbpt supports two flags:

* -l is used to set the logging level which may be set to:
  * 'e' for error
  * 'w' for warnings
  * 'i' for information
  * 'd' for debug
* -printConfig is used to print the Azure config. In this mode log level is set
  to warning and the application stops as soon as the configuration is printed.

### Examples

To set logging to 'information' run:

```console
foo@bar:~$ ./gbpt -l=i
```

To print the Azure configuration run:

```console
foo@bar:~$ ./gbpt -printConfig
```