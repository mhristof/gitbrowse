# gitbrose

command line utility that is able to translate any local git file to the web page of its remote.

Supported types of remotes are:

* Github
* Gitlab
* AWS CodeCommit

## Examples

### Github with ssh remote

```
$ git config --get remote.origin.url
git@github.com:mhristof/gitbrowse.git
$ gitbrowse Makefile
https://github.com/mhristof/gitbrowse/blob/master/Makefile
```
