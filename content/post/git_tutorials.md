---
title: "git usage guide"
date: 2017-11-10T11:34:11+08:00
draft: true
tags: [git]
description: ""
---

# Git使用指南

## 基本配置
开始使用git之前一般需要配置两个项目，操作命令如下：
```sh
git config --global user.name "Dylan Deng"
git config --global user.email "dengxingxian@tianyisc.com"
```
分别把两行命令中引号里的内容改成你自己的名字和邮箱。

## 一些使用案例
```sh
git checkout -b neimengv1.0 neimengv1.0  # 取得tag`neimengv1.0`的内容并创建为本地分支neimengv1.0
git checkout -b master origin/master     # 取得master分支的内容并创建为本地分支master
git fetch origin refs/changes/34/34/1 && git checkout -b change-34 refs/changes/34/34/1 # 取得编号为34的改动的内容并创建为本地分支change-34
git push origin HEAD:refs/for/common     # 上传当前提交到代码审阅，审阅通过会合并到common分支
```

## 初始化仓库
### 获得远程仓库

```sh
git clone <repository_url> [-b <branch_name>] [directory]
```
repository_url为远程仓库地址，为必须参数。
-b 此选项为可选参数，如果不提供这个参数，下载完成之后默认检出master分支。
directory 为可选参数，如果不提供这个参数，git会根据url，在当前路径创建一个与仓库同名文件夹，然后将仓库克隆到此文件夹。
example:
```sh
git clone ssh://dylan@192.168.99.89:29418/aml_02 -b common aml_02
```

### 初始化本地仓库
如果远程仓库不存在或者远程仓库存在但是为空的，我们需要在本地创建一个git仓库然后上传到远程。操作过程如下：
```sh
cd $worktree                     # 进入到目标工作目录
git init                         # 初始化仓库
git add filename1 filename2      # 添加需要使用版本控制来追踪的文件
git commit -m "initial commit"   # 第一个commit
git remote add origin <repository_url>   # 和远程仓库关联
git push origin HEAD:refs/heads/mater    # 上传本地commit到远程
```

## 开始开发
### 查看修改

查看修改过的文件：
```sh
git status
```
查看具体修改：
```sh
git diff
```
查看指定文件的修改：
```sh
git diff -- filename1
```

### 添加文件

```sh
git add filename1 filename2 filename3
```
### 提交
```
git commit [-a | --interactive | --patch] [-s] [-v] [-u<mode>] [--amend]
	   [--dry-run] [(-c | -C | --fixup | --squash) <commit>]
	   [-F <file> | -m <msg>] [--reset-author] [--allow-empty]
	   [--allow-empty-message] [--no-verify] [-e] [--author=<author>]
	   [--date=<date>] [--cleanup=<mode>] [--[no-]status]
	   [-i | -o] [-S[<keyid>]] [--] [<file>…​]
```
只提交指定的文件
```sh
git commit -- filename1 filename2 filename3
```
refer [Git - git-commit Documentation](https://git-scm.com/docs/git-commit).

### 上传提交
```sh
git push <remote> [local_branch:]<remote_branch>
```
如果remote未指定，则默认为origin；如果local_branch未指定，则默认为HEAD，即当前最新提交；如果remote_branch未指定则默认为当前分支。
example:
```sh
git push origin HEAD:master
```
上传到代码审阅
```sh
git push origin HEAD:refs/for/<branch>     # <branch>替换成要上传的分支
```
更多用法参考[Git - git-push Documentation](https://git-scm.com/docs/git-push)

### 获取别人提交

获取修改：
```sh
git fetch [<options>] [<repository> [<refspec>…​]]
```
通常写作`git fetch`，即获取仓库origin的所有分支修改。

获取修改并更新工作空间：
```sh
git pull [<repository> [<refspec>…​]]
```
git pull实际等于git fetch + git rebase 或 git fetch + git merge，也就是说，更新工作空间实际由git merge或git rebase完成，具体取决于你的git 配置。相关配置项：pull.rebase，如果pull.rebase=true，将会采用rebase策略来更新工作空间，否则使用merge策略更新工作空间。推荐rebase策略，因为merge策略一般会生成一个类似_Merge branch 'common' of 172.16.8.220:aml02 into common_的提交，这个提交不包含有用信息，是多余的。

git pull必须在当前分支有upstream，并且工作空间干净时候才可以执行。如果当前不处于分支上，使用`git fetch origin; git rebase origin/<branch>`来更新工作空间。如果工作空间不干净，可以先执行`git stash save`来暂存更改，然后执行`git pull`，然后执行`git stash pop`来恢复之前的更改。

更多用法参考：
- [Git - git-fetch Documentation](https://git-scm.com/docs/git-fetch)
- [Git - git-pull Documentation](https://git-scm.com/docs/git-pull)

### 查看修改记录
查看修改记录的命令为git log，通常有以下几种用法：
```sh
git log --all --graph --decorate
```
- --all 查看所有分支修改（默认为当前分支）
- --graph 以图形形式展示分支历史
- --decorate 加上此选项，分支名会在log中显示

#### 其他常用参数
- -p 显示每个提交的具体修改
- -n n为数字，指定要查看的log条数
- --oneline 每个提交压缩为一行展示
- [-- path] 查看指定文件或路径的修改记录，eg: `git log -- filename1`
- --stat 统计每个提交修改过的文件

### 分支
创建一个新分支：
```sh
git branch <new_branch>
```
切换一个分支：
```sh
git checkout <branch_name>
```
从指定提交创建一个分支并切换到这个分支：
```sh
git checkout -b <branch_name> <refspec>
```
删除一个分支：
```sh
git branch -d <branch_name>
```
强制删除一个分支：
```sh
git branch -D <branch_name>
```
上传一个分支到远程仓库：
example:
```sh
git push origin branch_name:branch_name
```
删除一个远程仓库分支：
example:
```sh
git push origin :branch_name
```
### git-merge
git merge的效果一般像这样：
```
  A---B---C topic                        A---B---C topic
 /                      ------------>   /         \
D---E---F---G master                   D---E---F---G---H master
```
### git-rebase
rebase之前：
```
          A---B---C topic
         /
    D---E---F---G master
```
执行`git rebase master topic`之后
```

                  A'--B'--C' topic
                 /
    D---E---F---G master
```
### git-revert
git revert会通过生成一个新的提交的方式来撤销指定的之前的提交。
- [Git - git-revert Documentation](https://git-scm.com/docs/git-revert)
### git-reset
git reset可用于设置HEAD指针位置或者恢复文件，重写历史等。
- [Git - git-reset Documentation](https://git-scm.com/docs/git-reset)
### git-stash
应用场景：当前工作空间有更改，需要获取别人的修改并更新工作空间，当前工作空间修改不适宜提交。
暂存修改：
```sh
git stash save
```
查看stash列表：
```sh
git stash list
```
将暂存的修改恢复到工作空间并丢弃这个暂存：
```sh
git stash pop
```
相关配置项
rebase.autoStash

### 常用配置项
- core.editor string值为可以为vim,nano,gedit等，会影响commit时候编辑COMMIT_MSG的编辑器等
- core.autocrlf input,auto,false，commit时候对文本文件换行符的处理
- user.name string用户名
- user.email string用户邮箱
- pull.rebase bool执行git pull时候更新工作空间的策略
- rebase.autoStash bool执行rebase时候是否自动暂存未保存的更改
- http.proxy http代理
- user.signingKey 指定用于commit和tag签名的gpg私钥
- commit.gpgSign bool是否在commit自动签名

### git revision
某一个提交：
- HEAD 当前位置
- FETCH_HEAD 从服务器上获取到的指定的位置
- &lt;branch&gt; 分支
- <refspec> refs/changes/30/30/5或refs/heads/master这样的形式
- <commit id> 某一提交的commit id

某些时候需要输入*revision range*，格式为`<revison1>..<revision2>`，例：`origin/common..HEAD`
使用案例：
1. `git diff origin/master..master` 比较服务器版本master分支和本地master分支差异
2. `git log 9697889b0b58029debdd8a8310f1ce863bfd7dcd..f6502daaf3008a7cfe1d8a52cb698e1557de812d` 查看指定两个提交之间的修改日志

## 更多
Please refer [git-doc](https://git-scm.com/docs).
