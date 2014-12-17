#!/bin/bash
############################################

# File Name : osx_install.sh

# Purpose :

# Creation Date : 12-17-2014

# Last Modified : Wed 17 Dec 2014 11:49:42 PM UTC

# Created By : Kiyor 

############################################

openssl_ver='1.0.1j'
unbound_ver='1.5.1'

wget -N https://www.openssl.org/source/openssl-${openssl_ver}.tar.gz
tar xvzf openssl-${openssl_ver}.tar.gz
cd openssl-${openssl_ver}
./Configure darwin64-x86_64-cc
make
sudo make install


cd ..
wget -N https://unbound.net/downloads/unbound-${unbound_ver}.tar.gz
tar xvzf unbound-${unbound_ver}.tar.gz
cd unbound-${unbound_ver}
./configure
make
sudo make install

