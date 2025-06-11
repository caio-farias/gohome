# gohome

**bash-integrated shift tracker** (8h shift by default)

## Setup

1. **Install the binary**:

   ```bash
   curl -L https://raw.githubusercontent.com/caio-farias/gohome/main/.pontobin -o ~/.pontobin
   chmod +x ~/.pontobin
   ```

2. **Create a `.ponto` file** in your user home:

   ```
   08:33 11:44
   15:17
   ```

   * Each line must be in a `ENTRY EXIT` or `ENTRY ` format.

3. **Update your `.bashrc`** to display shift status in your prompt:

   ```bash
   _update_ponto_status() {
     PONTO_STATUS=$(~/.pontobin)
   }
   PROMPT_COMMAND="_update_ponto_status; $PROMPT_COMMAND"
   PS1='\[\e[1;36m\]${PONTO_STATUS}\[\e[0m\] \[\e[1;32m\]\u@\h\[\e[0m\]:\[\e[1;34m\]\w\[\e[0m\]\$ '
   ```

4. **Reload your shell**:

   ```bash
   source ~/.bashrc
   ```
   
## Preview

* **Time remaining** (e.g., 1h15m left; `r-` prefix):
  ![Remaining Time](https://github.com/user-attachments/assets/c47de973-8ec8-4fa9-97e2-82d962d6acbd)

* **Shift exceeded** (e.g., 1m over; `w-` prefix):
  ![Exceeded Time](https://github.com/user-attachments/assets/e441d2b7-e0c9-4bff-8daf-882b8fd7a64a)
