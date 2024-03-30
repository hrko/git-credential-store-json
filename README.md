# git-credential-store-json
Helper to store Git credentials on disk in JSON format.

## Installation
### Pre-built binaries
Pre-built binaries are available on the [releases page](https://github.com/hrko/git-credential-store-json/releases). Download the binary for your platform and place it in a directory that is in your `PATH`.

```sh
# Example for Linux (x86_64)
# Replace `<tag>` with the latest release version
wget https://github.com/hrko/git-credential-store-json/releases/download/<tag>/git-credential-store-json_linux-amd64.zip
unzip git-credential-store-json_linux-amd64.zip
sudo mv git-credential-store-json /usr/local/bin/
sudo chmod +x /usr/local/bin/git-credential-store-json
```

### From source
Alternatively, you can install the helper using `go get`:
```sh
go install github.com/hrko/git-credential-store-json
```

## Synopsis
```sh
git config credential.helper 'store-json [<options>]'
```

## Description
This helper stores credentials on disk in JSON format. It is useful when you want to store `oauth_refresh_token` or other properties that are not supported by the default `store` helper. Credentials are stored in plain text, so you should understand the potential security risks and make sure the file is only readable by the user.

## Options
- `-f <path>` - Path to the file where the credentials will be stored. Default: `~/.git-credentials.json`
- `-v` - Verbose mode.

## Examples
Use with git-credential-oauth:
```sh
git config --global --unset-all credential.helper
git config --global --add credential.helper store-json
git config --global --add credential.helper oauth
```
