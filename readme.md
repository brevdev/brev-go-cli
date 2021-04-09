# Brev CLI

go version >= 1.16

## Distribute to Homebrew

Step 1: bump version (see top of Makefile)

Step 2: create homebrew distribution

```
> API_COTTER_KEY=... make dist-homebrew
```

Step 3: create GitHub release

Step 4: upload resultant tar.gz to GitHub release

Step 5: copy sha256 (output from step 2) and use it in a new update to https://github.com/brevdev/homebrew-tap

## Completions

To get completion to work:

brev completion bash > /usr/local/etc/bash_completion.d/brev

// this may not be needed
echo "autoload -U compinit; compinit" >> ~/.zshrc

brev completion zsh > "${fpath[1]}/\_brev"
