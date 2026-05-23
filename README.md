# sshsite

My personal site, served over SSH. Connect from any terminal:

```
ssh 188.93.147.250
```

Built in Go with [Wish](https://github.com/charmbracelet/wish) and
[Bubble Tea](https://github.com/charmbracelet/bubbletea). Edit `content.go`
to change what's shown. Deployed on [Fly.io](https://fly.io).

## Run locally

```
./run.sh        # build + start on port 2222
ssh -p 2222 localhost
```

## Deploy

```
flyctl deploy --remote-only -a omer-sshsite
```
