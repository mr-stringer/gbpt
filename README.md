# GBPT

**This project is in it's early stages.**

GBPT is a project that aims to offer a 'Pricing-as-Code' solution for VM and
disk infrastructure for Azure.

## Why is Pricing-as-Code required?

Getting an idea of how much Azure infrastructure might cost can be tricky.
One option is to use the
[Azure Calculator](https://azure.microsoft.com/en-gb/pricing/calculator/), which
is a web-based tool involving lots of clicking and drop down menus. Azure
Calculator allows you to quickly get pricing for any Azure service, but
it difficult to manage at scale.

Using Azure Calculator for very large solutions becomes problematic. As more
systems are added, the web interface slows down. When working with hundreds of
lines, the web interface is unusable. Moreover, when changes are required, such
as deploying systems to a different region, the user is required to make lots
of manual changes. On large configurations, this can take it many hours. 