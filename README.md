# ittybittyfeaturechecker
A minimalist, opinionated feature switch service.

## Feature Switches

This service allows applications to check if a given feature
is enabled or not.  Typically this is used for gating access
to experimental new features.  This is the primary use case
for this service.

## Minimalist and Opininated

Features are stored in memory and loaded from a json file. There
is no database or means of adding or updating feature switches
while the service is running. 

This simplifies deployment and operations since there is no external
database to setup and maintain.

Feature changes should be done via a PR/deployment as part of your
normal release processes.  This keeps a single source of truth and
various environments in sync. No more, "oops I forgot to add the feature
to production".  

Kill switches/scramming a feature is best checked for on the client
side.

Features are defined using expressions that need to evaluate to true
or false. 

