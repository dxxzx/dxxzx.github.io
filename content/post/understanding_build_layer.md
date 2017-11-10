---
title: "understanding android build layer"
date: 2017-11-10T11:34:11+08:00
draft: false
tags: [build,android]
description: ""
---

# 理解构建层
先看看Android官方的解释

> ## Understand Build Layers
> 
> The build hierarchy includes the abstraction layers that correspond to the physical makeup of a device. These layers are described in the table below. Each layer relates to the one above it in a one-to-many relationship. For example, an architecture can have more than one board and each board can have more than one product. You may define an element in a given layer as a specialization of an element in the same layer, thus eliminating copying and simplifying maintenance.
> 
> Layer | Example | Description
> ----- | ------- | -----------
> Product | myProduct, myProduct_eu, myProduct_eu_fr, j2, sdk | The product layer defines the feature specification of a shipping product such as the modules to build, locales supported, and the configuration for various locales. In other words, this is the name of the overall product. Product-specific variables are defined in product definition Makefiles. A product can inherit from other product definitions, which simplifies maintenance. A common method is to create a base product that contains features that apply for all products, then creating product variants based on that base product. For example, you can have two products that differ only by their radios (CDMA vs GSM) inherit from the same base product that does not define a radio.
> Board/Device | sardine, trout, goldfish | The device/board layer represents the physical layer of plastic on the device (i.e. the industrial design of the device). For example, North American devices probably include QWERTY keyboards whereas devices sold in France probably include AZERTY keyboards. This layer also represents the bare schematics of a product. These include the peripherals on the board and their configuration. The names used are merely codes for different board/device configurations.
> Arch | arm, x86, mips, arm64, x86_64, mips64 | The architecture layer describes the processor configuration and ABI (Application Binary Interface) running on the board.

## 产品层
典型的产品定义文件夹一般会有以下文件

- vendorsetup.sh
- AndroidProducts.mk
- product.mk

### vendorsetup.sh
这个文件非常简单，只是添加产品条目而已，比如像下面这样：
```bash
add_lunch_combo p201_iptv-eng
add_lunch_combo p201_iptv-user
add_lunch_combo p201_iptv-userdebug
```

### AndroidProducts.mk
这个文件也很简单，只是指定产品的构建文件而已，例：
```Makefile
PRODUCT_MAKEFILES := $(LOCAL_DIR)/sample.mk
```
这里可以指定多个构建文件，一个文件对应一种产品，文件名必须与产品名相同。

### product.mk
product即为产品名，也就是AndroidProducts.mk中指定的构建文件。这个文件中一般定义了产品构建所必要的所有特性。

> Parameter | Description | Example
> --------- | ----------- | -------
> PRODUCT_AAPT_CONFIG | aapt configurations to use when creating packages | 
> PRODUCT_BRAND | The brand (e.g., carrier) the software is customized for, if any | 
> PRODUCT_CHARACTERISTICS | aapt characteristics to allow adding variant-specific resources to a package. | tablet,nosdcard
> PRODUCT_COPY_FILES | List of words like source_path:destination_path. The file at the source path should be copied to the destination path when building this product. The rules for the copy steps are defined in config/Makefile | 
> PRODUCT_DEVICE | Name of the industrial design. This is also the board name, and the build system uses it to locate the BoardConfig.mk. | tuna
> PRODUCT_LOCALES | A space-separated list of two-letter language code, two-letter country code pairs that describe several settings for the user, such as the UI language and time, date and currency formatting. The first locale listed in PRODUCT_LOCALES is used as the product's default locale. | en_GB de_DE es_ES fr_CA
> PRODUCT_MANUFACTURER | Name of the manufacturer | acme
> PRODUCT_MODEL | End-user-visible name for the end product | 
> PRODUCT_NAME | End-user-visible name for the overall product. Appears in the Settings > About screen. | 
> PRODUCT_OTA_PUBLIC_KEYS | List of Over the Air (OTA) public keys for the product | 
> PRODUCT_PACKAGES | Lists the APKs and modules to install. | Calendar Contacts
> PRODUCT_PACKAGE_OVERLAYS | Indicate whether to use default resources or add any product specific overlays | vendor/acme/overlay
> PRODUCT_PROPERTY_OVERRIDES | List of system property assignments in the format "key=value" | 

## 设备层
典型的设备层定义文件夹至少包含一个文件BoardConfig.mk，一般设备层还包含另外一个文件AndroidBoard.mk。这两个文件，BoardConfig.mk一般只是定义一些设备相关的宏，用于在其他地方使用，比如kernel,bootloader的编译。AndroidBoard.mk则是定义一些设备相关的编译规则和目标等。在很多例子里，设备定义路径与产品定义路径相同。比如我们项目中的device/amlogic/p201_iptv文件夹。

参考：

[Adding a New Device &nbsp;|&nbsp; Android Open Source Project](https://source.android.com/source/add-device)
