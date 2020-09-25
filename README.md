# gitbrose

command line utility that is able to translate any local git file to the web page of its remote.

Supported types of remotes are:

* Github
* Gitlab
* AWS CodeCommit

## Installation

```shell
wget https://github.com/mhristof/gitbrowse/releases/latest/download/gitbrowse.$(uname) -O ~/bin/gitbrowse
chmod +x ~/bin/gitbrowse
```

## Usage

### Github with ssh remote

```
$ git config --get remote.origin.url
git@github.com:mhristof/gitbrowse.git
$ gitbrowse Makefile
https://github.com/mhristof/gitbrowse/blob/master/Makefile
```

## Vim

You can use `gitbrowse` from you Vim with this snippet

```vim
function GitBrowse()
    let line=line(".") + 1
    exec "silent !open $(gitbrowse " . expand('%') . " --line " . line . ")"
    exec ":redraw!"
endfunction
```

and if you want to abberviate it
```vim
cabbrev bb :call GitBrowse()<cr>
```
