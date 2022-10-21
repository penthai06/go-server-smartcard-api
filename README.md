# go-server-smartcard

## go build to product

```
export GOOS=linux
export GOARCH=amd64
go build .
```

## แหล่งข้อมูลการทำ swagger

https://github.com/Napat/sscard
https://www.facebook.com/groups/584867114995854/posts/1101152216700672/

ไทย
https://github.com/gogetth/sscard/blob/master/main/example_thidcard.go

eng
https://github.com/ebfe/scard/blob/master/example_test.go
https://github.com/sf1/go-card/blob/master/examples/apdu/main.go

printer
https://github.com/augustopimenta/escpos

Install package on ubuntu lite
(drive smartcard)[https://www.cnx-software.com/2019/08/11/reading-id-card-data-in-ubuntu-with-ez100pu-smart-card-reader-thai-id-edition/]
`sudo apt install libccid pcscd pinentry-gtk2 pcsc-tools libpcsclite-dev libreadline-dev coolkey`

(CUPS Print Server)[https://linuxhint.com/cups_print_server_ubuntu/]
`sudo apt install cups –y`
