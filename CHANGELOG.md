# Changelog

## 2.3.0 (2017/12/27)

* Add ability to set a custom work path (Issue #35)
* Move Node.js installation to `app` subfolder
* Remove UNC path while setting PATH and NODE_PATH
* Use SetConsoleTitle instead of exec cmd
* Switch to [Golang Dep](https://github.com/golang/dep) as dependency manager

> 💡 UPGRADE NOTE : Move everything in a new folder named `app` except `cache`, `work` and `nodejs-portable.*`.

## 2.2.1 (2017/10/01)

* Improve coding quality
* Add `ImmediateMode` (Issue #27)
* Translated `README` into Chinese

## 2.2.0 (2017/08/29)

* Resolve absolute paths in `customPaths` (Issue #30)
* Merge `gitPath` and `pythonPath` in `customPaths` (Issue #30)
* Remove `node_modules/npm` from PATH (Issue #30)

> 💡 UPGRADE NOTE : `gitPath` and `pythonPath` have been removed from the config file and must be moved to `customPaths`

## 2.1.2 (2017/08/26)

* Upgrade to Go 1.9
* Add Building instructions (Issue #29) 
* Add config file to customize Git and Python path (Issues #10 #29) 

## 2.1.1 (2017/08/22)

* Add NODE_PATH to the shell (Issue #29)
* Tmp folder not created while launching shell (Issue #28)

## 2.1.0 (2017/07/30)

* Download zip format if exists
* Create Node.js configuration in `etc` folder
* Add cache path in root folder
* Bug while writing npmrc configuration

## 2.0.0 (2017/07/27)

* BIG rewrite in Golang! (Issue #22)
* Use lessmsi instead of msiexec for long path issue
* Check Node.js version before downloading
* Automatically retrieve the latest stable version of Node.js as default choice
* Inject Node and Git (if exists) paths in system environment PATH while launching shell
* Add new release notification
* Add support guidelines
* MIT license

## 1.13 (2017/07/27)

* Crash when invalid characters in PATH variable (Issue #26)
* Push to working directory (Issue #25)
* "\Git\cmd" syntax error on launch (Issue #23)
* Switch to Node.js 6.11.1

## 1.12 (2017/04/27)

* Switch to Node.js 6.10.2
* Add CHANGELOG

## 1.11 (2016/01/30)

* Bug with npm init command (Issue #21)
* Error after installing: Registry key not found (Issue #12)

## 1.10 (2014/12/20)

* Prompt for choosing version and arch before install (Issue #20)

## 1.9 (2015/04/20)

* Proxy support (Issue #16)
* Switch to Node.js 0.10.35 by default

## 1.8 (2014/12/20)

* Switch to Node.js 0.10.34 by default

## 1.7 (2014/08/04)

* Where is git installed? Set temporary path (Issue #11)
* Switch to Node.js 0.10.29 by default

## 1.6 (2014/05/22)

* Allow menu selection from command line parameter (Issue #9)

## 1.5 (2014/02/20)

* x86 arch by default
* Add choice of OS architecture (Issue #7)

## 1.4 (2014/01/21)

* Switch to Node.js 0.10.24 by default

## 1.3 (2013/11/18)

* Remove trailing slash on nodejsPath
* Change install dir %CD% to %~dp0 (Issue #5)
* Avoid issue with space in path (Issue #4)
* Switch to Node.js 0.10.22 by default

## 1.2 (2013/08/31)

* Add progress during download (Issue #2)

## 1.1 (2013/08/16)

* Add LGPL LICENSE
* Switch to Node.js 0.10.7 by default

## 1.0 (2013/04/18)

* Initial version
