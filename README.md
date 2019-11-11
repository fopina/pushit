# pushit

Similar to [notify-push](https://github.com/fopina/notify-push) (in python), pushit is a CLI tool to push notifications straight to your phone.

The usual sysadmin method is to send alerts via email but there are tons of services out there with mobile apps that allow us to send push notifications, why not use those for faster and cleaners alerts?

## Instalation


Use `go get`:

```
go get github.com/fopina/pushit
```

Or download a pre-built binary from [releases](https://github.com/fopina/pushit/releases).  
The pre-built binary has the `--update` flag to allow easier updates in the future.

# Usage

```bash
pushit -h
Usage of pushit:
  -c, --conf string      TOML configuration file (default "~/.pushit.conf")
  -h, --help             this
  -l, --lines int        number of lines of the input that will be pushed - ignored if --stream is used (default 10)
  -o, --output           echo input - very useful when piping commands
  -p, --profile string   profile to use
  -s, --stream           stream the output, sending each line in separate notification
  -u, --update           Auto-update pushit with latest release
  -v, --version          display version
```

The configuration file is required but an example is available [here](https://github.com/fopina/pushit/blob/master/pushit.conf.example) in the repo.

```yaml
default: pushitbot-demo

profiles:
  pushitbot-demo:
    service: pushitbot
    params:
      token: 105e48ff92b92263f3397ed55f275a81
      format: Markdown
```

* `default` is the default profile when `--profile` is not specified
* `profiles` is a mapping with configured profiles
  * `pushitbot-demo` is the profile name, free choice
  * `service` is the notification service, check below for options
  * `params` is a mapping of options specific to the service, check below for options

Pushing messages is now as simple as:

```bash
$ pushit hello world
```

Or using `stdin` (useful for shell piping):

```bash
$ pushit
hello
$ echo hello world | pushit
```

A flag worth highlighting is `--stream` that will post a message per line read from `stdin`, just try:

```bash
$ (echo one; sleep 1; echo done) | pushit --stream
```

### Services

`Params` are the possible values to use in configuration file `params` profile entry

| Service    | Description | Params |
| ---------- | ------------- | -------------|
| slack      | Use [Slack Incoming Webhooks](https://api.slack.com/messaging/webhooks)  | **url**: webhook URL *(required)*<br>**channel**: channel to post the message to<br>**username**: username to display<br>**icon_emoji**: emoji to use as bot picture<br>**icon_url**: URL to use as bot picture
| pushitbot  | Use [PushItBot](http://fopina.github.io/tgbot-pushitbot/) | **token**: token provided by @PushItBot *(required)*<br>**format**: blank (default), HTML or Markdown, as defined in the service documentation
| telegram   |  Use [Telegram Bot API](https://core.telegram.org/bots/api#sendmessage) | **token**: token provided by @BotFather *(required)*<br>**chat_id**: target user/group ID - use @myidbot to find your ID or a group ID *(required)*
| requestbin | Use [Requestbin.com](https://requestbin.com/) - demo service for testing/debugging | **url**: requestbin.com generated endpoint<br>Any other params defined will be posted to the endpoint as well

