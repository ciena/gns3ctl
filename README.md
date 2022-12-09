# GNS3 Control Tool

This repository holds the source to a control tool that can be used to 
manipulate GNS3 networks.

## Creating a project/network

The `load` command is used to create a network as a GNS3 project. The input
to this command is a `YAML` formated file that contains the declaration of the
nodes and the links between those nodes. An example can be seen in the file
`example-network.yaml`.

## Configuration

The GNS3 command tool uses a configuration file that defaults to
`$HOME/.gns3ctl.yaml`. This file can be used to establish persistent flags
to the gns3ctl command, such as `project`, `compute`, or any of the _Global
Flags_.

## WIP - Work In Progress

This tool is very much a work in progress, so use the `--help` option to
discover what can be done and don't be surprised if things break. Also,
please feel free to submit issues and, better yet, sumbit a merge request.

## Clean GIT History

While this project has been under development internally for a while,
when it was made public in December 2022 it was done so without any
commit history. This was done out of an abundance of caution to ensure
nothing company propriary snuck in the history.

## Help

The help command can be used to get more information about the tool's
capabilities.

```
$ ./gns3ctl help
Allows a user to manipulate the GNS3 network simulator,
including the ability to create example networks and extract
information about those networks.

Usage:
  gns3ctl [command]

Available Commands:
  close       Closes subresources
  completion  Generate the autocompletion script for the specified shell
  delete      Deletes a subresource
  get         Fetch or query subresources
  help        Help about any command
  import      Import information form external systems
  load        Loads a project into the GNS3 environment
  open        Open a subresource
  start       Start the execution of subresources
  stop        Stop the execution of subresources
  suspend     Suspend a list of subresources
  version     Display the GNS3 server version

Flags:
  -a, --address string                Service and port on which to contact the server (default "localhost:3080")
  -d, --base-directory string         Default project name to use when performing project specific operations (default "/home/dbainbri/GNS3")
  -c, --compute string                Default compute to use when creating projects (default "local")
      --config string                 config file (default is $HOME/.gns3ctl.yaml)
      --download-buffer-size string   size of in memory buffer to use for file downloads (default "10M")
  -h, --help                          help for gns3ctl
  -k, --insecure-skip-verify          Skip verifing the TLS cert of the host (default true)
  -w, --password string               Password for basic authentication (default "admin")
  -p, --project string                Default project name to use when performing project specific operations (default "default")
  -t, --timeout duration              Timeout for http requests (default 20s)
  -u, --username string               Username for basic authentication (default "admin")

Use "gns3ctl [command] --help" for more information about a command.
```
