autoload -U add-zsh-hook

export ENVY_SHELL=zsh
export ENVY_SESSION_KEY={{.SessionKey}}

envy_chpwd_hook() {
  # execute the undo script if it exists
  if [[ -f {{.UndoFilepath}} ]]; then
    . {{.UndoFilepath}}
  fi

  # generate the shell scripts then execute the load script
  envy gen
  . {{.LoadFilepath}}
}

add-zsh-hook chpwd envy_chpwd_hook

envy_zshexit_hook() {
  # remove all envy files for this session
  rm {{.LoadFilepath}}
  rm {{.UndoFilepath}}
}

add-zsh-hook zshexit envy_zshexit_hook

envy_chpwd_hook