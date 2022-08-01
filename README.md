# Go语言学习
[![Build Status](https://travis-ci.org/pibigstar/go-demo.svg?branch=master)](https://travis-ci.org/pibigstar/go-demo)
[![Code Coverage](https://codecov.io/gh/pibigstar/go-demo/branch/master/graph/badge.svg)](https://codecov.io/gh/pibigstar/go-demo/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pibigstar/go-demo)](https://goreportcard.com/report/github.com/pibigstar/go-demo)
[![License](https://img.shields.io/github/license/pibigstar/go-demo.svg?style=flat)](https://github.com/pibigstar/go-demo)
[![go-demo](https://img.shields.io/badge/go-demo-green)](https://github.com/pibigstar/go-demo)

- [base](base): Go语言基础
- [pprof](pprof): Go性能分析
- [design](design): Go实现常用设计模式
- [interview](interview): Go面试题及详解
- [sdk](sdk): Go对接第三方工具（mqtt、elastic、kafka...)
- [utils](utils): GoWeb开发常用工具类
- [blockchain](blockchain): Go实现一个简单的区块链
- [proxy](proxy): Go实现内网穿透工具
- [spider](spider): Go实现爬虫(QQ协议登录，QQ自动领礼物)
- [algo](algo): Go实现LeetCode中的算法题

## 项目结构
<details>
<summary>展开查看</summary>
<pre><code>
├─base
│  ├─context
│  ├─csv
│  ├─file
│  ├─flag
│  ├─goroutine
│  ├─http
│  │  ├─get_post
│  │  ├─restful
│  │  ├─server
│  │  └─url
│  ├─json
│  ├─mail
│  ├─mysql
│  ├─net
│  │  ├─client
│  │  └─server
│  ├─reflect
│  ├─regexp
│  ├─shell
│  ├─sort
│  ├─string
│  ├─sync
│  │  └─atomic
│  ├─time
│  ├─xml
│  └─zip
│      └─test
├─blockchain
│  ├─core
│  └─server
├─design
│  ├─adaptor
│  ├─chain
│  ├─decorator
│  ├─facade
│  ├─factory
│  │  ├─abstract
│  │  └─simple
│  ├─observer
│  ├─proxy
│  ├─singleton
│  ├─strategy
│  └─template
├─interview
│  ├─handpick
│  └─others
├─leetcode
│  ├─difficulty
│  ├─medium
│  └─simple
├─sdk
│  ├─alipay
│  ├─elasticsearch
│  ├─kafka
│  ├─mqtt
│  ├─nsq
│  │  ├─nsqio
│  │  └─test
│  ├─oss
│  ├─rabbitmq
│  ├─redis
│  ├─shortdomain
│  ├─sms
│  └─weixin
├─spider
│  ├─agent
│  ├─gift
│  │  ├─auto
│  │  └─hand
│  └─qq
│      ├─client
│      └─server
└─utils
    ├─cmp
    ├─cron
    ├─disk
    ├─encrypt
    ├─error
    ├─images
    ├─ip
    │  └─address
    ├─markdown
    ├─metadata
    ├─mock
    ├─multiconfig
    ├─name
    ├─pool
    ├─qrcode
    ├─rand
    ├─retry
    ├─seq
    ├─token
    ├─word
    └─xlsx
</pre></code>
</details>