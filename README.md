

## 简介

[祁启云免费验证系统](http://demo.developer.qqcloudcom.cn/) 一个使用golang语言、Web框架beego、前端Naive-Ui-Admin开发的免费网络验证系统

## 版本
- 当前版本1.11

## 更新方法
- 请直接将本目录中的verification.exe/verification直接覆盖到你服务器部署的目录，更新前，请先关闭正在运行的验证服务
- 新版本优化了批量创建激活码方法，所以在创建激活码之前，请先移除之前创建的所有激活码，才能正常添加激活码

## 更新内容
- 修复版本号管理 设置禁用不生效
- 修复软件配置更换之后旧配置依旧生效的bug
- 修复代理删除不清空卡密/用户的问题
- 修复批量删除显示删除成功，但是实际零删除的问题
- 修复删除软件，不清空归属订单/卡密/用户的问题
- 增加激活码/会员管理显示代理归属的功能
- 增加代理提卡之前先查询消耗金额的功能
- 增加全局刷新后台缓存的功能
- 增加订单列表，可查看代理提卡记录

## 特性
- 可跨平台部署，支持windows/linux双系统部署
- 部署简单，无需安装php环境、nginx、mysql等繁琐步骤，一步到位
- 自适应后台界面，操作简单，自带api文档，对接方便
- 支持aes/rsa加签，避免返回包被篡改

## 默认账号密码
- admin
- 112233
- 部署成功之后，访问地址为：[http://服务器ip:9960](http://ip:9960)，请注意默认端口为9960，记得修改默认账号密码


## API及使用说明文档
- [文档地址](https://console-docs.apipost.cn/preview/8848a5303c2cfbb3/e0393e2595dc925d)
- https://console-docs.apipost.cn/preview/8848a5303c2cfbb3/e0393e2595dc925d

## 免责声明
- 本程序开发初衷仅供学习交流使用，请勿将程序应用于任何非法场景，请在下载24小时后自行删除，因使用者所造成的所有法律后果，程序开发者概不负责！


## demo测试站
- [demo测试站](http://demo.developer.qqcloudcom.cn/)
- http://demo.developer.qqcloudcom.cn

## 更新历史
- 1.11 修复版本号管理 设置禁用不生效，修复软件配置更换之后旧配置依旧生效的bug，修复代理删除不清空卡密/用户的问题，增加订单列表...
- 1.1 修复无法修改账号资料、列表筛选bug、增加代理功能、增加角色权限
- 1.0 初发布的版本
