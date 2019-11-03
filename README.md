# cronic

Interactive cron invocation tool

`cronic` lists the crontab files installed on the local machine.
Once the user chooses a specific crontab file, it shows the commands included inside the file.
A user can select a command and choose to execute the command.
The command will be run in a suitable environment, and as the specific user as specified in the crontab file.

# Building

Requires Go 1.13 or later.

```sh
go get
go build cmd/cronic/main.go
```

# Running

```
cmd/cronic/cronic [OPTIONS]
```

## Options

* `-path STRING`
    * path to crontabs, defaults to `CRONIC_PATH` environment variable, or `/etc/cron.d`

## Usage

The interface is split into three panes:

* to the left is a list of known crontab files
* to the top is a list of environment variable declarations in the selected crontab
* to the bottom is a list of commands contained in the selected crontab

The crontab files pane is focused when the application starts up.

To execute a specific crontab command:

1. Use the arrow keys to select the correct crontab file
2. Press enter to bring focs to the commands pane
3. Use the arrow keys to select the command you want to run
4. Press the enter key to run the command

The command is launched in the background.
All output is redirected to the system's log with a prefix of `cronic`.

When on the list of commands, press the `Esc` key to navigate back to the list
of files.

When on the list of files, press the `Esc` key to exit (or `^C` anywhere).

# Errata

## Temporary files

The script creates temporary files under `/tmp/`.
The filenames have the prefix `cronic`.

## Supported crontab formats

The application assumes all crontabs are in Vixie format, specifically:

1. the first 5 whitespace-delimited tokens on a line are the time specification

   Shortcuts such as (`@daily`) are also supported
2. the next token is the user under which the command is run
3. the remainder of the line is the command to be run

## Known issues

* Errors while launching the script cause the application to exit

  These are usually due to incorrect permissions.

* Percent (`%`) symbols in the commands are not un-escaped

* Only Vixie cron format is supported

## Future Improvements

* Cleaner UI + better shortcuts
