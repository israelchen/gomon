# gomon
* Implementation of performance counters for Go (based on description [here](https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx)).
* Context-based Telemetry for Go with handlers.

Telemetry enables us to measure different aspects of contextual, logical operations by one or more handlers.

Handlers are trigged by the start and end of the logical operations and can trace, measure, aggregate and publish
the data as desired. One example can be the output of performance data using expvar, another example can be tracing, etc.

Please note that this implementation may not be idiomatic for Go and quite frankly completely stupid
but I still felt the urge to hook up telemetry for logical operations with decent performance counters and expvar.

Use at your own risk.
