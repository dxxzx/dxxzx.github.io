---
title: "jenkins ldap auth"
date: 2017-11-10T13:36:11+08:00
draft: false
tags: [jenkins,ldap]
description: ""
---

# Jenkins Authorize with LDAP

## Install LDAP plugin
login to jenkins, go to `Manage Jenkins -> Manage Plugins -> Available` to install `LDAP plugin`.

## Configure LDAP plugin
login to jenkins, go to `Manage Jenkins -> Configure Global Security -> Access Coutrol` to configure LDAP.

- Server: &lt;host name of ldap server&gt;, possible value: *192.168.1.101*, *ldap.example.com*
- root DN: usually be *dc=my-domain,dc=com*
- User search base: usually be *ou=People*
- User search filter: usually be *uid={0}*
- Manager DN: usually be *cn=Manager,dc=my-domain,dc=com* or *cn=admin,dc=my-domain,dc=com*
- Manager Password: Your Manager password
- Display Name LDAP attribute: usually be displayname
- Email Address LDAP attribute: usually be mail
