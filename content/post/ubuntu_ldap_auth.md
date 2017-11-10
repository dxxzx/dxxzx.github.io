---
title: "ubuntu ldap auth"
date: 2017-11-10T11:34:11+08:00
draft: false
tags: [ubuntu,ldap]
topics: []
description: ""
---

# ubuntu网络认证

## 安装必要软件
```sh
sudo apt install libnss-ldap ldapscripts
```
某些配置条目可能如下：

	base dc=tianyisc,dc=com
	uri ldap://192.168.99.89/
	binddn cn=Manager,dc=tianyisc,dc=com
	bindpw password
	rootbinddn cn=Manager,dc=tianyisc,dc=com

Now configure the LDAP profile for NSS:
```sh
sudo auth-client-config -t nss -p lac_ldap
```
Configure the system to use LDAP for authentication:
```sh
sudo pam-auth-update
```
## 添加用户
### ldapscript（方式一）

> Install the package:
> ```sh
> sudo apt install ldapscripts
> ```
> Then edit the file /etc/ldapscripts/ldapscripts.conf to arrive at something similar to the following:
> ```
> SERVER=localhost
> BINDDN='cn=admin,dc=example,dc=com'
> BINDPWDFILE="/etc/ldapscripts/ldapscripts.passwd"
> SUFFIX='dc=example,dc=com'
> GSUFFIX='ou=Groups'
> USUFFIX='ou=People'
> MSUFFIX='ou=Computers'
> GIDSTART=10000
> UIDSTART=10000
> MIDSTART=10000
> ```
> Now, create the ldapscripts.passwd file to allow rootDN access to the directory:
> ```sh
> sudo sh -c "echo -n 'secret' > /etc/ldapscripts/ldapscripts.passwd"
> sudo chmod 400 /etc/ldapscripts/ldapscripts.passwd
> ```
> Replace “secret” with the actual password for your database's rootDN user.
> 
> The scripts are now ready to help manage your directory. Here are some examples of how to use them:
> 
> Create a new user:
> ```sh
> sudo ldapadduser george example
> ```
> This will create a user with uid george and set the user's primary group (gid) to example
> 
> Change a user's password:
> ```sh
> sudo ldapsetpasswd george
> Changing password for user uid=george,ou=People,dc=example,dc=com
> New Password: 
> New Password (verify):
> ```
> Delete a user:
> ```sh
> sudo ldapdeleteuser george
> ```
> Add a group:
> ```sh
> sudo ldapaddgroup qa
> ```
> Delete a group:
> ```sh
> sudo ldapdeletegroup qa
> ```
> Add a user to a group:
> ```sh
> sudo ldapaddusertogroup george qa
> ```
> You should now see a memberUid attribute for the qa group with a value of george.
> 
> Remove a user from a group:
> ```sh
> sudo ldapdeleteuserfromgroup george qa
> ```
> The memberUid attribute should now be removed from the qa group.
> 
> The ldapmodifyuser script allows you to add, remove, or replace a user's attributes. The script uses the same syntax as the ldapmodify utility. > For example:
> ```sh
> sudo ldapmodifyuser george
> # About to modify the following entry :
> dn: uid=george,ou=People,dc=example,dc=com
> objectClass: account
> objectClass: posixAccount
> cn: george
> uid: george
> uidNumber: 1001
> gidNumber: 1001
> homeDirectory: /home/george
> loginShell: /bin/bash
> gecos: george
> description: User account
> userPassword:: e1NTSEF9eXFsTFcyWlhwWkF1eGUybVdFWHZKRzJVMjFTSG9vcHk=
> 
> # Enter your modifications here, end with CTRL-D.
> dn: uid=george,ou=People,dc=example,dc=com
> replace: gecos
> gecos: George Carlin
> ```
> The user's gecos should now be “George Carlin”.

__这种方法添加的用户需要额外添加一些属性，尤其mail属性，这会在其他使用到__
添加mail属性
```sh
sudo ldapmodifyuser george
# About to modify the following entry :
dn: uid=george,ou=People,dc=example,dc=com
objectClass: account
objectClass: posixAccount
cn: george
uid: george
uidNumber: 1001
gidNumber: 1001
homeDirectory: /home/george
loginShell: /bin/bash
gecos: george
description: User account
userPassword:: e1NTSEF9eXFsTFcyWlhwWkF1eGUybVdFWHZKRzJVMj
# Enter your modifications here, end with CTRL-D.
dn: uid=george,ou=People,dc=example,dc=com
add: mail
mail: username@example.com
```

### ldapadd（方式二）
编辑文件add_user.ldif
```sh
cat > add_user.ldif << "EOF"
dn: uid=dylan,ou=People,dc=tianyisc,dc=com
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
uid: dylan
sn: Deng
givenName: Dylan
cn: Dylan Deng
displayName: Dylan Deng
uidNumber: 10000
gidNumber: 5000
gecos: Dylan Deng
loginShell: /bin/bash
homeDirectory: /home/dylan
mail: dengxingxian@tianyisc.com
EOF
```
添加用户
```sh
ldapadd -LLL -D cn=Manager,dc=tianyisc,dc=com -h 192.168.99.89 -b dc=tianyisc,dc=com -W -f add_user.ldif
```

__这样添加的用户仅存在于网络上，不能使用passwd修改密码完善方法如下：__

在/etc/passwd中添加对应的用户信息

	dylan:x:10000:5000:Dylan Deng:/home/dylan:/bin/bash


参考：
- [OpenLDAP Server](https://help.ubuntu.com/lts/serverguide/openldap-server.html#openldap-auth-config)
