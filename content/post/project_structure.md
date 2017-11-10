---
title: "android source directory structure"
date: 2017-11-10T11:34:11+08:00
draft: true
tags: [android,platform]
topics: []
description: ""
---

- abi
- art
- bionic
- bootable --用于放置可引导程序源码，此处为recovery
- build --构建系统定义目录
- common --kernel
- cts
- dalvik --android中的Java虚拟机
- developers
- development
- device --各厂商定义自己设备和产品的目录
- docs
- external
- frameworks
- hardware --硬件相关，各种驱动
- libcore --系统核心库，比如dalvik的native的类加载器就实现在这里
- libnativehelper
- Makefile --入口Makefile，由build/core/root.mk复制过来
- ndk
- packages --通常为Android中app的源码放置处
- pdk
- prebuilts --预编译好的程序
- sdk
- system
