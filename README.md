<!-----------------------------

- File Name : README.md

- Purpose :

- Creation Date : 12-15-2014

- Last Modified : Wed 17 Dec 2014 11:54:11 PM UTC

- Created By : Kiyor

------------------------------->

#	digger | dig more than ever

##	How to use

###	this require unbound installed in system

-	in centos, please install `unbound` `unbound-libs` `unbound-devel`
-	in OSX, please run `osx_install.sh`

```bash

go get -u github.com/kiyor/digger
go install github.com/kiyor/digger
cd ${GOPATH}/src/github.com/kiyor/digger
sudo cp -R ./reslov /usr/local/etc/
digger google.com

```

##	Sample output

![img](http://ccnacdn.s3.amazonaws.com/img/2014-12-15_README.md__notegosrcgithub.comkiyordigger_-_VIM__ssh__14144_11-39-47.png)

##	Option

-	use digger and [ip2loc](https://github.com/kiyor/ip2loc) will help you a lot
-	change or add `/usr/local/etc/reslov/xxx.conf` if you need
