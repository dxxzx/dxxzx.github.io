---
title: "gerrit service establish guide"
date: 2017-11-10T11:34:11+08:00
draft: true
tags: [gerrit]
topics: []
description: ""
---

# Gerrit服务器搭建
安装gerrit之前请确保mariadb服务和openldap服务已经安装并配置好，如果未安装，参考这里[mariadb服务器搭建](mariadb_server.md)，[openldap服务器搭建](openldap_service.md)
## 配置gerrit
## 配置tomcat
```
[gerrit]
	basePath = /var/lib/gerrit
	serverId = 408f5af4-fcd0-4e39-823a-5b83a514e542
	canonicalWebUrl = http://192.168.99.89:8080/gerrit/
[database]
	type = mariadb
	database = reviewdb
	hostname = localhost
	username = gerrit
[index]
	type = LUCENE
[download]
	command = checkout
	command = cherry_pick
	command = pull
	command = format_patch
	scheme = ssh
[auth]
	type = LDAP
	gitBasicAuthPolicy = LDAP
[gc]
	startTime = 1:00
	interval = 1 w
[gitweb]
	project = ?p=${project}
	branch = ?p=${project}
	revision = ?p=${project}
	filehistory = ?p=${project}
	roottree = ?p=${project}
	file = ?p=${project}
[changeCleanup]
	startTime = 4:00
	interval = 1 w
[ldap]
	server = ldap://127.0.0.1
	username = cn=Manager,dc=tianyisc,dc=com
	accountBase = ou=People,dc=tianyisc,dc=com
	accountPattern = (&(objectClass=person)(uid=${username}))
	accountFullName = displayName
	accountEmailAddress = mail
	groupBase = ou=Groups,dc=tianyisc,dc=com
	groupMemberPattern = (&(objectClass=group)(member=${dn}))
[sendemail]
    smtpServer = smtp.mxhichina.com
    smtpSeverPort = 25
    smtpEncryption = ssl
    smtpUser = iptv@tianyisc.com
    sslVerify = false
    from = Gerrit Review <iptv@tianyisc.com>
[receive]
	enableSignedPush = true
[container]
	user = tomcat
	javaHome = /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.144-0.b01.el7_4.x86_64/jre
[automerge]
	botEmail = liming@tianyisc.com
[commentlink "change"]
    match = "#/c/(\\d+)"
    html = "<a href=\"/#/c/$1/\">$1</a>"
[sshd]
	listenAddress = *:29418
[httpd]
	listenUrl = http://*:8080
[cache]
	directory = /var/cache/gerrit
[commitmessage]
	rejectTooLong = true
[plugins]
	allowRemoteAdmin = true
```