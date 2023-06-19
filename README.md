# ethscan

仿[以太坊浏览器](https://etherscan.io/)后端，相关API用例在HTTPClient目录下。

本项目实现了对nonce, balance, block, transaction的查询，以及单个block下所有transactions的分页查询。`send_sign`和`send_raw`接口用于广播已签名的交易。具体用例可查看HTTPClient下的http文件，Goland可直接运行，VSCode安装插件后可运行用例。