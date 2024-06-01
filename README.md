# Terraform provider for Open vSwitch

该Terraform提供程序管理本地Open vSwitch桥和端口。

## 使用示例

来自 [examples/sample-bridge](./examples/sample-bridge/):

```
terraform {
  required_providers {
    openvswitch = {
      source = "example.com/local/openvswitch"
      version = "1.0.0"
    }
  }
}

provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
  ip_address = "192.168.100.1/24"
}
```

## 重要说明

- ip、ovs-vsctl、ovs-ofctl 命令都需要 sudo 或 root 权限
- Error handling is currently broken

>  配置无密码的sudo，以允许Terraform用户无密码运行ovs-vsctl命令。在/etc/sudoers文件中添加以下行（使用visudo编辑器）：
>
> ```
> terraform_user ALL=(ALL) NOPASSWD: /usr/bin/ovs-vsctl
> ```
> 将terraform_user替换为运行Terraform的实际用户名。


## Installation from source

Requirements:

* Go 1.15.x or later
* GNU Make
* Terraform v0.12.* (This doesn't work on v0.13 yet)

克隆此仓库，然后执行以下操作：

```
$ go install
$ go build -o terraform-provider-openvswitch
```

## 本地安装Terraform Provider

Terraform Provider的全网标识符

Terraform的Provider在全网的的标识符由三部分组成，分别为hostname，namespace和type组成，即`<hostname>/<namespace>/<type>`。hostname是指分发、下载Provider的域名，默认为`registry.terraform.io`。namespace是指提供、开发Provider的组织的命名空间，默认为`hashicorp`。`type`是指Provider的具体类型。

例如有以下Terraform模板：

```
terraform {
    required_providers {
         alicloud = {
          source = "aliyun/alicloud"
          version = "1.126.0"
        }
    }
}
```

上述模板使用terraform init命令会默认去registry.terraform.io下载aliyun开发的alicloudProvider的1.126.0版本。

如果使用本地安装插件有两种方法。首先两种方法都需要将下载的Provider或者本地编译完成的Provider放置在以下文件目录层级：

```
XX(e.g. /usr/share/terraform/providers/)
└── <hostname>(e.g. registry.terraform.io)
    └── <namespace>(e.g. aliyun)
        └── <type>(e.g. alicloud)
            └── <version>(e.g. 1.127.0)
                └── <your OS>(e.g. linux_amd64)
                    └── <binary file>(e.g. terraform-provider-alicloud)
```

### 方法一：使用terraform init的自带参数

第一种方法，使用terraform init的plugin-dir参数：

```sh
terraform init -plugin-dir=/usr/share/terraform/providers
```

### 方法二：编写配置文件

第二种方法，编写`./terraformrc`配置文件，该文件需要放在`$HOME/`目录下：

```
provider_installation {
  filesystem_mirror {
    path    = "/usr/share/terraform/providers"
    include = ["registry.terraform.io/*/*"]
  }
}
```

其中include字段是指符合该通配符全网标识符的Provider，需要去/usr/share/terraform/providers查找本地Provider。./terraformrc的编写更详细的参数可以参考[官网](https://developer.hashicorp.com/terraform/cli/config/config-file?spm=a2c6h.12873639.article-detail.7.31045713cSFwkJ#provider-installation)
