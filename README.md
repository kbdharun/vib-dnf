# vib-dnf

This repository contains the Vib (Vanilla Image Builder) plugin for the DNF (Dandified YUM) package manager.

## Usage

### Using DNF Plugin with Vib

It can be used in a workflow with the following syntax:

```yml
- uses: vanilla-os/vib-gh-action@v1.0.0
  with:
    recipe: 'recipe.yml'
    plugins: 'kbdharun/vib-dnf:v1.0.0'
```

### Using the Plugin in a Recipe

```yml
modules:
  - name: "Install base packages"
    type: "dnf"
    options:
      allowerasing: true
      skip_broken: true
      skip_unavailable: false
      allow_downgrade: false
      downloadonly: false
      security: false
      bugfix: false
      enhancement: false
      extra_flags:
        - "--nogpgcheck"
        - "--refresh"
    sources:
      - packages:
          - vim
          - git
          - httpd
          - python3
      - path: "/path/to/packages-list.txt"
```

#### Examples

- Basic Package Installation:

```yml
modules:
  - name: "Install Development Tools"
    type: "dnf"
    options: {}
    sources:
      - packages:
          - gcc
          - make
          - automake
          - autoconf
```

- Apply Security Updates (on existing packages):

```yml
modules:
  - name: "Apply Security Updates"
    type: "dnf"
    options:
      security: true
    sources:
      - packages:
          - "*"
```

- Install packages from a file:

```yml
modules:
  - name: "Install Packages from List"
    type: "dnf"
    options:
      skip_unavailable: true
    sources:
      - path: "dir/packages.txt"
```

## Building the plugin

The plugin can be built locally with the following commands:

```sh
go get ./...
go build -trimpath -buildmode=plugin -o dnf.so -v ./...
```
