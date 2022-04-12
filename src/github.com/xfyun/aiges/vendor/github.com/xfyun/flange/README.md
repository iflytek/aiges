# Flange Golang

This repo is the Golang verison of Flange instrumentation. In this readme file, we will describle how to integrate it within your golang-based application/component.

 
## Flange Golang Specification

### How to use ?

1. download or checkout this repo, extract it into local disk;
2. add the root directory into your local `GOPATH`
3. code and run

### Normal RPC Integration
 
* In every service/component to be traced, you need initialize `Flange` before using  `Span`, normally at the entry, which just like the LOG component before;

```go

// call the function at the entry to init the global 'fc'
func initTracer() {
  // set local dump for example
  flange.DumpEnable = false
  flange.DeliverEnable = true
  // set the Spill Dir
  flange.SpillDir = "./log/span"
  flange.SpillEnable = true
  // set the flume ip:port, and backend-consumer number inside
  flange.Init("172.16.51.29", "4545", 4, "0.1.0.1", "8088", "sis")
}

 ```
 
* at the end of service, you can finish this component  by calling `flange.Fini()`, if needed

* While coding, the only thing you need to notice is `Span`. Such as create a root span by calling `flange.NewSpan()`, or create a child span by calling `Span.Next()`. **NOTICE: you must call Span.WithName(string) to set the span name.**;

1. create a root span by calling `flange.NewSpan()`;

 ``` go
 // create a root span
 // with deploy server info [ip:port:serverName]
 // and  set the span type, and sample switch 
sisSpan := flange.NewSpan(flange.SERVER, false).WithName("get-uid").Start()

 // your code 

 // finish the trace and report
sisSpan.End()
flange.Flush(sisSpan)
 ```

 2. create a child span by calling `Span.Next()`;

 ``` go
 // create a child span of current span
sis2iatSpan := sisSpan.Next(flange.CLIENT).WithName("rec-start").Start()

 // your code

 // finish the trace and report
sis2iatSpan.End()
flange.Flush(sis2iatSpan)
 ```

 * While in a normal RPC, **the span meta info must be passed by your RPC protocol**;

1. at the upstream service/component, get the span meta info, and pass it within you RPC protocol;

``` go
 // get the span meta info
 var meta = sis2iatSpan.Meta();

 // PASS THE META IN YOUR RPC
```

 2. at the downstream service/component, reconstruct the span status by `flange.FromMeta()`;
 
 ``` go
 // reconstruct the server-side spanlet
 // with meta info passed by RPC
 // and deploy server info [ip:port:serverName]
 // and set the span type
sis2iatSpan := flange.FromMeta(meta, flange.SERVER).WithName("rec-start").Start()

 // your code

 // finish the trace and report
sis2iatSpan.End()
flange.Flush(sis2iatSpan)
 ```

### RPC over MQ

* While in a MQ process, the usage of `Start/End` of normal RPC, should be changed to `Send` or `Recv`;

1. at the Producer-side，call `Span.Send()` to record the message send time;

```go
// your code

// record the producer span and report
mqSpan := sisSpan.Next(flange.PRODUCER).WithName("message-send").Send()
// report the span 
flange.Flush(mqSpan)
```

2. at the Consumer-side，call `Span.Recv()` to record the message consume time. **NOTICE: the consumer-side span is also reconstructed by span's meta info, which means that the span's meta info should be coverd in the message.**

```go
// your code

// record the consumer span and report
span := flange.FromMeta(meta, flange.CONSUMER).WithName("get-iat").Recv()
// report the span
flange.Flush(span)
```

### Other auxiliary info

 * There are others auxiliary info supported by Span, such as LocalComponent，ServerAddress as well as custom Tag. **NOTICE: the Span.Tag() and Span.tagLocalComponent()/Span.tagServerAddr() should appear in pairs.**


 ```go
// create a span
sis := sisSpan.Next(flange.CLIENT).WithName("get-uid").Start()

 // add a custom tag info
sis = sis.WithTag("sid", "iat@0000")

// add a local component tag info
sis = sis.WithLocalComponent().WithTag("localComponent", "zkr")

// add a server address tag info
sis = sis.WithServerAddr().WithTag("serverAddress", "192.168.1.1")

// add a ret tag info
sis = sis.WithRetTag("10030")

// add a error tag info
sis = sis.WithErrorTag("no server response")
 ```
 
 * Also, there are sdk-level span-gather-mechanism, the usage are :
 
 ```go
// create a span which is sampled by setting parameter abandon=true, so it will be abandon by the sample-mechanism
sampleSpan := flange.NewSpan(flange.SERVER, true).WithName("get-uid").Start()
```

moreover, you can enable the `ForceDeliver` to gather the sample-span forcible, which is the default setting.
```go
// enable the ForceDeliver to gather sample-span
flange.ForceDeliver = true
```

## Flange Golang RPC Example

Supposed the `sis` may  call `zkr`，demonstrated as below：

```
     ┌─────┐     ┌─────┐
---->│ sis │---->│ zkr │ 
     └─────┘     └─────┘
```

* in component of `sis`, function `get_iat` may be called by upstream component/service, therefore you can trace `sis` as below;


```go
// code in sis

// init the singleton 'fc' in sis entry
func initTracer() {
  // dump local for example
  flange.DumpEnable = false
  flange.DeliverEnable = true
  flange.Init("172.16.51.29", "4545", 4, "0.1.0.1", "8088", "sis")
}

// rpc interface called by upstream
func get_iat(metaStr string, ...) {
  // create a Span from meta
  span := flange.FromMeta(metaStr, flange.SERVER).WithName("get-iat").Start()

  // call downstream rpc interface with meta
  // when start a rpc , you need create a new Span, while client-side Span use the default SpanKind.CLIENT
  childSpan := span.Next(flange.CLIENT).WithName("get-iat-server").Start()
  // and pass this span's meta
  zrk.get_iat_server(childSpan.Meta(), ...)
  // and stop this rpc span
  childSpan.End()
  flange.Flush(childSpan)

  // your code

  span.End()
  flange.Flush(span)
}
```

* in component of `zkr`, function `get_iat_server` may be called by `sis`, therefore you can trace `zkr` as below;

```go
// code in zkr

// init the singleton tracer in sis entry
func initTracer() {
  // dump local for example
  flange.DumpEnable = true
  flange.DeliverEnable = false
  flange.Init("172.16.51.30", "4545", 4, "0.1.0.1", "8088", "zkr")
}

// rpc interface called by upstream
func get_iat_server(metaStr string, ...) {
  // create a Span from meta
  childSpan := flange.FromMeta(metaStr, flange.SERVER).WithName("get-iat-server").Start()

  // your code

  childSpan.End();
  flange.Flush(childSpan)
}
```

NOTICE:  the `childSpan` in component of `sis` and the `childSpan` in component of `zkr`, will be treaded as a `logically complete span` Span at the server-side of APM.


## Performance

with CPU:Intel(R) Xeon(R) CPU E5-2650 v4 @ 2.20GHz and Memory:128G, each test run 1000000 messages by 200 producer  and 4 consumer, we get following result (ops/sec):

|   ops/sec |   cpu(%) |   mem(%) |
|   ------  |   ------  |   :------:   |
|   119273    | 2314    |   0.1 |
|   121610    | 2420    |   0.1 |
|   119736    | 2399   |   0.1 |
|   121687    | 311.3   |   0.1 |
|   119476    | 2402   |   0.1 |
|   118484    | 1858   |   0.1 |
|   119566    | 2416   |   0.1 |
|   117337    | 2395   |   0.1 |
|   117849    | 100   |   0.0 |
|   121333    | 2046   |   0.1 |