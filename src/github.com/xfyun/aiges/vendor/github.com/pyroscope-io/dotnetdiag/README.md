# .NET diagnostics

The package provides means for .Net runtime diagnostics implemented in Golang:
 - [Diagnostics IPC Protocol](https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#transport) client.
 - [NetTrace](https://github.com/microsoft/perfview/blob/main/src/TraceEvent/EventPipe/EventPipeFormat.md) decoder.

### Diagnostic IPC Client

```
# go get github.com/pyroscope-io/dotnetdiag
```

Supported .NET versions:
 - .NET 5.0
 - .NET Core 3.1

Supported platforms:
 - [x] Windows
 - [x] Linux
 - [x] MacOS

Implemented commands:
 - [x] StopTracing
 - [x] CollectTracing
 - [ ] CollectTracing2
 - [ ] CreateCoreDump
 - [ ] AttachProfiler
 - [ ] ProcessInfo
 - [ ] ResumeRuntime

### NetTrace decoder

```
# go get github.com/pyroscope-io/dotnetdiag/nettrace
```

Supported format versions: <= 4

The decoder deserializes `NetTrace` binary stream to the object sequence. The package contains an example stream
handler implementation that processes events from **Microsoft-DotNETCore-SampleProfiler** provider. See [examples](examples) directory.
