#!/bin/bash
OS="`uname`"
OS=$(echo $OS | tr '[A-Z]' '[a-z]')
curl -o /usr/local/bin/remem https://raw.githubusercontent.com/yehohanan7/remem/master/downloads/remem-${OS}
$chmod +x /usr/local/bin/remem
echo "remem s" >> ~/.profile

