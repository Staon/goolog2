# GoOLog2

_GoOLog2_ is a Go implementation of a logging framework I have made
and used for years in [Aveco](http://www.aveco.com/). The original
framework has been written in C++.

## Severity and verbosity

The main concept making the framework different to other ordinary logging
frameworks is distinction between severity and verbosity:

* __severity__ describes type of a logged message and its severity to
  logging process.
* __verbosity__ describes importance of a logged message.

Obviously, both concepts are not independent. One can hardly imagine
logging of a critical error at a low level of importance or
logging of a debugging message at a high level of importance.
I guess this dependency is the reason why most of the frameworks don't
distinct them.

However, the concepts are actually different. The best example is
the information severity (Info). This kind of messages covers entire
spectrum of verbosity levels: from highly important messages like server
starting and stopping to low important messages ment only for developers.

The framework distinguishes severity levels:

* __critical__ an error which cannot be recovered and which cause immediate
  ending of the process/service. The critical errors are usually last
  items in the log because they are logged just before the forced end
  caused by the error.
* __error__ severity marks errors of the process's state which can be
  recovered. An example can be loosing of connection to another service.
  The process can work somehow until the connection is re-established again.
  A common mistake is logging of errors returned to clients (if the logging
  process is a server) with the error severity. These conditions are not
  an errors of the server itself, hence they shouldn't be reported so.
* __warning__ severity reports some suspicious or potentially dangerous
  conditions. A using of a default for a missing configuration variable
  is an example.
* __info__ severity is the most usable one. All ordinary information about
  run of the process are reported at this severity (but at different
  verbosity levels).
* __debug__ severity is dedicated only for developers. Messages at this
  level shouldn't be seen in production environment.

The verbosity levels are defined by numbers. In theory, number of levels
is not limited. Practically, we have used 6 levels:

* __0__ logging disabled,
* __1,2__ messages important to operations,
* __3,4__ messages important to developers,
* __5__ too often periodical messages which would exhaust disk space
  so they are disabled at standard circumstances.

## Subsystem

Every message can be attached to a subsystem. The subsystem can be a library
or a module or special kind of log. Loggers attached to a subsystem accept
only messages of the subsystem. Loggers not attached to any subsystem
accept any message.

## Loggers

A _logger_ is and abstraction of a logging target. Currently there are two
loggers implemented: a _file logger_ and a _console logger_. The original
framework implemented two other loggers: logging into an Aveco's
proprietary logging system and logging into standard Unix _syslog_.
The last one is going to be implemented in the near future.

## Usage

Usage of the framework is easy. There is one global logging dispatcher.
At the beginning of the process, the dispatcher must be initialized:

```go
import(
  olog2 "github.com/Staon/goolog2"
)

olog2.Init("my_process")
```

The _Init_ function accepts a name called _system_. The name should be
unique allowing distinction of process's messages in a combined log
like _syslog_.

Next step of the initialization is creation of requested loggers:

```go
olog2.AddFileLogger("file", "", olog2.MaskAll, 4, "file.log", false)
olog2.AddConsoleLoggerStderr("console", "", olog2.MaskStd, 2)
```

The first argument of the __Add*__ functions is a logger name. The name
should be unique among all loggers. In the future the name is going to
be served as an address for dynamic changing of parameters of loggers.

The second argument is a subsystem. It can be empty if the logger
isn't attached to a subsystem.

The third and fourth arguments specify severities (as a mask) and maximal
verbosity of messages, which the logger accepts and logs.

Now the framework is ready to log:

```go
olog2.Info1("server started")
olog2.Critical1("No one has programmed this piece of SW!! I cannot work!")
olog2.Info1("server stopped") 
```
There is a set of convenient functions following simple name scheme:
_SeverityVerbosity_. The names can end with _s_ or _f_ - those
function allow specification of a subsystem or the message can be
formatted by standard _printf_ format description:

```go
olog2.Error2fs("special_module", "connection failed: %s", err)
```
Most of the time the convenient functions are good enough. However,
sometimes there is a need to specify severity and verbosity
dynamically. Then the functions _LogMessage()_ and _LogMessagef()_
can be useful.

At the end of the process the framework should be correctly cleaned. 
Hence the openede logging files can be correctly flushed and closed.

```go
olog2.Destroy()
```

__Warning:__ the initialization (Init and adding of loggers) phase and
destrucion phase are not thread safe! Be careful that all threads
have already stopped before you invoke the _Destroy()_ function.
