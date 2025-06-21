# HTTP Inspector

Return Your HTTP Information.

## IP Data Source

[P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb)

## Deployment Guide

1. Prepare an environment to execute this program.
   ```shell
   useradd -r -s /usr/sbin/nologin -c "HTTP Inspector service user" http-inspector
   mkdir -p /opt/http-inspector
   chown http-inspector:http-inspector /opt/http-inspector
   chmod 750 /opt/http-inspector
   ```
2. Download the database from [ GeoLite2-City.mmdb ](https://github.com/P3TERX/GeoLite.mmdb/releases).
3. Place it to `/opt/http-inspector/ GeoLite2-City.mmdb`.
4. Download the executable file [releases](https://github.com/YogiLiu/http_inspector/releases).
5. Place the executable file to `/opt/http-inspector/http_inspector`.
6. Download [the systemd service file](https://github.com/YogiLiu/http_inspector/main/http-inspector.service).
7. Place it to `/etc/systemd/system/http-inspector.service`.
8. Start it.
   ```shell
   systemctl daemon-reload
   systemctl start http-inspector
   systemctl enable http-inspector
   ```

Then, the program will serve at `:5202`.