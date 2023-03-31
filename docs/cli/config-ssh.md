<!-- DO NOT EDIT | GENERATED CONTENT -->

# config-ssh

Add an SSH Host entry for your workspaces "ssh coder.workspace"

## Usage

```console
coder config-ssh [flags]
```

## Description

```console
  - You can use -o (or --ssh-option) so set SSH options to be used for all your
    workspaces:

      $ coder config-ssh -o ForwardAgent=yes

  - You can use --dry-run (or -n) to see the changes that would be made:

      $ coder config-ssh --dry-run
```

## Options

### -n, --dry-run

|             |                                 |
| ----------- | ------------------------------- |
| Type        | <code>bool</code>               |
| Environment | <code>$CODER_SSH_DRY_RUN</code> |

Perform a trial run with no changes made, showing a diff at the end.

### --ssh-config-file

|             |                                     |
| ----------- | ----------------------------------- |
| Type        | <code>string</code>                 |
| Environment | <code>$CODER_SSH_CONFIG_FILE</code> |
| Default     | <code>~/.ssh/config</code>          |

Specifies the path to an SSH config.

### --ssh-host-prefix

|      |                     |
| ---- | ------------------- |
| Type | <code>string</code> |

Override the default host prefix.

### -o, --ssh-option

|             |                                     |
| ----------- | ----------------------------------- |
| Type        | <code>string-array</code>           |
| Environment | <code>$CODER_SSH_CONFIG_OPTS</code> |

Specifies additional SSH options to embed in each host stanza.

### --use-previous-options

|             |                                              |
| ----------- | -------------------------------------------- |
| Type        | <code>bool</code>                            |
| Environment | <code>$CODER_SSH_USE_PREVIOUS_OPTIONS</code> |

Specifies whether or not to keep options from previous run of config-ssh.

### -y, --yes

|      |                   |
| ---- | ----------------- |
| Type | <code>bool</code> |

Bypass prompts.