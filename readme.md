To get completion to work:

brev completion bash > /usr/local/etc/bash_completion.d/brev

// this may not be needed
echo "autoload -U compinit; compinit" >> ~/.zshrc

brev completion zsh > "${fpath[1]}/\_brev"
