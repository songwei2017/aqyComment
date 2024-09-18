学习
eg : https://www.xxyy180.com/mf/?url=https://www.iqiyi.com/v_1828bo3g16w.html 

1.启动服务
Q: 进人service/go 目录，go run main.go , 端口 1188，不用改，占用了需要改的话，tampermonkey/aqyComment.js 里面的端口也要改一下

2.谷歌油猴
Q: tampermonkey/aqyComment.js ，添加脚本.

3.打开链接，后面url为aiqiyi链接，对应更换
https://www.xxyy180.com/mf/?url=https://www.iqiyi.com/v_1828bo3g16w.html 

3.页面上出现的 aqyID 怎么获取
Q: 链接 https://iqiyi.com/v_12345.html , 则 id 为 v_12345

4.如果谷歌默认不允许http跨域
Q：chrome://flags/，搜索 Block insecure private network requests ，设置为 disable ，重启
