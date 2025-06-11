# gohome
bash integrated shift calculator

# how to use
- Download the binary using: 
  ```
  curl -o ~/.pontobin URL_TO_BIN
  ```
- Create a filed called .ponto and then add shift info like:
  ```
  08:33 11:44
  15:17
  ```
- Add the following config to your `.bashrc`:
  ```
  _update_ponto_status() {
          PONTO_STATUS=$(~/.pontobin)
  }
  PROMPT_COMMAND="_update_ponto_status; $PROMPT_COMMAND"
  PS1='\[\e[1;36m\]${PONTO_STATUS}\[\e[0m\] \[\e[1;32m\]\u@\h\[\e[0m\]:\[\e[1;34m\]\w\[\e[0m\]\$ '
  ```
- Source from new `.bashrc` or repone bash