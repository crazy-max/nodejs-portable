> For more portable applications, checkout [Portapps](https://portapps.io) !

<p align="center"><a href="https://github.com/crazy-max/nodejs-portable" target="_blank"><img width="100" src="https://github.com/crazy-max/nodejs-portable/blob/master/res/logo.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/nodejs-portable.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/nodejs-portable/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/actions?workflow=build"><img src="https://img.shields.io/github/workflow/status/crazy-max/nodejs-portable/build?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/nodejs-portable"><img src="https://goreportcard.com/badge/github.com/crazy-max/nodejs-portable?style=flat-square" alt="Go Report"></a>
  <br /><a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

Read this in other languages: [English](README.md), [简体中文](README.zh-cn.md).

## :warning: Abandoned project

This project is not maintained anymore and is abandoned. Feel free to fork and make your own changes if needed.

Thanks to everyone for their valuable feedback and contributions.

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

## Getting started

Run `nodejs-portable.exe`, then choose task:
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
* Run `mage build` inside. The artifact will be available in `bin`

## Contributing

Want to contribute? Awesome! The most basic way to show your support is to star the project, or to raise issues.

You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) or by
making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely!

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.<br />
USB icon credit to [Dakirby309](http://dakirby309.deviantart.com/).
