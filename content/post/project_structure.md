---
title: "android source directory structure"
date: 2017-11-10T11:34:11+08:00
draft: false
tags: [android,platform]
description: ""
---

- abi --This folder contains a sub folder called cpp which actually contains many C++ files linked to many places.
- art --it is the folder that deals with the compilation of the latest android ART runtime. If you're looking into source directories of some other android versions, you won't fins it obviously
- bionic --the C-runtime for Android. Note that Android is not using glibc like most Linux distributions. Instead the c-library is called bionic and is based mostly on BSD-derived sources. In this folder you will find the source for the c-library, math and other core runtime libraries.
- bootable --boot and startup related code. Some of it is legacy, the fastboot protocol info could be interesting since it is implemented by boot loaders in a number of devices such as the Nexus ones.
- build --the build system implementation including all the core make file templates. An important file here is the envsetup.sh script that will help you a lot when working with the platform source. Running this script in a shell will enable commands to setup environment variables, build specific modules and grep in source code files.
- cts -- the compatability tests. The test suite to ensure that a build complies with the Android specification.
- dalvik --the source code for the implementation of the Dalvik Virtual Machine
- developers --
- development --projects related to development such as the source code for the sdk and ndk tools. Normally not a folder you touch when working with the platform for a target.
- device --product specific code for different devices. This is the place to find hardware modules for the different Nexus devices, build configurations and more.
- docs --I contains an important sub-folder called source.android.com. Contains tutorials, references, and miscellaneous information relating to the Android Open Source Project (AOSP). The current iteration of this site is fully static HTML (notably lacking in javascript and doxygen content), and is and/or was maintained by skyler (illustrious intern under Dan Morrill and assistant to the almighty JBQ).
- external --contains source code for all external open source projects such as SQLite, Freetype and webkit.
- frameworks --this folder is essential to Android since it contains the sources for the framework. Here you will find the implementation of key services such as the System Server with the Package- and Activity managers. A lot of the mapping between the java application APIs and the native libraries is also done here.
- hardware --hardware related source code such as the Android hardware abstraction layer specification and implementation. This folder also contains the reference radio interface layer (to communicate with the modem side) implementation.
- libcore --Apache Harmony.
- libnativehelper --Helper functions for use with JNI.
- Makefile --Makefile, build/core/root.mk
- ndk --native development kit
- out - the build output will be placed here after you run make. The folder structure is out/target/product/. In the default build for the emulator the output will be placed in out/target/product/generic. This is where you will find the images used by the emulator to start (or to be downloaded and flashed to a device if you are building for a hardware target).
- packages --contains the source code for the default applications such as contacts, calendar, browser.
- pdk -- I believe that 'pdk' is the Platform Development Kit, it's basically an SDK/set of tools that Google sends to OEMs to evaluate their framework ahead of each major Android upgrade since Android 4.1.
- prebuilts --contains files that are distributed in binary form for convenience. Examples include the cross compilations toolchains for different development machines.
- sdk --This directory contains lots of apps that are not part of operating system. There are quite useful apps that developers can leverage on and can be enhanced further as part of the operating system.
- system --source code files for the core Android system. That is the minimal Linux system that is started before the Dalvik VM and any java based services are enabled. This includes the source code for the init process and the default init.rc script that provide the dynamic configuration of the platform
- tools --Various IDE tools.
- vendor --This directory contains vendors specific libraries. Most of the proprietary binary libraries from non-open source projects are stored here when building AOSP.
