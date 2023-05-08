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

## 介绍

本工具可用于码表生成与维护，主要功能有词组自动编码（对于带有音码的方案可以自动注音），对已有词库进行编码校验。

编码规则非常强大，几乎支持所有方案：`一般4码定长方案` 形码（支持构词码与单字编码不同）、双拼定长（如小鹤音形）、两笔，也支持 `顶功`（如星空键道 6，星空两笔），支持 `不定长`方案，如红辣椒五笔，全拼等。

## 使用

直接拖动 配置文件.ini 到程序图标上。或者可以用命令行：

```sh
.\lilac.exe path.ini
```

编辑 `pinyin-data/correct.txt` 可以修正读音。

## 配置文件详解

配置文件分为以下几个部分：

- `Config`: 主要配置
- `Char`: 用于生成词组编码的数据
- `Mapping`: 拼音映射，主要用于音码方案
- `Dict`: 主要词库部分
- `Check`: 想要校验的词库

### Config

#### _Rule_

自动生成编码的规则，由 `,` 分隔不同词长，`<词长>:<编码规则>`，省略词长则作为默认的规则。

规则中每个编码由 `+` 分隔，编码由一个字母和一个数字表示，字母表示取第几个字（a 或 A 为第一个字，b 或 B 为第二个字，z 或 Z 为最后一个字），大写字母表示取音码(源自`Mapping`)，小写表示取形码(源自`Char`)，数字则为取第几码，省略数字则全取。

整句有特殊规则，以 `..(分隔符)` 结尾，如 `..'` 表示以 `'` 分隔，分隔符可以省略。

<details>
<summary>例子</summary>

```ini
; 形码
Rule = 2:a1+a2+b1+b2, \
    3:a1+b1+c1+c2, \
    :a1+b1+c1+z1

; 双拼 2码音 2码形
Rule = 2:A+B, \
    3:A1+A2+B1+C1, \
    :A1+B1+C1+Z1

; 键道6顶功 2码音 4码形
Rule = 2:A+B+a1+b1, \
    3:A1+B1+C1+a1+b1+c1, \
    :A1+B1+C1+Z1+a1+b1

; 两笔 1码音 3码形
Rule = 2:A+a1+B+b1, \
    3:A+a1+B+C, \
    :A+B+C+Z

; 星空两笔顶功 1码音 5码形
Rule = 2:A+a1+B+b1+a2+b2, \
    3:A+B+C+a2+b2+c2, \
    :A+B+C+Z+a2+b2

; 整句第一个字母随便填，只和大小写有关，反正都要取遍每一个字
; 红辣椒五笔
Rule = :n1+n2..

; 双拼整句？空格分隔（用_代替）
Rule = :N.._
```

</details>

#### _保留单字全码(true|false)_

生成简码后是否需要保留单字全码，保留后会追加到词库末尾。

#### _单字简码规则_

由 `,` 分隔不同编码长，`<编码长>:<最大重码数>`，省略重码数则为无限重码，若无对应的编码长，则重码数为 1。

#### _词组简码规则_

同上

#### _Sort(true|false)_

是否按照按编码重新排序。

### Char

构词码，每个字可以有多个编码，一般是形码  
对形码来说是单字全码  
对双拼，辅助码（可以生成单字编码）  
对两笔，后三码

> 支持 `>>()` 语法引用文件。

### Mapping

双拼或者两笔生成词组所需的拼音映射表

两种格式：

- `大写字母<tab>音素`: 大写字母为键位，音素一般是声母或者韵母，多个音素用空格分隔。
- `音节<tab>键位`: 都为小写字母

<details>
<summary>例子：大牛双拼</summary>

```ini
; 按键(大写)基础映射
Q	q ua ian
W	w ei vn
E	e
R	r ou
T	t iu
Y	y un
U	sh u
I	ch i
O	zh o uo
P	p ie
A	zh a
S	s ao
D	d an
F	f ang
G	g uai ing
H	h ai ue
J	j eng van
K	k en ia
L	l ong iong
Z	z uan
X	x ve uang
C	c ian
V	sh v ui
B	b in
N	n ui iang
M	m iao

; 自定义音节，零声母都需要自定义
a	ea
ai	eh
an	ed
ang	ef
ao	es
e	ee
ei	ew
en	ek
er	eu
o	eo
ou	er
; 一个音节可以映射多种按键组合
shi	ui vi
```

</details>

<details>
<summary>例子：哲豆音形</summary>

```ini
Q	q
W	w
; E
R	r
T	t
Y	y
U	q
I	z zh
; O
P	p
; A
S	s sh
D	d
F	f
G	g
H	h
J	j
K	k
L	l
Z	z zh
X	x
C	c ch
; V
B	b
N	n
M	m

; 自定义音节，零声母都需要自定义
a	a
ai	a
an	a
ang	a
ao	a
e	e
ei	e
en	e
er	e
o	o
ou	o
```

</details>

### Dict

分为几种格式：

- `字词<tab>编码`: 不做处理
- `字词`: 根据规则自动生成编码，若要用到拼音也会自动注音
- `?字词<tab>拼音`: 根据给定的拼音生成编码
- `>>(文件路径)`: 引入另一个文件路径
- `?>>(文件路径)`: 引入带拼音的词库辅助生成编码

### Check

校验已有词库，格式为多多 `字词<tab>编码`

> 支持 `>>()` 语法引用文件。

