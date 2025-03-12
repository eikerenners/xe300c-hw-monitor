## Simple Monitor application for the iNET XE300C 4G Wifi Router 
This application uses the internal API of the router available on localhost. Currently, only one call and one response object is supported, though it should be fairly easy to extend. As a result, the current temperature, battery level and load average is printed to the console and the "status.log" file in the same directory as the application. 

Note: The Rust version is just a skeleton; I haven't gotten around to actually implementing anything...

Request: 
```sh
curl -k http://127.0.0.1/rpc -H 'glinet: 1' -d '{"jsonrpc":"2.0","id":1,"method":"call","params":["","system","get_info"]}'
```

Response: 
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "mac": "94:83:C4:4E:DE:FD",
    "disable_guest_during_scan_wifi": false,
    "hardware_version": "",
    "country_code": "",
    "sn_bak": "3b0785ef83b61dca",
    "software_feature": {
      "ipv6": true,
      "adguard": false,
      "passthrough": false,
      "repeater_eap": true,
      "vpn": true,
      "ids_ips": false,
      "bark": false,
      "tor": false,
      "secondwan": false,
      "sms_forward": true,
      "nas": false
    },
    "vendor": "GL.iNet",
    "hardware_feature": {
      "reset_button": "gpio-3",
      "nand": true,
      "bluetooth": false,
      "wan": "eth1",
      "usb_reset": "",
      "switch_button": "",
      "radio": "mac80211",
      "lan": "eth0",
      "usb": "1-1.3",
      "build_in_modem": "1-1.2",
      "noled": false,
      "hwnat": false,
      "microsd": "1-1.1",
      "modem_reset": 2,
      "fan": false,
      "mcu": true,
      "nowds": false
    },
    "cpu_num": 1,
    "board_info": {
      "architecture": "Qualcomm Atheros QCA9533 ver 2 rev 0",
      "hostname": "GL-XE300",
      "kernel_version": "5.10.176",
      "openwrt_version": "OpenWrt 22.03.4 r20123-38ccc47687",
      "model": "GL.iNet GL-XE300(NOR\\/NAND)"
    },
    "firmware_date": "2024-08-23 12:28:02",
    "model": "xe300",
    "ddns": "sfedefd",
    "sn": "20526461fe60b24e",
    "firmware_type": "release1",
    "firmware_version": "4.3.18"
  }
}
```

## Build GO App for MIPS 
Since the XE300C is a MIPS based microcontroller architecture, the golang compiler needs to be instructed to compile the application accordingly. 

Building for GL.iNet gl-xe300. Find architecture of the device: 
```sh
uname -a

Linux 
GL-XE300 
5.10.176 
#0
Sun Apr 9 12:27:46 2023 
mips 
GNU/Linux
```

To build app for wrt xe300c (mips architecture)
```sh
GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o xe300c-hw-monitor 
```

## install as service
SSH into the router, install SCP, copy files via SCP, create new system service, start up manually, run for a while, stop. 
Remember, the router has VERY limited resources. Don't run this service too long and clutter the internal memory. 

```sh
ssh root@192.168.8.1

opkg update
opkg install openssh-sftp-server
exit

GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o xe300c-hw-monitor
scp -r ./xe300c-hw-monitor root@192.168.8.1:/tmp
scp -r ./xe300c-hw-monitor root@192.168.8.1:/tmp

```