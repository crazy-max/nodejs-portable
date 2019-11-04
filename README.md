> For more portable applications, checkout [Portapps](https://portapps.io) !

<p align="center"><a href="https://github.com/crazy-max/nodejs-portable" target="_blank"><img width="100" src="https://github.com/crazy-max/nodejs-portable/blob/master/res/logo.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/nodejs-portable.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/nodejs-portable/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/actions"><img src="https://github.com/crazy-max/nodejs-portable/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/nodejs-portable"><img src="https://goreportcard.com/badge/github.com/crazy-max/nodejs-portable?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/nodejs-portable"><img src="https://img.shields.io/codacy/grade/03ea4cd8c645497aba77b5e462b5118c.svg?style=flat-square" alt="Code Quality"></a>
  <br /><a href="https://www.patreon.com/crazymax"><img src="https://img.shields.io/badge/donate-patreon-f96854.svg?logo=patreon&style=flat-square" alt="Support me on Patreon"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

Read this in other languages: [English](README.md), [简体中文](README.zh-cn.md).

## About

A single EXE written in [Go](https://golang.org/) to make [Node.js](http://nodejs.org/) portable on Windows systems.<br />
Tested on Windows 7, Windows 8.1 and Windows 10.

![](res/screenshots/main-20170915.gif)
> Main window of Node.js Portable

Configuration file `nodejs-portable.conf` is generated at first launch:

![](res/screenshots/files-20171227.png)

## Installation

* Download the [latest release](https://github.com/crazy-max/nodejs-portable/releases/latest).
* Put `nodejs-portable.exe` in a new empty folder.

> Add it to a white list if you receieved a warning from anti-virus software.

## Getting started

Run `nodejs-portable.exe`, then choose task :
* **1** to install node.js by entering version and architecture.
* **2** to launch Node.js shell.

> If you have already installed Node.js, create a folder named `app` and put your node.js environment inside, then launch `nodejs-portable.exe`. (See [#35](https://github.com/crazy-max/nodejs-portable/issues/35))

### Config file `nodejs-portable.conf`

* `workPath` : Shell working dir (can be relative to `nodejs-portable.exe`).
* `customPaths` : Array of custom paths that will be placed in the `PATH` environment variable (paths can be relative to `nodejs-portable.exe`).
* `immediateMode`: Set this to `true` if you want to use node shell immediately.
* `shell`: Shell to be used. Can be `cmd` (default) or `powershell`.

> If an exception happened, take a look into `nodejs-portable.log` for further information.

### Command line

Node.js Portable can be used through the command line to inject arguments directly to node:

```
$ nodejs-portable.exe --version
v9.5.0
```

> Take a look into `nodejs-portable.log` if you have any issue.

## Building

* Install [Go](https://golang.org/dl/) 1.12+
* Add Go to your PATH (ex. `C:\Go\bin`)
* Install latest version of [Mage](https://github.com/magefile/mage/releases/latest)
* Add Mage to your PATH
* Clone this repository
* Run `mage build` inside. The artifact will be available in `bin\release`

## How can I help ?

All kinds of contributions are welcome :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Support me on Patreon](res/patreon.png)](https://www.patreon.com/crazymax) 
[![Paypal](res/paypal.png)](https://www.paypal.me/crazyws)

## License

MIT. See `LICENSE` for more details.<br />
USB icon credit to [Dakirby309](http://dakirby309.deviantart.com/).
