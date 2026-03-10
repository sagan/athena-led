- [athena-led](#athena-led)
  - [安装方法](#安装方法)
    - [针对官方固件](#针对官方固件)
    - [针对其他固件](#针对其他固件)
  - [按键控制](#按键控制)
  - [使用示例](#使用示例)
  - [构建](#构建)

# athena-led

自用的控制京东云 AX6600 雅典娜路由器 (JDCloud AX6600 athena, RE-CS-02) 的 LED
点阵显示屏的用户态程序。这是对[原项目](https://github.com/NONGFAH/athena-led)的 fork。加入一些新特性，包括：

- 增加支持显示的字符数量。现在所有 ASCII 字符都能够正常显示。也支持部分特殊 Unicode 符号，包括：
  - ♥ (heart), ☀ (sunny), ☾ (moon), ☁ (cloudy), ⛆ (little rain), 🌧 (rainy), ⛈ (thunderstorm), ❄ (snow), 🌫 (fog)
  -  ←, →, ↑, ↓, ↗, ↘, ✓, ✗
- 默认写入 pid 文件 `/var/run/athena-led.pid`。
- `-backend` 参数，用于选择屏幕驱动后端：
  - `auto`（默认）：启动时优先使用 `/dev/tmp1628-led`，失败则回退 GPIO。
  - `tmp1628`：使用 `/dev/tmp1628-led`，设备不存在时启动报错。
  - `gpio`：使用原来的 GPIO 方案。
- `-option` flag 增加几种可选显示内容：
  - cpu : CPU 占用率。
  - mem : 内存占用率。
  - upload : 当前上传速度。
  - download : 当前下载速度。
  - updl : 当前的上传和下载速度。上传数据在前。最大程度精简了格式，能够完整显示在屏幕上而无需滚动。
  - nic : 当前各个物理网口的链路状态。网口显示顺序为 `wan,lan4,lan3,lan2,lan1`。每个网口的链路状态用一个字母表示：`O` - disconnected; `B` - 10 Mbps; `H` - 100 Mbps; `G` - 1000 Mbps; `S` - 2500 Mbps。
  - countdown : 显示倒计时计数。
  - dino : 显示 [Chrome Dino](https://en.wikipedia.org/wiki/Dinosaur_Game) 风格的动画。
- `-option` flag 部分原有选项也有修改：
  - date : 可以自定义日期显示格式（Go 时间格式字符串）。支持两种自定义的代表星期几的时间格式符：`%w` - 使用点数(dots)表示；`%W` - 使用罗马数字表示。
  - time : 默认时间格式改为 `15:04:05`，显示秒数。同样可以自定义时间显示格式。
- 点阵屏幕两侧的4个 LED 状态灯(status)现在会动态变化：
  - time : 表示系统 CPU 占用率。占用率越高闪烁越快。
  - medal : 表示当前路由器 Internet 连接是否正常，如果正常则亮。默认使用的测试 url: `http://www.google.com/generate_204`。
  - upload & download : 表示当前路由器网口的实时上传/下载状态。传输速率越大闪烁越快。默认显示 Internet 出口网口的信息。
- `-option` flag 里（空格分隔）的每个 option 支持独立设置内容(`:content` 格式后缀) 和/或 显示时长(`#5` 格式后缀)参数，例如：
  - `string:i_have_a_dream` : 显示 "i have a dream" 文字 (将内容里的空格替换为 _)。
  - `url:https://ipinfo.io/json#3` : 显示 `https://ipinfo.io/json` 这个 URL 的内容，显示时长为3秒。
  - `upload:lan1` : 显示 `lan1` 这个网络接口的上传速度。
  - `date:01-02 %W` : 显示 `01-02` 格式的日期以及代表当前是星期几的罗马字符(Go 时间格式字符串)。
  - `time:15_04` : 显示 "15 04" 格式的时间。
  - `temp:24` : 显示区域 2 (`/sys/class/thermal/thermal_zone2`) 和 4 (`/sys/class/thermal/thermal_zone4`) 的温度。
  - `dino#15` : 显示 15秒的恐龙动画。
- 支持重复传入多个 `-option` 参数。每个参数作为一个 Profile。程序启动后默认使用第一个 Profile。
- 支持通过 signal 信号控制程序。例如：`kill -SIGUSR1 $(cat /var/run/athena-led.pid )`。程序启动后会写入该 pid 文件。
  - `SIGUSR1` : 切换使用的 Profile。
  - `SIGUSR2` : 切换屏幕的关闭 / 打开状态。
  - `SIGHUP` : 打开屏幕 / 重置为使用默认 Profile。
- 支持通过 `TZ` 环境变量修改显示的日期/时间的时区。
- 通过 url 获取的显示内容默认缓存至少 60 秒；支持通过 `Cache-Control` 响应头设置缓存有效期。

## 安装方法
### 针对官方固件
由于官方固件默认目录空间有限，下载 `athena-led` (arm版本) 可执行文件后，如果`/usr/sbin/`空间不足，也可以放到 `/opt/bin` 等其他合适的目录下（需要自行增加 `PATH` 目录）。

此外，官方固件的 `/etc/init.d/jdcloudbi` 会直接控制 LED 显示屏。因此需要禁用该服务，否则会和本程序冲突。

```sh
/etc/init.d/jdcloudbi disable
/etc/init.d/jdcloudbi stop
```

最后，为了可以让 `athena-led` 开机启动，需要将其作为启动服务，创建 `/etc/init.d/athena_led` 文件并设置可执行权限 `chmod a+x /etc/init.d/athena_led`。参考内容如下：

```sh
#!/bin/sh /etc/rc.common

START=99
STOP=10

USE_PROCD=1
NAME=athena_led
PIDFILE=/var/run/athena-led.pid
PROG=/path/to/athena-led       # 此处替换成实际路径

start_service() {
    procd_open_instance
    procd_set_param command "$PROG" -backend auto -pidFile "$PIDFILE" -seconds 5 -option time # 可根据需要修改启动参数
    procd_set_param respawn 3600 5 0
    procd_set_param stdout 1
    procd_set_param stderr 1
    procd_close_instance
}

stop_service() {
    service_stop "$PROG"
}

reload_service() {
    stop
    start
}
```

其中 `PROG=` 需要先改成 `athena-led` 可执行文件的实际路径，例如 `/usr/sbin/athena-led` 或 `/opt/bin/athena-led`。默认启动参数会持续显示时间；如果你需要自定义显示内容，直接修改脚本里的 `procd_set_param command` 那一行即可。

写入脚本后执行：

```sh
/etc/init.d/athena_led enable
/etc/init.d/athena_led start
```

### 针对其他固件
下载 `athena-led` (arm64版本)可执行文件然后放到 `/usr/sbin/athena-led`
(替换原文件)即可。本程序命令行参数与原版保持兼容，所以仍然可以用
[luci-app-athena-led](https://github.com/NONGFAH/luci-app-athena-led) 控制。但一些新特性需要手动编辑
`/etc/config/athena_led` 和/或 `/etc/init.d/athena_led` 文件里的参数才能启用，无法通过 luci / Web UI 设置。

运行 `athena-led -h` 查看详细帮助。

## 按键控制

雅典娜路由器顶部有2个圆形物理按键：右侧的（较大的）按键是 `wps` 键。左侧的（较小的）按键是 `BTN_0` 键，在 OpenWrt
固件里默认没有功能。通过配合使用本程序的多 Profile 功能，可以将 `BTN_0` 键配置为切换屏幕显示内容。参考 OpenWrt
的[文档](https://openwrt.org/docs/guide-user/hardware/hardware.button)。在
[这个](https://github.com/ZqinKing/wrt_release/releases) ImmortalWrt 固件中测试工作，其它固件未测试。

创建 `/etc/rc.button/BTN_0` 文件并设置可执行权限 `chmod a+x /etc/rc.button/BTN_0`。内容如下：

```sh
#!/bin/sh

# pressed / timeout / released
if [ "$ACTION" = "pressed" ]
then
    kill -SIGUSR1 $(cat /var/run/athena-led.pid )
    return 2
elif [ "$ACTION" = "timeout" ]
then
    kill -SIGUSR2 $(cat /var/run/athena-led.pid )
fi

```

功能：

- 短按 `BTN_0` 按键：切换屏幕显示内容。
- 长按 `BTN_0` 按键 2 秒：关闭或打开屏幕。注意长按触发 `timeout` 事件的所需秒数由脚本里 `pressed` 事件的返回值决定。

## 使用示例

修改 `/etc/init.d/athena_led`, 将 `procd_set_param` 行内容改为：

```
procd_set_param command $PROG -seconds 5 -option "date#2 time:15:04#2 url upload#2 cpu#2 mem#2" -option time -option upload -url "https://ipinfo.io/ip" -tempFlag "4"
```

以上传入了3个 `-option` 参数设置了 3 个 Profile:

- Profile 0 (默认): 显示2秒日期、然后显示2秒 15:04 格式的时间、显示5秒 https://ipinfo.io/ip 内容、显示2秒上传速度、显示2秒CPU占用率、最后显示2秒内存占用率。然后循环回到开始。
- Profile 1 : 一直显示时间。
- Profile 2 : 一直显示上传速度。

默认使用第一个 Profile。通过设备顶部的按键切换其他 Profile。也可以自己写脚本发送 signal 控制。

修改后需要 `service athena_led restart` 重启服务。

## 构建

1. 设置环境变量
GOARCH=arm64;GOOS=linux
2. 编译
go build -ldflags="-s -w" -trimpath -o athena-led
3. 打包
复制`athena-led` 到`luci-app-athena-led/root/usr/sbin/`
