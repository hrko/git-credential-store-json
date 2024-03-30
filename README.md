# git-credential-store-json
Helper to store Git credentials on disk in JSON format.

## Installation
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
