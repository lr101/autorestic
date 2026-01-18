<p align="center">
  <br>
  <br>
  <br>
  <img align="center" src="https://github.com/cupcakearmy/autorestic/raw/master/.github/logo.png" height="50" alt="autorestic logo">
  <br>
  <br>
  
  <p align="center">
    Config driven, easy backup cli for <a href="https://restic.net/">restic</a>.
    <br>
    <strong><a href="https://autorestic.vercel.app/">»»» Docs & Getting Started »»»</a></strong>
  <br><br>
  <a target="_blank" href="https://discord.gg/wS7RpYTYd2">
    <img src="https://img.shields.io/discord/252403122348097536" alt="discord badge" />
    <img src="https://img.shields.io/github/contributors/cupcakearmy/autorestic" alt="contributor badge" />
    <img src="https://img.shields.io/github/downloads/cupcakearmy/autorestic/total" alt="downloads badge" />
    <img src="https://img.shields.io/github/v/release/cupcakearmy/autorestic" alt="version badge" />
  </a>
  </p>
</p>


> **Fork Information:**
> <br>
> This is a modified fork of the original autorestic. It includes the following specific changes:
> * **InfluxDB Metrics:** Added integration for sending backup metrics and updates to InfluxDB. Add:
>```yml
>monitors:
>  stats:
>    type: influx
>    env:
>      server_tag:
>      influx_url:
>      influx_token:
>      influx_org:
>      influx_bucket:
>
>locations:
>  home:
>    from: 
>    <...>
>    monitors:
>      - stats
>
>```
> * **Docker Cron Functionality:** Added Docker cron scheduling capabilities including Python helper scripts.
>
>```yml
>services:
>  autorestic:
>    image: ghcr.io/lr101/autorestic:latest
>    container_name: autorestic
>    restart: unless-stopped
>    privileged: true                # (Optional)      
>    env_file: .env
>    environment:
>      - TZ=Europe/Berlin            # Change to your timezone
>      - CRON_SCHEDULE=10 2 * * *    # Cron schedule for backups
>    volumes:
>      - ${DEVICE_FOLDER_PATH}:/data          # path to .autorestic.yml
>      - ${LOCAL_BACKUP_PATH}:/backup
>      - ./logs:/var/log/autorestic           # (Optional) Logs
>      - ./rclone.conf:/root/.rclone.conf:ro  # (Optional) Rclone
>      - ~/.ssh:/root/.ssh:ro                 # (Optional) SFTP
>```

<br>
<br>

### 💭 Why / What?

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficult to manage if you have many different locations that you want to backup to multiple locations. This utility is aimed at making this easier 🙂.

### 🌈 Features

- YAML config files, no CLI
- Incremental -> Minimal space is used
- Backup locations to multiple backends
- Snapshot policies and pruning
- Fully encrypted
- Before/after backup hooks
- Exclude pattern/files
- Cron jobs for automatic backup
- Backup & Restore docker volume
- Generated completions for `[bash|zsh|fish|powershell]`
- **InfluxDB Metric Reporting** (Fork specific)
- **Docker Cron with Python Scripts** (Fork specific)

### ❓ Questions / Support

Check the [discussions page](https://github.com/cupcakearmy/autorestic/discussions) or [join on discord](https://discord.gg/wS7RpYTYd2)

## Contributing / Developing

PRs, feature requests, etc. are welcomed :)
Have a look at [the dev docs](./DEVELOPMENT.md)