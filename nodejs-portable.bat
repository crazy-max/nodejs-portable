@ECHO OFF

:: 
::
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
::  Usage: nodejs-portable.bat [action] [target_dir] [work_dir]                   ::
::                                                                                ::
::    action:     1 = launch Node.js, 2 = Install Node.js                         ::
::                                                                                ::
::    target_dir: where to install/locate Node.js. Default = '.'                  ::
::                                                                                ::
::    work_dir:   directory where work will be performed                          ::
::                                                                                ::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

SET currentBatchTitle=Node.js Portable v1.10
TITLE %currentBatchTitle%

:: Settings
SET nodejsVersion=0.12.4
SET nodejsArch=x64
::SET nodejsArch=x86
::SET proxyUrl=<url>:<port>
::SET proxyUser=<domain>\<user>
::SET proxyPwd=<password>

:: Parameters (can be provided as command line arguments)
SET nonInteractiveMode=%~1
SET nodejsTask=%~1
SET nodejsPath=%~dpnx2
IF "%nodejsPath%"=="" SET nodejsPath=%CD%
IF "%nodejsPath:~-1%"=="\" SET nodejsPath=%nodejsPath:~0,-1%
SET nodejsWork=%~dpnx3
IF "%nodejsWork%"=="" SET nodejsWork=%nodejsPath%\work

:: Batch vars (no edits necessary)
SET npmPath=%nodejsPath%\node_modules\npm
SET npmGlobalConfigFilePath=%npmPath%\npmrc
SET nodejsMsiPackage=node-v%nodejsVersion%-%nodejsArch%.msi
IF %nodejsArch%==x64 SET nodejsUrl=http://nodejs.org/dist/v%nodejsVersion%/x64/%nodejsMsiPackage%
IF %nodejsArch%==x86 SET nodejsUrl=http://nodejs.org/dist/v%nodejsVersion%/%nodejsMsiPackage%

:: Check if the menu selection is provided as a command line parameter
IF NOT "%nodejsTask%"=="" GOTO :ACTION



::::::::::::::::::::::::::::::::::::::::
:MENU
::::::::::::::::::::::::::::::::::::::::
CLS
ECHO.
ECHO # %currentBatchTitle%
ECHO.
ECHO Target: Node.js version %nodejsVersion% in %nodejsPath% 
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
IF "%nodejsTask%" == "1" (
 CALL :LAUNCH
 GOTO :ACTION_DONE
)
IF "%nodejsTask%" == "2" (
 CALL :INSTALL
 GOTO :ACTION_DONE
)
IF "%nodejsTask%" == "9" GOTO :EXIT

echo Unknown action: %nodejsTask%

::::::::::::::::::::::::::::::::::::::::
:ACTION_DONE
::::::::::::::::::::::::::::::::::::::::
IF DEFINED nonInteractiveMode GOTO :EXIT

GOTO :MENU



::::::::::::::::::::::::::::::::::::::::
:INSTALL
::::::::::::::::::::::::::::::::::::::::

:: Check if Node.js is installed
IF EXIST "%nodejsPath%\node.exe" (
  ECHO Node.js is already installed in %nodejsPath%
  IF NOT DEFINED nonInteractiveMode PAUSE
  GOTO :EOF
)

:: Relocate and create temp dir (workaround for permission issue)
SET TEMP=%nodejsPath%\%~n0.tmp
SET nodejsInstallVbs=%TEMP%\nodejs_install.vbs
SET nodejsMsiPackageTempFile=%TEMP%\%nodejsMsiPackage%
IF NOT EXIST "%TEMP%" MKDIR "%TEMP%"

:: Prepare cscript to download Node.js
ECHO WScript.StdOut.WriteLine "Downloading " ^& "%nodejsUrl%" >%nodejsInstallVbs%
ECHO WScript.StdOut.WriteLine "         to " ^& "%nodejsMsiPackageTempFile%" >>%nodejsInstallVbs%
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
ECHO .savetofile "%nodejsMsiPackageTempFile%", 2 >>%nodejsInstallVbs%
ECHO end with >>%nodejsInstallVbs%

:: Download latest version in the current folder
cscript.exe /NoLogo %nodejsInstallVbs%

:: Extract the MSI package
ECHO Installing Node.js %nodejsVersion% in %nodejsPath%
msiexec /a "%nodejsMsiPackageTempFile%" /qn TARGETDIR="%nodejsPath%"
XCOPY "%nodejsPath%\nodejs" "%nodejsPath%" /s /e /i /h /y

:: Clean folders
RMDIR /s /q "%nodejsPath%\nodejs"
RMDIR /s /q "%TEMP%"
IF EXIST "%nodejsPath%\%nodejsMsiPackage%" DEL "%nodejsPath%\%nodejsMsiPackage%"

:: Finish installation
ECHO.
IF EXIST "%nodejsPath%\node.exe" ECHO Node.js successfully installed in %nodejsPath%
IF NOT EXIST "%nodejsPath%\node.exe" (
  ECHO An error occurred during the installation.
) ELSE (
  CALL :PREPARE
)

IF NOT DEFINED nonInteractiveMode PAUSE
GOTO :EOF


::::::::::::::::::::::::::::::::::::::::
:LAUNCH
::::::::::::::::::::::::::::::::::::::::
ECHO Configuring environement for Node.js %nodejsVersion% in %nodejsPath%

:: Check if Node.js is installed
IF NOT EXIST "%nodejsPath%\node.exe" (
  ECHO Node.js is not installed in %nodejsPath%. Please install first...
  IF NOT DEFINED nonInteractiveMode PAUSE
  GOTO :EOF
)
CALL :PREPARE

:: Where is git installed? Set temporary path.
ECHO Looking for Git...
SET GIT_HOME=
where git
IF ERRORLEVEL 1 (
  ECHO Git is not found in PATH, looking in Registry where it is installed...
  SET WHEREISGIT=
  IF /i NOT "%PROCESSOR_ARCHITECTURE%"=="x86" SET WHEREISGIT=\Wow6432Node
  FOR /F "tokens=2*" %%F in ('REG QUERY HKLM\SOFTWARE%WHEREISGIT%\Microsoft\Windows\CurrentVersion\Uninstall\Git_is1 /v InstallLocation') DO SET GIT_HOME=%%G
) ELSE (
  ECHO Git found in PATH
)
IF "%GIT_HOME%" == "" GOTO :DONE_WITH_GIT
SET GIT_CMD=%GIT_HOME%cmd
ECHO Adding Git to PATH %GIT_CMD%
SET PATH=%GIT_CMD%;%PATH%
:DONE_WITH_GIT

PATH=%nodejsPath%;%PATH%

:: Init node vars
IF DEFINED nonInteractiveMode (
  ECHO Configuring current command line processor...
  cd "%nodejsWork%"
  call "%nodejsPath%\nodevars.bat"
  call "%nodejsPath%\npm" config set globalconfig "%npmGlobalConfigFilePath%" --global
) ELSE (
  ECHO.
  ECHO Starting new command line processor...
  cmd.exe /k "cd "%nodejsWork%" && "%nodejsPath%\nodevars.bat" && ECHO. && ECHO Type 'exit' when you are done working with Node.js && CALL "%nodejsPath%\npm" config set globalconfig "%npmGlobalConfigFilePath%" --global"
)

GOTO :EOF



::::::::::::::::::::::::::::::::::::::::
:PREPARE
::::::::::::::::::::::::::::::::::::::::
ECHO Preparing Node.js %nodejsVersion% in %nodejsPath%

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
GOTO :EOF



::::::::::::::::::::::::::::::::::::::::
:EXIT
::::::::::::::::::::::::::::::::::::::::

ENDLOCAL

