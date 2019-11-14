> For more portable applications, checkout [Portapps](https://portapps.io) !

<p align="center"><a href="https://github.com/crazy-max/nodejs-portable" target="_blank"><img width="100" src="https://github.com/crazy-max/nodejs-portable/blob/master/res/logo.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/nodejs-portable.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/nodejs-portable/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/nodejs-portable/actions"><img src="https://github.com/crazy-max/nodejs-portable/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/nodejs-portable"><img src="https://goreportcard.com/badge/github.com/crazy-max/nodejs-portable?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/nodejs-portable"><img src="https://img.shields.io/codacy/grade/03ea4cd8c645497aba77b5e462b5118c.svg?style=flat-square" alt="Code Quality"></a>
  <br /><a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

阅读此文档其他语言的版本： [English](README.md), [简体中文](README.zh-cn.md).

## 关于

这是一个用 [Go 语言](https://golang.org/) 写的小程序,可以绿化 Windows 系统上的 [Node.js](http://nodejs.org/) 开发环境<br />
已经在 Windows 7 , Windows 8.1 和 Windows 10 上完成测试。

![](res/screenshots/main-20170915.gif)
> Node.js Portable 的主窗口

配置文件 `nodejs-portable.conf` 会在初次启动时被创造：

![](res/screenshots/files-20171227.png)

## 安装

* 下载 [最新的发布版本](https://github.com/crazy-max/nodejs-portable/releases/latest) 。
* 将 `nodejs-portable.exe` 放入一个空文件夹。

> 非常不推荐将 `nodejs-portable.exe` 放入带中文的路径中,可能会报各种诡异的错误Orz

## 开始使用

运行 `nodejs-portable.exe` ,然后按提示输入：
* **1** 自动安装, 输入版本号和系统架构, 然后程序会自动安装 Node.js 环境。
* **2** 自动配置并运行 Node.js 开发环境。

> 如果你已经安装完成了 Node.js , 新建一个 `app` 文件夹, 将你的环境放入其中, 再执行 `nodejs-portable.exe` 即可 ([#35](https://github.com/crazy-max/nodejs-portable/issues/35))

###  `nodejs-portable.conf` 配置文件

* `workPath` : 环境的工作目录 (可以是相对于 `nodejs-portable.exe` 的相对路径)。
* `customPaths` : 一组用于放入 `PATH` 环境变量 的路径 (可以是相对于 `nodejs-portable.exe` 的相对路径)。
* `immediateMode`: 立即模式, 将其设置为 `true` 来直接打开运行时环境。
* `shell`: Shell to be used. Can be `cmd` (default) or `powershell`.

> 如果出现了异常, 请检查或提供 `nodejs-portable.log` 来获取更多信息.

### Command line

Node.js Portable can be used through the command line to inject arguments directly to node:

```
$ nodejs-portable.exe --version
v9.5.0
```

> Take a look into `nodejs-portable.log` if you have any issue.

## 构建

* 安装 [Go](https://golang.org/dl/) 1.12+
* 将 Go 加入你的 PATH 环境变量 (例如 `C:\Go\bin`)
* Install latest version of [Mage](https://github.com/magefile/mage/releases/latest)
* Add Mage to your PATH
* Clone this repository
* Run `mage build` inside. The artifact will be available in `bin`

## 我怎么支持项目？

All kinds of contributions are welcome :raised_hands:! The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon: You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) :clap: or by making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely! :rocket:

Thanks again for your support, it is much appreciated! :pray:

## 许可证

MIT。阅读 `LICENSE` 来获得更多细节。<br />
USB 图标感谢 [Dakirby309](http://dakirby309.deviantart.com/) 。<br />
中文翻译 [Retomehere](https://github.com/xiazeyu)。
