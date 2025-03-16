# Hallaca

> A framework to analyze Node.js packages, written in Go

## About

Hallaca is primarily a security analysis tool for Node.js packages. Inspired by ESLint - but for your dependencies and your projects meta information.

## Motivation

Node.js is a very popular platform for developers. It has a rich ecosystem via npm, its package manager. However, due to its popularity it is also constantly targeted by supply chain attacks.
Hallaca aims to provide an open source framework to analyze package metadata and dependencies, ensuring they are safe to use.

Just like ESLint, Hallaca is meant to be extensible. So users can write and enforce their own policies.

## Status

Very much WIP.

## FAQ

### Why Go?

To get a different perspective of writing Node.js tools. I've made a prototype with TypeScript and it was fine. But after a decade of writing JavaScript, there's only so much you can learn.
Go has proven to be a viable language for writing Node.js tools (e.g. esbuild) so this project serves as an opportunity for me to evaluate this approach. Also who says no to potentially faster tooling?

Having worked with C, I was yearning to do more system level programming again and Go might just look to scratch that itch.
