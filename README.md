> Node.js Portable is now integrated on [Neard](http://neard.io) !

<p align="center"><a href="https://github.com/crazy-max/nodejs-portable" target="_blank"><img width="100" src="https://github.com/crazy-max/nodejs-portable/blob/master/res/logo.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/nodejs-portable.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/nodejs-portable/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://ci.appveyor.com/project/crazy-max/nodejs-portable"><img src="https://img.shields.io/appveyor/ci/crazy-max/nodejs-portable.svg?style=flat-square" alt="AppVeyor"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/nodejs-portable"><img src="https://goreportcard.com/badge/github.com/crazy-max/nodejs-portable?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/nodejs-portable"><img src="https://img.shields.io/codacy/grade/03ea4cd8c645497aba77b5e462b5118c.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=QEEZEYZ6QTKGU"><img src="https://img.shields.io/badge/donate-paypal-7057ff.svg?style=flat-square" alt="Donate Paypal"></a>
</p>

Read this in other languages: [English](README.md), [简体中文](README.zh-cn.md).

## About

A single EXE written in [Go](https://golang.org/) to make [Node.js](http://nodejs.org/) portable on Windows systems.<br />
Tested on Windows 7, Windows 8.1 and Windows 10.

![](res/screenshots/main-20170915.gif)
> Main window of Node.js Portable

Configuration file `nodejs-portable.conf` is generated at first launch :

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

> If an exception happened, take a look into `nodejs-portable.log` for further information.

### Command line

Node.js Portable can be used through the command line to inject arguments directly to node :

```
$ nodejs-portable.exe --version
v9.5.0
```

> Take a look into `nodejs-portable.log` if you have any issue.

## Building

* Install [Go](https://golang.org/dl/) 1.11+
* Add Go to your PATH (ex. `C:\Go\bin`)
* Install [Java SE Development Kit](http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html) 1.8+
* Add Java to your PATH (ex. `C:\Program Files (x86)\Java\jdk1.8.0_144\bin`)
* Install [Apache Ant](http://ant.apache.org/bindownload.cgi) 1.9+
* Add Ant to your PATH (ex. `C:\apache-ant\bin`)

Then,

* Clone this repository to `$GOPATH/src/github.com/crazy-max/nodejs-portable`
* Run `ant release`. The artefact will be available in `bin\release`

If you don't want to use Java/Ant to build the project, run :

```
set GOARCH=386
go get -u github.com/Masterminds/glide
glide install -v
go generate -v
go build -v -ldflags "-s -w"
```

## How can i help ?

All kinds of contributions are welcomed :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Paypal](res/paypal-donate.png)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=QEEZEYZ6QTKGU)

## License

MIT. See `LICENSE` for more details.<br />
USB icon credit to [Dakirby309](http://dakirby309.deviantart.com/).
