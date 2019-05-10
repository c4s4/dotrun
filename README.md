# DotRun Command

DotRun command loads dotenv file, loads its environment and runs given command in that environment.

## Installation

Drop the binary for your platform in the *bin* directory of the archive somewhere in your `PATH` and rename it *dotrun*.

## Usage

To run command *foo* (with its arguments) in the environment defined in *.env* file in current directory, type:

```bash
$ dotrun foo args...
```

*.env* file might define environment such as:

```bash
FOO=BAR
SPAM=EGGS
```

*foo* will then be able to access this environment defined in *.env* file.

You might also specify another dotenv file with *-env* option before the command to run:

```bash
$ dotrun -env /etc/foo.env foo args...
```

This way, *dotrun* won't load environment from *.env* file in current directory but from specified file instead.

*Enjoy!*
