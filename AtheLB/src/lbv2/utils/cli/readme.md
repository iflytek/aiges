## 节点获取 ##

- 节点获取
	- perf表示是否启用压测模式

```
cli.exe -all 1 -c 1 -n 10 -lbname lbv2 -nbest 10 -perf 0 -svc svc -subsvc sms -mode 0 -uid 12949

cli.exe -all 0 -c 1 -n 1 -lbname xsf-lbv2 -nbest 1 -perf 0 -svc svc -subsvc sms -mode 0
```


## 强制下线 ##

- 下线

```
cli.exe -mode 0 -action 1
```
