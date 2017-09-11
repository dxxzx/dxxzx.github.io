# configure window edge
## Related files
1. include/configs/gxl_p211_v1.h, line 145,
```c
#define CONFIG_EXTRA_ENV_SETTINGS \
        "firstboot=1\0"\
        "upgrade_step=0\0"\
        "jtag=apee\0"\
        "loadaddr=1080000\0"\
        "outputmode=720p50hz\0" \
        "hdmimode=720p50hz\0" \
        "cvbsmode=576cvbs\0" \
        "1080p_h=1037\0" \
        "1080p_w=1843\0" \
        "1080p_x=38\0" \
        "1080p_y=21\0" \
        "1080i_h=1037\0" \
        "1080i_w=1843\0" \
        "1080i_x=38\0" \
        "1080i_y=21\0" \
        "720p_h=691\0" \
        "720p_w=1222\0" \
        "720p_x=32\0" \
        "720p_y=14\0" \
        "576p_h=545\0" \
        "576p_w=680\0" \
        "576p_x=20\0" \
        "576p_y=15\0" \
        "576i_h=545\0" \
        "576i_w=680\0" \
        "576i_x=20\0" \
        "576i_y=15\0" \
        "480p_h=433\0" \
        "480p_w=649\0" \
        "480p_x=35\0" \
        "480p_y=23\0" \
        "480i_h=433\0" \
        "480i_w=649\0" \
        "480i_x=35\0" \
        "480i_y=23\0" \
        "uimode=1080p\0" \
        "display_width=1920\0" \
        "display_height=1080\0" \
        "display_bpp=16\0" \
        ...
```
As you see, variables related to window edge are defined in this file.
2. build/include/config.h, line 4,`#include <configs/gxl_p211_v1.h>`

This file include the file include/configs/gxl_p211_v1.h, make the macro CONFIG_EXTRA_ENV_SETTINGS visible for other files.
3. include/env_default.h, line 109,
```c
#ifdef  CONFIG_ENV_VARS_UBOOT_CONFIG                        
    "arch="     CONFIG_SYS_ARCH         "\0"                
    "cpu="      CONFIG_SYS_CPU          "\0"                
    "board="    CONFIG_SYS_BOARD        "\0"                
    "board_name="   CONFIG_SYS_BOARD        "\0"            
#ifdef CONFIG_SYS_VENDOR                                    
    "vendor="   CONFIG_SYS_VENDOR       "\0"                
#endif                                                      
#ifdef CONFIG_SYS_SOC                                       
    "soc="      CONFIG_SYS_SOC          "\0"                
#endif                                                      
#endif                                                      
#ifdef  CONFIG_EXTRA_ENV_SETTINGS                           
    CONFIG_EXTRA_ENV_SETTINGS                               
#endif                                                      
    "\0"
```

In this file, macro CONFIG_EXTRA_ENV_SETTINGS are included in macro CONFIG_ENV_VARS_UBOOT_CONFIG, this file are referenced in many files.

## Make uboot
we use predefined configuration file to make uboot
```sh
cd ${UBOOT_SOURCE_DIRECTORY}                    # enter source directory
make gxl_p211_v1_defconfig                      # prepare configuration file
make -j16                                       # start make
```

generated file are listed in direcotry fip
- fip/u-boot.bin
- fip/u-boot.bin.sd.bin   
- fip/u-boot.bin.usb.tpl
- fip/u-boot.bin.usb.bl2
- fip/gxl/u-boot.bin
- fip/gxl/u-boot.bin.sd.bin   
- fip/gxl/u-boot.bin.usb.tpl
- fip/gxl/u-boot.bin.usb.bl2

## Using generated uboot to make android system

To using generated uboot, just copy them to directory 'upgrade/1080p', or 'upgrade/720p' under device defination directory, like this:

```sh
cp -v fip/gxl/u-boot.bin{,.sd.bin,.usb.bl2,.usb.tpl} ../device/amlogic/p201_iptv/upgrade/1080p
```
the particular destination directory was decide by file factory.mk under device defination directory.
eg:
- device/amlogic/p211_iptv/factory.mk, line 17-34

```
ifeq ($(BUILD_CHINA_TELECOM_JICAI_APKS),true)
UPGRADE_FILES := \
    aml_sdc_burn.ini \
    1080p/u-boot.bin.sd.bin \
    1080p/u-boot.bin.usb.bl2 \
    1080p/u-boot.bin.usb.tpl\
    platform.conf \
    $(PACKAGE_CONFIG_FILE)
else
UPGRADE_FILES := \
    aml_sdc_burn.ini \
    720p/u-boot.bin.sd.bin \
    720p/u-boot.bin.usb.bl2 \
    720p/u-boot.bin.usb.tpl\
    platform.conf \
    $(PACKAGE_CONFIG_FILE)
endif
UPGRADE_FILES := $(addprefix $(TARGET_DEVICE_DIR)/upgrade/,$(UPGRADE_FILES))
```