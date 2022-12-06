# Dev Diary 4

## Problem

Is it possible to benchmark test against the database?

## Pretext

Scaffolding test data depends upon randomly generating user names *or* HTTP calls; it's difficult to benchmark both of these as they introduce a dependence on something that *isn't* code we have written.

## Decision

Decide on a method of benchmarking even with the introduction of "external" code; benchmark tests can still provide valuable insights into whether or not a proposed change is an improvement or not, even if its not an objective measure of the code's performance.