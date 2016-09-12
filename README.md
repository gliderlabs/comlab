# Gosper

Glider Labs app framework and development utility.

Gosper is a tool for setting up and developing Golang apps (optionally with React)
based on [github.com/gliderlabs/pkg/com](https://github.com/gliderlabs/pkg/tree/master/com).

This means highly modularized, extensible, and otherwise fairly opinionated. Our
first priority is serving Glider Labs apps before making this a general use
framework.

## Current Status

At this moment it is housing the dev runner / harness used by several apps. Next
major addition would be templates for creating new Gosper projects. Over time,
development tools and workflows used by projects will be encoded here.

## Packages

Each project may have its own `pkg` directory for potentially re-useable Go
packages. Common packages live at `github.com/gliderlabs/pkg`, including `com`.
However, packages dependent on Gosper structure and opinionated uses of those
shared packages would live in the `pkg` directory here for now.
