---
title: "openldap service establish guide"
date: 2017-11-10T11:34:11+08:00
draft: false
tags: [openldap,ldap]
description: ""
---

# establish openldap server
## Install and start service
```sh
yum install openldap-servers openldap-clients
cp /usr/share/openldap-servers/DB_CONFIG.example /var/lib/ldap/DB_CONFIG
chown ldap. /var/lib/ldap/DB_CONFIG
systemctl start slapd
systemctl enable slapd
```
## setup OpenLDAP manager password
generate encrptyed password:
```sh
# slappasswd    
New password:   
Re-enter new password:   
{SSHA}2aaO8Jrm2AkRYmI8dMptxesNsQ9bI2y8
```

string {SSHA}xxxxxxxxxxxxxxxxxxxxxxxx are encrypted password, it will be used later. 
then, create file like below. 

```sh
cat > chrootpw.ldif << "EOF"
dn: olcDatabase={0}config,cn=config
changetype: modify
add: olcRootPW
olcRootPW: {SSHA}2aaO8Jrm2AkRYmI8dMptxesNsQ9bI2y8
EOF
```
import this file:
```sh
# ldapadd -Y EXTERNAL -H ldapi:/// -f chrootpw.ldif  
SASL/EXTERNAL authentication started  
SASL username: gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth  
SASL SSF: 0  
modifying entry "olcDatabase={0}config,cn=config"
```

## 导入基本 Schema（可以有选择的导入）

```sh
cd /etc/openldap/schema/  
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f cosine.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f nis.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f collective.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f corba.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f core.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f duaconf.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f dyngroup.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f inetorgperson.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f java.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f misc.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f openldap.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f pmi.ldif
ldapadd -Y EXTERNAL -H ldapi:/// -D "cn=config" -f ppolicy.ldif
```

## 设置自己的Domain Name
首先要生成经处理后的目录管理者明文密码：
```sh
# slappasswd  
New password:   
Re-enter new password:   
{SSHA}2aaO8Jrm2AkRYmI8dMptxesNsQ9bI2y8
```
之后，再新建如下文件，文件内容如下，注意，要使用你自己的域名替换掉文件中所有的 "dc=***,dc=***"，并且使用刚刚生成的密码，替换文中的 "olcRootPW" 部分： 
```sh
cat > chdomain.ldif << "EOF"
# replace to your own domain name for "dc=***,dc=***" section  
# specify the password generated above for "olcRootPW" section  
dn: olcDatabase={1}monitor,cn=config  
changetype: modify  
replace: olcAccess  
olcAccess: {0}to * by dn.base="gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth"  
  read by dn.base="cn=Manager,dc=tianyisc,dc=com" read by * none  
  
dn: olcDatabase={2}hdb,cn=config  
changetype: modify  
replace: olcSuffix  
olcSuffix: dc=tianyisc,dc=com  
  
dn: olcDatabase={2}hdb,cn=config  
changetype: modify  
replace: olcRootDN  
olcRootDN: cn=Manager,dc=tianyisc,dc=com  
  
dn: olcDatabase={2}hdb,cn=config  
changetype: modify  
add: olcRootPW  
olcRootPW: {SSHA}ZhmO2UeH4tsyy5ly0fTwdkO10WJ69V6U  
  
dn: olcDatabase={2}hdb,cn=config  
changetype: modify  
add: olcAccess  
olcAccess: {0}to attrs=userPassword,shadowLastChange by  
  dn="cn=Manager,dc=tianyisc,dc=com" write by anonymous auth by self write by * none  
olcAccess: {1}to dn.base="" by * read  
olcAccess: {2}to * by dn="cn=Manager,dc=tianyisc,dc=com" write by * read
EOF
```
之后再导入该文件：
```sh
# ldapmodify -Y EXTERNAL -H ldapi:/// -f chdomain.ldif    
SASL/EXTERNAL authentication started  
SASL username: gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth  
SASL SSF: 0  
modifying entry "olcDatabase={1}monitor,cn=config"  
  
modifying entry "olcDatabase={2}hdb,cn=config"  
  
modifying entry "olcDatabase={2}hdb,cn=config"  
  
modifying entry "olcDatabase={2}hdb,cn=config"  
  
modifying entry "olcDatabase={2}hdb,cn=config"  
```
然后再新建如下文件,  文件内容如下，注意，要使用你自己的域名替换掉文件中所有的 "dc=***,dc=***"：
```sh
cat > basedomain.ldif << "EOF"
dn: dc=tianyisc,dc=com
objectClass: top
objectClass: dcObject
objectClass: organization
o: TianYi
dc: tianyisc

dn: ou=People,dc=tianyisc,dc=com
objectClass: organizationalUnit
ou: People

dn: ou=Groups,dc=tianyisc,dc=com
objectClass: organizationalUnit
ou: Groups

dn: cn=developers,ou=Groups,dc=tianyisc,dc=com
objectClass: posixGroup
cn: developers
gidNumber: 5000
EOF
```
最后导入该文件：
```sh
# ldapadd -x -D cn=Manager,dc=tianyisc,dc=com -W -f basedomain.ldif  
Enter LDAP Password:   
adding new entry "dc=tianyisc,dc=com"  
  
adding new entry "cn=Manager,dc=tianyisc,dc=com"  
  
adding new entry "ou=People,dc=tianyisc,dc=com"  
  
adding new entry "ou=Group,dc=tianyisc,dc=com"  
```
## 允许防火墙访问 LDAP 服务
```sh
firewall-cmd --zone=internal --add-source=192.168.99.0/24 --permanent   # 使用internal区域，并将ip范围加入internal区域
firewall-cmd --zone=internal --add-service=ldap --permanent             # 允许ldap服务
firewall-cmd --reload                                                   # 重新加载防火墙规则
```


references:

- [\[原创\] CentOS7 下 OpenLDAP Server 安装和配置及使用 phpLDAPadmin 和 Java LDAP 访问 LDAP Server - Leo's Blog - ITeye博客](http://yhz61010.iteye.com/blog/2352672)
