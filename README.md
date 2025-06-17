# gohome

#### How to use
- Build the `main.go` file or download the binary:
  ```
  go build -o .pontobin main.go && cp .pontobin ~/.pontobin
  ```
  or
  ```
  curl -L https://raw.githubusercontent.com/caio-farias/gohome/main/.pontobin -o ~/.pontobin
  chmod +x ~/.pontobin
  ```
- Create a filed called .ponto and then add shift info like:
  ```
  08:33 11:44
  15:17
  ```
- Add the following config to your `.bashrc`:
  ```
  _update_ponto_status() {
          PONTO_OUTPUT=$(~/.pontobin -no-prefix=true -show-final=true)
  }
  PROMPT_COMMAND="_update_ponto_status; $PROMPT_COMMAND"
  PS1='${PONTO_OUTPUT} \[\e[1;32m\]\u@\h\[\e[0m\]:\[\e[1;34m\]\w\[\e[0m\]\$ '

  ```
- Source from new `.bashrc` or reopen bash.
