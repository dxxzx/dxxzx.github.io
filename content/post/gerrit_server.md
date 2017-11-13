---
title: "gerrit service establish guide"
date: 2017-11-10T11:34:11+08:00
draft: false
tags: [gerrit]
description: ""
---

# Gerrit server establish
Before install gerrit, please ensure mariadb service and openldap service are installed and configured. If not yet, please refer [mariadb server establish](mariadb_server.md), [openldap server establish](openldap_service.md).
## init gerrit site
using below command to init gerrit site
```sh
java -jar gerrit.war init -d gerrit-site
```
output may look like below:
```
Using secure store: com.google.gerrit.server.securestore.DefaultSecureStore
[2017-11-10 21:25:18,220] [main] INFO  com.google.gerrit.server.config.GerritServerConfigProvider : No /home/dylan/gerrit-site/etc/gerrit.config; assuming defaults

*** Gerrit Code Review 2.14.4
*** 

Create '/home/dylan/gerrit-site' [Y/n]? y
```
### git repository path section
```
*** Git Repositories
*** 

Location of Git repositories   [git]: 
```
### database section
before answer gerrit, you should create review database and mariadb for gerrit.
execute below command in mysql prompt.
```sql
CREATE USER 'gerrit'@'localhost' IDENTIFIED BY 'secret';
CREATE DATABASE reviewdb DEFAULT CHARACTER SET 'utf8';
GRANT ALL ON reviewdb.* TO 'gerrit'@'localhost';
FLUSH PRIVILEGES;
```
then, anwser question in gerrit installation
```
*** SQL Database
*** 

Database server type           [h2]: mariadb

Gerrit Code Review is not shipped with MariaDB Connector/J 1.5.9
**  This library is required for your configuration. **
Download and install it now [Y/n]? y
Downloading https://repo1.maven.org/maven2/org/mariadb/jdbc/mariadb-java-client/1.5.9/mariadb-java-client-1.5.9.jar ... OK
Checksum mariadb-java-client-1.5.9.jar OK
Server hostname                [localhost]: 
Server port                    [(mariadb default)]: 
Database name                  [reviewdb]: 
Database username              [dylan]: gerrit
gerrit's password              : 
              confirm password : 
```
### index type
```
*** Index
*** 

Type                           [lucene/?]: 
```
### authentication section
before anwser these question, please setup ldap service first. see [openldap_service](openldap_service)
```
*** User Authentication
*** 

Authentication method          [openid/?]: ldap
Git/HTTP authentication        [http/?]: ldap
LDAP server                    [ldap://localhost]: 
LDAP username                  : cn=Manager,dc=my-domain,dc=com
cn=Manager,dc=my-domain,dc=com's password : 
              confirm password : 
Account BaseDN                 : ou=People,dc=my-domain,dc=com
Group BaseDN                   [ou=People,dc=my-domain,dc=com]: ou=Groups,dc=my-domain,dc=com
Enable signed push support     [y/N]? y
```
### review labels
```
*** Review Labels
*** 

Install Verified label         [y/N]? y
```
### send mail setting
If your server could send out email, just leave these question as theirs default anwser.
```
*** Email Delivery
*** 

SMTP server hostname           [localhost]: 
SMTP server port               [(default)]: 
SMTP encryption                [none/?]: 
SMTP username                  : 
```
### container section
if you run gerrit as a standalone user, please create this user first. If you run gerrit in tomcat container, fill it with tomcat
```
*** Container Process
*** 

Run as                         [dylan]: gerrit
Java runtime                   [/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.144-7.b01.fc27.x86_64/jre]: 
Copy gerrit-2.14.4.war to gerrit-site/bin/gerrit.war [Y/n]? 
Copying gerrit-2.14.4.war to gerrit-site/bin/gerrit.war
```
### sshd section
just leave it as default.
```
*** SSH Daemon
*** 

Listen on address              [*]: 
Listen on port                 [29418]: 
Generating SSH host key ... rsa... dsa... ed25519... ecdsa 256... ecdsa 384... ecdsa 521... done
```
### httpd section
```
*** HTTP Daemon
*** 

Behind reverse proxy           [y/N]? 
Use SSL (https://)             [y/N]? 
Listen on address              [*]: 
Listen on port                 [8080]: 
Canonical URL                  [http://176.74.176.187:8080/]: 
```
### cache
```
*** Cache
*** 
```
### plugins
```
*** Plugins
*** 

Installing plugins.
Install plugin commit-message-length-validator version v2.14.4 [y/N]? y
Installed commit-message-length-validator v2.14.4
Install plugin download-commands version v2.14.4 [y/N]? y
Installed download-commands v2.14.4
Install plugin hooks version v2.14.4 [y/N]? y
Installed hooks v2.14.4
Install plugin replication version v2.14.4 [y/N]? y
Installed replication v2.14.4
Install plugin reviewnotes version v2.14.4 [y/N]? y
Installed reviewnotes v2.14.4
Install plugin singleusergroup version v2.14.4 [y/N]? y
Installed singleusergroup v2.14.4
Initializing plugins.
```

## run gerrit in tomcat container
### init gerrit site
as last section [init gerrit site](#init-gerrit-site)
### deploy gerrit.war
generally, you can just copy gerrit.war to path */var/lib/tomcat/webapps/*, tomcat will deploy it automaticlly.
after deployment, you should modify file content of *gerrit-launcher/workspace-root.txt* under your gerrit app directory, let its content point to gerrit site path.
```
/home/dylan/gerrit-site
```
### configure gerrit app context
add below content to tomcat config file *server.xml*, element Host.
```xml
<Context path="/gerrit" docBase="/var/lib/tomcat/webapps/gerrit" reloadable="true">
    <Resource
        name="jdbc/ReviewDb"
        type="javax.sql.DataSource"
        driverClassName="org.mariadb.jdbc.Driver"
        username="root"
        password="password"
        maxWait="10000"
        maxIdle="30"
        maxActive="100"
        url="jdbc:mariadb://localhost:3306/ReviewDB?autoReconnect=true"
        auth="Container" />
</Context>
```
### restart tomcat
```sh
systemctl restart tomcat
```
## a gerrit config sample
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
	smtpUser = gerrit@tianyisc.com
	sslVerify = false
	from = Gerrit Review <gerrit@tianyisc.com>
[receive]
	enableSignedPush = true
[container]
	user = tomcat
	javaHome = /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.144-0.b01.el7_4.x86_64/jre
[automerge]
	botEmail = gerrit@tianyisc.com
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