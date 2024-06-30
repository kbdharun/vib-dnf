# vib-dnf

This repository contains the Vib (Vanilla Image Builder) plugin for the DNF (Dandified YUM) package manager.

## Usage

It can be used in a workflow with the following syntax:

```yml
- uses: vanilla-os/vib-gh-action@v0.7.2
  with:
    recipe: 'recipe.yml'
    plugins: 'kbdharun/vib-dnf:v0.1'
```

## Building the plugin

The plugin can be built locally with the following commands:

```sh
go get ./...
go build -trimpath -buildmode=plugin -o dnf.so -v ./...
```
