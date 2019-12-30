# DotRun Command

DotRun command loads dotenv file, loads its environment and runs given command in that environment.

## Installation

### Unix users (Linux, BSDs and MacOSX)

Unix users may download and install latest *dotrun* release with command:

```bash
sh -c "$(curl http://sweetohm.net/dist/dotrun/install)"
```

If *curl* is not installed on you system, you might run:

```bash
sh -c "$(wget -O - http://sweetohm.net/dist/dotrun/install)"
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

*foo* will then be able to access this environment defined in *.env* file.

You can specify another dotenv file with `-env file` option before the command to run:

```bash
dotrun -env /etc/foo.env foo args...
```

You can also load multiple dotenv files, repeating `-env file` option on command line :

```bash
dotrun -env /etc/foo.env -env /etc/bar.env foo args...
```

The environment files are evaluated in the order of the command line, so that in previous example variables defined in *bar.env* would overwrite those defined in *foo.env*.

## Important Note

Let's say you have following *.env* file:

```bash
FOO=BAR
```

And run command `dotrun echo $FOO`. You would probably expect *BAR*, but this is not the case. Because `$FOO` will be evaluated before running this command and replaced with its value on command line. To have expected behavior, you must run:

```bash
$ dotrun 'echo $FOO'
BAR
```

In this case, `$FOO` won't be evaluated by the shell because it is protected by single quotes. Then *dotrun* will run command `echo $FOO` in a shell.

*Enjoy!*
