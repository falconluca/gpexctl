# gpexctl

## Get Started

1. 申请 YouTube API Key [传送门](https://console.developers.google.com/)
2. 将 API Key 配置到 `config.yaml` 文件的 `apiKey` 字段
3. 网站无法正常访问时, 可配置 HTTP 代理 (参见 `config.yaml` 文件的 `HTTPProxy` 字段)
4. 安装 gpexctl

```sh
make install
```

5. 查看帮助文档 `gpexctl youtube searh --help`

```bash
Usage:
  gpexctl youtube search [flags]

Flags:
  -h, --help               help for search
  -p, --period int         how long are videos from today. (default 10)
  -s, --size int           the size of videos. (default 10)
  -q, --term stringArray   search keywords for videos.

Global Flags:
      --debug   enable debug mode.
      --table   output as UI table format.
```

6. 搜索推荐最值得一看的内容

```bash
gpexctl youtube search -q 'GoodNotes' -q '生产力' --table
```

7. have fun :)

```bash
==>  🍾 共抓取视频 15 个
===============================
关键期 'GoodNotes,生产力' 最值得播放视频
===============================
编号    得分                    标题                                                    播放量   订阅量   🚪 传送门                                  
No.1    26751.666666666668     take notes with me #shorts #notetaking #stu...         16051   457     https://www.youtube.com/watch?v=AJ5ZsuE8ds8
No.2    2219.375               MUJI-inspired 2022 digital journal for #goo...         3551    574     https://www.youtube.com/watch?v=K8ycnJ99Gbo
No.3    340.2                  Add Stickers to Goodnotes in SECONDS!                  1701    1530    https://www.youtube.com/watch?v=opK8oXHrlPc
No.4    -                      2022 아이패드 다이어리 📔 종이느낌 가득한 심...                 115501  204000  https://www.youtube.com/watch?v=SbU37scy-uE
No.5    -                      FREE 2022 digital planner | complete, dated...         4310    10900   https://www.youtube.com/watch?v=ViRoYGQfEOA
No.6    -                      Fiz um planner digital #goodnotes #shorts              4933    114000  https://www.youtube.com/watch?v=CyZ1EzlgbHs
No.7    -                      2022 JOURNAL &amp; PLANNER GRATIS!| Goodnot...             2191    224000  https://www.youtube.com/watch?v=Un9TQQr8dCI
No.8    -                      GoodNotes 5 : Jak wgrać podręczniki? ⎮Nowe ...         1924    6760    https://www.youtube.com/watch?v=007W91mFeOY
No.9    -                      Creative Ways to Use GoodNotes for Digital ...         19914   97500   https://www.youtube.com/watch?v=EwsVw5Qvyr8
No.10   -                      Digital Journal with me on Goodnotes 🌼 Nov...         1400    34400   https://www.youtube.com/watch?v=Hs7ODlf9LvE
No.11   -                      全球性的 #严重货币超发、#生产力不平衡、#通缩...               34      1130    https://www.youtube.com/watch?v=WNm1D_UghDE
No.12   -                      要颜值还是要生产力？MateBook 14s全都给！                    7       1000000 https://www.youtube.com/watch?v=oVQjxjziu5g
No.13   -                      游戏本？生产力？Redmi G 2021酷睿版体验                      3       315     https://www.youtube.com/watch?v=GLXi0Ildtm0
No.14   -                      Vlog 创作者生产力修正                                     0       1000000 https://www.youtube.com/watch?v=Eh-JR87qfDE
No.15   -                      土豆学 | 《历史通论》，4.3，英国道路与生产力决定论             0       2       https://www.youtube.com/watch?v=81zycKeT4kA
```