@ECHO OFF
SETLOCAL EnableDelayedExpansion

::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::                                                                                ::
::  Node.js Portable                                                              ::
::                                                                                ::
::  A DOS Batch script to make Node.js portable on Windows systems.               ::
::                                                                                ::
::  Copyright (C) 2013-2015 Cr@zy <webmaster@crazyws.fr>                          ::
::                                                                                ::
::  Node.js Portable is free software; you can redistribute it and/or modify      ::
::  it under the terms of the GNU Lesser General Public License as published by   ::
::  the Free Software Foundation, either version 3 of the License, or             ::
::  (at your option) any later version.                                           ::
::                                                                                ::
::  Node.js Portable is distributed in the hope that it will be useful,           ::
::  but WITHOUT ANY WARRANTY; without even the implied warranty of                ::
::  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the                  ::
::  GNU Lesser General Public License for more details.                           ::
::                                                                                ::
::  You should have received a copy of the GNU Lesser General Public License      ::
::  along with this program. If not, see http://www.gnu.org/licenses/.            ::
::                                                                                ::
::  Related post: http://goo.gl/gavL4                                             ::
::  Usage: nodejs-portable.bat                                                    ::
::                                                                                ::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

TITLE Node.js Portable v1.9

:: Settings
SET nodejsVersion=0.10.35
SET nodejsArch=x86
::SET proxyUrl=<url>:<port>
::SET proxyUser=<domain>\<user>
::SET proxyPwd=<password>

:: Batch vars (no edits necessary)
SET nodejsTask=%1
SET nodejsPath=%~dp0
SET nodejsPath=!nodejsPath:~0,-1!
SET nodejsWork=%nodejsPath%\work
SET npmPath=%nodejsPath%\node_modules\npm
SET npmGlobalConfigFilePath=%npmPath%\npmrc
SET nodejsInstallVbs=%TEMP%\nodejs_install.vbs
SET nodejsMsiPackage=node-v%nodejsVersion%-%nodejsArch%.msi
IF %nodejsArch%==x64 SET nodejsUrl=http://nodejs.org/dist/v%nodejsVersion%/x64/%nodejsMsiPackage%
IF %nodejsArch%==x86 SET nodejsUrl=http://nodejs.org/dist/v%nodejsVersion%/%nodejsMsiPackage%

:: Check if the menu selection is provided as a command line parameter
IF NOT "%nodejsTask%"=="" GOTO ACTION



::::::::::::::::::::::::::::::::::::::::
:MENU
::::::::::::::::::::::::::::::::::::::::
CLS
ECHO.
ECHO # Node.js Portable v1.9
ECHO.

ECHO  1 - Launch
ECHO  2 - Install
ECHO  9 - Exit
ECHO.
SET /P nodejsTask=Choose a task:
ECHO.



::::::::::::::::::::::::::::::::::::::::
:ACTION
::::::::::::::::::::::::::::::::::::::::
IF %nodejsTask% == 1 GOTO LAUNCH
IF %nodejsTask% == 2 GOTO INSTALL
IF %nodejsTask% == 9 GOTO EXIT
GOTO MENU



::::::::::::::::::::::::::::::::::::::::
:INSTALL
::::::::::::::::::::::::::::::::::::::::

:: Check if node.js is installed
IF EXIST "%nodejsPath%\node.exe" ECHO node.js is already installed... && GOTO EOF

:: Relocate and create temp dir (workaround for permission issue)
SET TEMP=%nodejsPath%\tmp
IF NOT EXIST "%TEMP%" MKDIR "%TEMP%"

:: Prepare cscript to download node.js
ECHO WScript.StdOut.Write "Download " ^& "%nodejsUrl%" ^& " ">%nodejsInstallVbs%
:: Switched to 'WinHttp.WinHttpRequest.5.1'
ECHO dim http: set http = createobject("WinHttp.WinHttpRequest.5.1") >>%nodejsInstallVbs%
IF DEFINED proxyUrl ECHO http.SetProxy 2, "%proxyUrl%", "localhost" >>%nodejsInstallVbs%
ECHO dim bStrm: set bStrm = createobject("Adodb.Stream") >>%nodejsInstallVbs%
:: Open in asynchronous mode.
ECHO http.Open "GET", "%nodejsUrl%", True >>%nodejsInstallVbs%
IF DEFINED proxyUser IF DEFINED proxyPwd ECHO http.SetCredentials "%proxyUser%", "%proxyPwd%", "1" >>%nodejsInstallVbs%
ECHO http.Send >>%nodejsInstallVbs%
:: Every second write a '.' until the download is complete
ECHO while http.WaitForResponse(0) = 0 >>%nodejsInstallVbs%
ECHO   WScript.StdOut.Write "." >>%nodejsInstallVbs%
ECHO   WScript.Sleep 1000 >>%nodejsInstallVbs%
ECHO wend >>%nodejsInstallVbs%
:: ECHO WScript.StdOut.Write vbCRLF >>%nodejsInstallVbs%
:: Write the HTTP status code onto console
ECHO WScript.StdOut.WriteLine " [HTTP " ^& http.Status ^& " " ^& http.StatusText ^& "]" >>%nodejsInstallVbs%
ECHO with bStrm >>%nodejsInstallVbs%
ECHO .type = 1 '//binary >>%nodejsInstallVbs%
ECHO .open >>%nodejsInstallVbs%
ECHO .write http.responseBody >>%nodejsInstallVbs%
ECHO .savetofile "%TEMP%\%nodejsMsiPackage%", 2 >>%nodejsInstallVbs%
ECHO end with >>%nodejsInstallVbs%

:: Download latest version in the current folder
cscript.exe /NoLogo %nodejsInstallVbs%

:: Extract the MSI package
ECHO Install node.js in %nodejsPath%...
msiexec /a "%TEMP%\%nodejsMsiPackage%" /qn TARGETDIR="%nodejsPath%"
XCOPY "%nodejsPath%\nodejs" "%nodejsPath%" /s /e /i /h /y

:: Clean folders
RMDIR /s /q "%nodejsPath%\nodejs"
RMDIR /s /q "%TEMP%"
IF EXIST "%nodejsPath%\%nodejsMsiPackage%" DEL "%nodejsPath%\%nodejsMsiPackage%"

:: Finish installation
ECHO.
IF EXIST "%nodejsPath%\node.exe" ECHO node.js successfully installed in '%nodejsPath%'
IF NOT EXIST "%nodejsPath%\node.exe" ECHO An error occurred during the installation.
GOTO PREPARE



::::::::::::::::::::::::::::::::::::::::
:LAUNCH
::::::::::::::::::::::::::::::::::::::::

:: Check if node.js is installed
IF NOT EXIST "%nodejsPath%\node.exe" ECHO node.js is not installed... Please install first... && GOTO EOF
IF NOT %nodejsTask% == 0 GOTO PREPARE

:: Where is git installed? Set temporary path.
SET WHEREISGIT=
IF /i NOT "%PROCESSOR_ARCHITECTURE%"=="x86" SET WHEREISGIT=\Wow6432Node
FOR /F "tokens=2*" %%F in ('REG QUERY HKLM\SOFTWARE%WHEREISGIT%\Microsoft\Windows\CurrentVersion\Uninstall\Git_is1 /v InstallLocation') DO SET GIT=%%G
SET PATH=%PATH%;%GIT%cmd

:: Init node vars
cmd.exe /k "cd "%nodejsWork%" && "%nodejsPath%\nodevars.bat" && "%nodejsPath%\npm" config set globalconfig "%npmGlobalConfigFilePath%" --global"
GOTO MENU



::::::::::::::::::::::::::::::::::::::::
:PREPARE
::::::::::::::::::::::::::::::::::::::::

:: Relocate and edit NPM global config file
ECHO prefix = %nodejsPath%\ >%npmGlobalConfigFilePath%
ECHO globalconfig = %npmPath%\npmrc >>%npmGlobalConfigFilePath%
ECHO globalignorefile = %npmPath%\npmignore >>%npmGlobalConfigFilePath%
ECHO init-module = %npmPath%\init.js >>%npmGlobalConfigFilePath%
ECHO cache = %npmPath%\cache >>%npmGlobalConfigFilePath%

IF NOT EXIST "%nodejsWork%" MKDIR "%nodejsWork%"
IF NOT EXIST "%npmPath%\npmignore" ECHO. 2>"%npmPath%\npmignore"
IF NOT EXIST "%npmPath%\init.js" ECHO. 2>"%npmPath%\init.js"
IF NOT EXIST "%npmPath%\cache" MKDIR "%npmPath%\cache"
IF %nodejsTask% == 1 SET nodejsTask=0 && GOTO LAUNCH
GOTO EOF



::::::::::::::::::::::::::::::::::::::::
:EOF
::::::::::::::::::::::::::::::::::::::::

ECHO.
PAUSE
GOTO MENU



::::::::::::::::::::::::::::::::::::::::
:EXIT
::::::::::::::::::::::::::::::::::::::::

ENDLOCAL
