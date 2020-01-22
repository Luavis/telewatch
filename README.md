# telewatch

Watch with telegram 🔭

## Install

```sh
$ git clone https://github.com/luavis/telewatch
$ cd telewatch
$ go build -o telewatch cmd/main.go
$ mkdir -p ~/.config/telewatch
$ echo "token: xxxxxx:XXXXXXXXXXX" > ~/.config/telewatch/token.yml
```

Telegram bot need to know chat id with you, so run `telewatch register` then
send any message to your bot.

```sh
$ telewatch register
Send any message to your bot, telewatch will save your chat id
Your telegram bot link: https://t.me/YOUR_BOT_NAME
Receive: hello world
Chat id: xxxxxx
Chat id is successfully saved
```

## Examples

- send message when changed
    ```
    $ telewatch command
    ```
- Refresh interval
    ```
    $ telewatch -n 60 command
    ```
- Daemonize
    ```
    $ telewatch -d command
    ```
- Alert when command done
    ```
    make && telewatch alert -m "Hello world"
    ```
