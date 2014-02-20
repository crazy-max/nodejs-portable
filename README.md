# Node.js Portable

A DOS Batch script to make [Node.js](http://nodejs.org/) portable on Windows systems.

Tested on Windows XP, Windows Vista and Windows 7.

## Requirements

* Have an Internet connection (for installation when downloading the application).
* [WSH (Windows Script Host)](http://support.microsoft.com/kb/232211) : Open a command prompt and type ``cscript`` to check.

## Installation

* Put the ``nodejs-portable.bat`` in the same directory as ``node.exe`` or in an empty directory for a new installation.
* Run ``nodejs-portable.bat``.
* Choose task 2 to install or task 1 to launch node.js if it's already installed.

## Configuration

In the ``Settings`` section line 32, you can change :
* ``nodejsVersion`` : The node.js version.
* ``nodejsArch`` : The OS architecture (``x86`` or ``x64``)

## Note

If you have already installed node.js, just copy the folder where you want and launch the script.

## License

LGPL. See ``LICENSE`` for more details.

## More infos

http://www.crazyws.fr/dev/applis-et-scripts/node-js-portable-JWSN9.html