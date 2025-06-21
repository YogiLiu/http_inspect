# HTTP Inspector

Return Your HTTP Information.

## IP Data Source

[lionsoul2014/ip2region](https://github.com/lionsoul2014/ip2region)

## Deployment Guide

1. Prepare an environment to execute this program.
   ```shell
   useradd -r -s /usr/sbin/nologin -c "HTTP Inspector service user" http-inspector
   mkdir -p /opt/http-inspector
   chown http-inspector:http-inspector /opt/http-inspector
   chmod 750 /opt/http-inspector
   ```
2. Download the database from [ip2region.xdb](https://github.com/lionsoul2014/ip2region/blob/master/data/ip2region.xdb).
3. Place it to `/opt/http-inspector/ip2region.xdb`.
4. Download the executable file [releases](https://github.com/YogiLiu/http_inspector/releases).
5. Decompress it and Place the executable file to `/opt/http-inspector/http_inspector`.
6. Download [the systemd service file](https://github.com/YogiLiu/http_inspector/main/http-inspector.service).
7. Place it to `/etc/systemd/system/http-inspector.service`.
8. Start it.
   ```shell
   systemctl daemon-reload
   systemctl start http-inspector
   systemctl enable http-inspector
   ```

Then, the program will serve at `:5202`.