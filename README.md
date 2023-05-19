# data2image
文件转图片

golang version >= 1.20
```shell
go install github.com/Rehtt/data2image@latest
```

## 文件转图片
```shell
data2image -i you/file/path
```
### 参数
- -i  输入路径（必须）
- -o  输出路径，默认：output
- -w  图片长度，默认：512
- -h  图片高度，默认：512
- -n  图片名称，默认：out%d.png


## 图片转文件
```shell
data2image -d -i you/image/path
```
### 参数
- -d  解码图片（必须）
- -i  输入路径（必须）
- -o  输出路径，默认：output
- -n  图片名称，默认：out%d.png