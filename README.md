
<div align="center">

<img src="logo.png"  width="150" height="150"> </img>

### 丁香码表生成器

[![GitHub Repo stars](https://img.shields.io/github/stars/flowerime/lilac)](https://github.com/flowerime/lilac/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/flowerime/lilac)](https://github.com/flowerime/lilac/network/members)
![GitHub repo size](https://img.shields.io/github/repo-size/flowerime/lilac)
![GitHub](https://img.shields.io/github/license/flowerime/lilac)
<!-- [![GitHub release (latest by date)](https://img.shields.io/github/v/release/flowerime/lilac)](https://github.com/flowerime/lilac/releases) -->
<!-- [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/flowerime/lilac/build.yml)](https://github.com/flowerime/lilac/actions/workflows/build.yml) -->

</div>

## TODO

- [x] 从全码表生成出简不出全码表（需要规则）
- [x] 从简码全码混和码表中提取单字全码
- [x] 根据单字全码表对词组进行编码<del>（对多编码的字进行笛卡尔积，需要一种编码规则）</del>
- [x] 根据单字全码表对词库进行错码校验
- [x] 根据全拼词库生成双拼定长码表
- [ ] 词条过滤器（例：词长>9、码长>=5、词频<10）
- [ ] 从单字全码表生成词库

基于拼音的编码生成。。

适配带飞键的双拼方案。

笛卡尔积。。。

## 成词规则

### 形码

2:AaAbBaBb AABB
3:AaBaCaCb ABCC
:AaBaCaZa AACZ

形音 希码 组词用不到音

### 纯双拼

成词只用到全拼  
只需要考虑四码定长  
第三码的规则
AABC
ABCC

### 音形结合的 只用到声母

需要额外定义一个全拼到声母()的映射表

音形 两笔 

### 音形 键道6 星空两笔

### 红辣椒五笔

取每个字前两码
ab...
