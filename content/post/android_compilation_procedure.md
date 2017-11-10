---
title: "understanding android compiling progress"
date: 2017-09-12T09:21:11+08:00
draft: true
---

- build/core/main.mk:108:`include $(BUILD_SYSTEM)/config.mk`
- build/core/config.mk:137:`include $(BUILD_SYSTEM)/envsetup.mk`
- build/core/config.mk:155:`include $(board_config_mk)`
- build/core/envsetup.mk:115:`include $(BUILD_SYSTEM)/product_config.mk`
- build/core/product_config.mk:179:`include $(BUILD_SYSTEM)/product.mk`
- build/core/product_config.mk:180:`include $(BUILD_SYSTEM)/device.mk`
- build/core/config.mk:45:`SRC_TARGET_DIR := $(TOPDIR)build/target`
- build/core/config.mk:139:`# Boards may be defined under $(SRC_TARGET_DIR)/board/$(TARGET_DEVICE)`
- build/core/config.mk:145:`               $(SRC_TARGET_DIR)/board/$(TARGET_DEVICE)/BoardConfig.mk \`
- build/core/product_config.mk:185:`    $(SRC_TARGET_DIR)/product/AndroidProducts.mk)`
- build/core/product_config.mk:267:`$(foreach runtime, $(product_runtimes), $(eval include $(SRC_TARGET_DIR)/product/$(runtime).mk))`
- build/core/product.mk:33:`  $(SRC_TARGET_DIR)/product/AndroidProducts.mk`