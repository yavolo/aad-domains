# aad-domains

Fetches domains related to an AAD tenant

Install

`go get -u github.com/yavolo/aad-domains`

Usage

`cat domains.txt | aad-domains`

Example

```console
yavolo@box:~/go/bin$ cat domains.txt
microsoft.com
yavolo@box:~/go/bin$ cat domains.txt | ./aad-domains
css.one.microsoft.com
dmarc.microsoft
cloudyn.com
cyberx-labs.com
volometrix.com
...
```
