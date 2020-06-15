# DotRun Command

[![Build Status](https://travis-ci.org/c4s4/dotrun.svg?branch=master)](https://travis-ci.org/c4s4/dotrun)
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/dotrun)](https://goreportcard.com/report/github.com/c4s4/dotrun)
[![Codecov](https://codecov.io/gh/c4s4/dotrun/branch/master/graph/badge.svg)](https://codecov.io/gh/c4s4/dotrun)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

DotRun command loads dotenv file, loads its environment and runs given command in that environment.

## Installation

### Unix users (Linux, BSDs and MacOSX)

Unix users may download and install latest *dotrun* release with command:

```bash
sh -c "$(curl https://sweetohm.net/dist/dotrun/install)"
```

If *curl* is not installed on you system, you might run:

```bash
sh -c "$(wget -O - https://sweetohm.net/dist/dotrun/install)"
```

**Note:** Some directories are protected, even as *root*, on **MacOSX** (since *El Capitan* release), thus you can't install *dotrun* in */usr/bin* for instance.

### Binary package

Otherwise, you can download latest binary archive at <https://github.com/c4s4/dotrun/releases>. Unzip the archive, put the binary of your platform somewhere in your *PATH* and rename it *dotrun*.

## Usage

To run command *foo* with its arguments in the environment defined in *.env* file in current directory, type:

```bash
dotrun foo args...
```

*.env* file might define environment such as:

```bash
FOO=BAR
SPAM=EGGS
```

Command *foo* will then be able to access the environment defined in *.env* file.

You can specify another dotenv file with `-env file` option:

```bash
dotrun -env /etc/foo.env foo args...
```

You can also load multiple dotenv files, repeating `-env file` option on command line :

```bash
dotrun -env /etc/foo.env -env /etc/bar.env foo args...
```

The environment files are evaluated in the order of the command line, so that in previous example variables defined in *bar.env* would overwrite those defined in *foo.env*.

## Shell

Let's say you have following *.env* file:

```bash
FOO=BAR
```

You would probably expect following:

```bash
$ dotrun echo $FOO
BAR
```

But this is not what happens:

```bash
$ dotrun echo $FOO

```

Because `$FOO` will be evaluated by the shell before running dotrun and replaced with its value on command line. To have expected behavior, you must run:

```bash
$ dotrun -shell 'echo $FOO'
BAR
```

In this case, command `echo $FOO` will not be evaluated until it runs in a shell. This shell will run in environment defined with dotenv file passed on command line and will print expected value on the console.

Note that you could try to obtain expected result with command `dotrun 'echo $FOO'`, but this won't work because dotrun will try to run command `echo $FOO` which doesn't exist.

On Unix, dotrun will run command in a shell with `sh -c command` and `cmd /c command` on Windows.

*Enjoy!*
