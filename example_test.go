// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zap_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// tweak NewDevelopmentConfig
func Example_wataash_zap_config() {
	outOrig := os.Stdout
	defer func() { os.Stdout = outOrig }()
	os.Stdout = os.Stderr

	// logger, _ := zap.NewDevelopment()
	// is:
	cfg := zap.Config{
		Level: zap.NewAtomicLevelAt(zap.DebugLevel),

		Development: true, // true: stack trace on warn
		// Development: false, // true: stack trace on error

		// DisableCaller: true,  // INFO	foo info
		// DisableCaller: false, // INFO	zap/example_test.go:117	foo info

		// DisableStacktrace: false,

		// same log only 2/sec ?
		// Sampling: &zap.SamplingConfig{
		// 	Initial:    2,
		// 	Thereafter: 100, // and re-logging after same 100 received?
		// },

		// Encoding: "console", // 14:56:12	INFO	zap/example_test.go:105	foo
		Encoding: "json",
		// {"L":"INFO","T":"14:56:19","C":"zap/example_test.go:105","M":"foo"}
		// sugar.warn:
		// {"L":"WARN","T":"14:56:19","C":"zap/example_test.go:106","M":"foo",
		//  "S":"go.uber.org/zap_test.ExampleWataashZap\n
		//         \t/Users/wsh/go/src/go.uber.org/zap/example_test.go:106\n
		//       testing.runExample\n
		//         \t/Users/wsh/go/src/go.googlesource.com/go/src/testing/example.go:121\n
		//       ..."
		// }

		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			MessageKey:    "M",
			StacktraceKey: "S",

			LineEnding: zapcore.DefaultLineEnding, // \n
			// LineEnding: "----------",

			EncodeLevel: zapcore.CapitalLevelEncoder,
			// zapcore.LevelEncoder
			// EncodeLevel: zapcore.LowercaseLevelEncoder,      // info
			// EncodeLevel: zapcore.LowercaseColorLevelEncoder, // info (blue)
			// EncodeLevel: zapcore.CapitalLevelEncoder,        // INFO
			// EncodeLevel: zapcore.CapitalColorLevelEncoder,   // INFO (blue)

			// EncodeTime: zapcore.ISO8601TimeEncoder,
			// zapcore.TimeEncoder
			// EncodeTime: zapcore.EpochTimeEncoder,       // 1.555307003305727e+09
			// EncodeTime: zapcore.EpochMillisTimeEncoder, // 1.555307012028605e+12
			// EncodeTime: zapcore.EpochNanosTimeEncoder,  // 1555307020546399000
			// EncodeTime: zapcore.ISO8601TimeEncoder,     // 2019-04-15T14:35:32.706+0900
			// EncodeTime: nil,                            // (none)
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("15:04:05"))
			},

			EncodeDuration: zapcore.StringDurationEncoder,
			// zapcore.DurationEncoder
			// EncodeDuration: zapcore.SecondsDurationEncoder, // ?
			// EncodeDuration: zapcore.NanosDurationEncoder,   // ?
			// EncodeDuration: zapcore.StringDurationEncoder,  // ?

			EncodeCaller: zapcore.ShortCallerEncoder,
		},

		OutputPaths: []string{"stderr"},
		// OutputPaths: []string{"stderr", "stderr", "stderr", "stdout", "/tmp/wataash/zap.debug"},

		ErrorOutputPaths: []string{"stderr"},

		// InitialFields: map[string]interface{}{"url": "http://www.empty.com"},
	}
	// TODO: yaml, json config
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Debugf("foo debug") // {"L":"DEBUG","T":"15:33:15","C":"zap/example_test.go:115","M":"foo debug"}
	sugar.Infof("foo info")   // {"L":"INFO","T":"15:33:15","C":"zap/example_test.go:116","M":"foo info"}
	sugar.Warnf("foo warn")   // {"L":"WARN","T":"15:34:23","C":"zap/example_test.go:117","M":"foo warn","S":"go.uber.org/zap_test.Example_..."}
	sugar.Errorf("foo error") // {"L":"ERROR","T":"15:34:31","C":"zap/example_test.go:118","M":"foo error","S":"go.uber.org/zap_test.Example_..."}

	// {"L":"INFO","T":"15:53:55","C":"zap/example_test.go:133","M":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
	sugar.Infow("Failed to fetch URL.",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)

	// // Output:
}

func Example_wataash_zap_product() {
	outOrig := os.Stdout
	defer func() { os.Stdout = outOrig }()
	os.Stdout = os.Stderr

	// logger, err := zap.NewProduction()
	// is:

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	// {"level":"info","ts":1555309053.257174,"caller":"zap/example_test.go:59","msg":"foo"}
	// {"level":"warn","ts":1555309053.257233,"caller":"zap/example_test.go:60","msg":"foo"}
	// {"level":"error","ts":1555309053.2572422,"caller":"zap/example_test.go:61","msg":"foo",
	//  "stacktrace":
	//  "go.uber.org/zap_test.Example_wataash_zap_product\n
	//    \t/Users/wsh/go/src/go.uber.org/zap/example_test.go:61\n
	//  testing.runExample\n
	//    \t/Users/wsh/go/src/go.googlesource.com/go/src/testing/example.go:121\n
	//  ..."
	// }
	sugar.Info("foo")
	sugar.Warn("foo")
	sugar.Error("foo")

	// // Output:
}

// for cli applications
// * stacktrace: error, fatal (by Development: false)
// * --log (default: ["stdout"]) colored, hh:MM:ss
// * --log-json (default: [])
// * "stdout" and "stderr"
// TODO: stackdriver sentry rollbar
//       https://github.com/blendle/zapdriver
func Example_wataash_zap_cli() {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapCfg.DisableCaller = false
	zapCfg.DisableStacktrace = false
	zapCfg.Sampling = nil
	zapCfg.Encoding = ""

	optLogs := []string{"stderr", "/dev/stderr"}
	zapCfg.Encoding = "console"
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCfg.EncoderConfig.EncodeTime =
		func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("15:04:05"))
		}
	zapCfg.OutputPaths = optLogs
	loggerConsole, err := zapCfg.Build()
	sugarConsole := loggerConsole.Sugar()
	_ = err

	optLogsJson := []string{"stderr"}
	zapCfg.Encoding = "json"
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapCfg.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	zapCfg.OutputPaths = optLogsJson
	loggerJson, err := zapCfg.Build()
	sugarJson := loggerJson.Sugar()
	_ = err

	sugarConsole.Warnw("foo")
	sugarJson.Warnw("foo")
	sugarConsole.Errorw("foo")
	sugarJson.Errorw("foo")

	// TODO: merge them with zapcore.NewTee

	// // Output:
}

func Example_presets() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"

	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.",
		// Structured context as strongly typed fields.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	// Output:
	// {"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
	// {"level":"info","msg":"Failed to fetch URL: http://example.com"}
	// {"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
}

func Example_basicConfiguration() {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	// same in yaml
	// spaces for indent are invalid...
	rawYaml := []byte(`
level: debug
encoding: json
outputPaths: [stdout, /tmp/logs]
errorOutputPaths: [stderr]
initialFields:
  foo: bar
encoderConfig:
  messageKey: message
  levelKey: level
  levelEncoder: lowercase
`)
	if err := yaml.Unmarshal(rawYaml, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	// Output:
	// {"level":"info","message":"logger construction succeeded","foo":"bar"}
}

func Example_advancedConfiguration() {
	outOrig := os.Stdout
	defer func() { os.Stdout = outOrig }()
	os.Stdout = os.Stderr

	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	defer logger.Sync()
	logger.Info("constructed a logger")
	// // Output:
}

func ExampleNamespace() {
	logger := zap.NewExample()
	defer logger.Sync()

	logger.With(
		zap.Namespace("metrics"),
		zap.Int("counter", 1),
	).Info("tracked some metrics")
	// Output:
	// {"level":"info","msg":"tracked some metrics","metrics":{"counter":1}}
}

func ExampleNewStdLog() {
	logger := zap.NewExample()
	defer logger.Sync()

	std := zap.NewStdLog(logger)
	std.Print("standard logger wrapper")
	// Output:
	// {"level":"info","msg":"standard logger wrapper"}
}

func ExampleRedirectStdLog() {
	logger := zap.NewExample()
	defer logger.Sync()

	undo := zap.RedirectStdLog(logger)
	defer undo()

	log.Print("redirected standard library")
	// Output:
	// {"level":"info","msg":"redirected standard library"}
}

func ExampleReplaceGlobals() {
	logger := zap.NewExample()
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("replaced zap's global loggers")
	// Output:
	// {"level":"info","msg":"replaced zap's global loggers"}
}

func ExampleAtomicLevel() {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("info logging enabled")

	atom.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","msg":"info logging enabled"}
}

func ExampleAtomicLevel_config() {
	// The zap.Config struct includes an AtomicLevel. To use it, keep a
	// reference to the Config.
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("info logging enabled")

	cfg.Level.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","message":"info logging enabled"}
}

func ExampleLogger_Check() {
	logger := zap.NewExample()
	defer logger.Sync()

	if ce := logger.Check(zap.DebugLevel, "debugging"); ce != nil {
		// If debug-level log output isn't enabled or if zap's sampling would have
		// dropped this log entry, we don't allocate the slice that holds these
		// fields.
		ce.Write(
			zap.String("foo", "bar"),
			zap.String("baz", "quux"),
		)
	}

	// Output:
	// {"level":"debug","msg":"debugging","foo":"bar","baz":"quux"}
}

func ExampleLogger_Named() {
	logger := zap.NewExample()
	defer logger.Sync()

	// By default, Loggers are unnamed.
	logger.Info("no name")

	// The first call to Named sets the Logger name.
	main := logger.Named("main")
	main.Info("main logger")

	// Additional calls to Named create a period-separated path.
	main.Named("subpackage").Info("sub-logger")
	// Output:
	// {"level":"info","msg":"no name"}
	// {"level":"info","logger":"main","msg":"main logger"}
	// {"level":"info","logger":"main.subpackage","msg":"sub-logger"}
}

func ExampleWrapCore_replace() {
	// Replacing a Logger's core can alter fundamental behaviors.
	// For example, it can convert a Logger to a no-op.
	nop := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	})

	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("working")
	logger.WithOptions(nop).Info("no-op")
	logger.Info("original logger still works")
	// Output:
	// {"level":"info","msg":"working"}
	// {"level":"info","msg":"original logger still works"}
}

func ExampleWrapCore_wrap() {
	// Wrapping a Logger's core can extend its functionality. As a trivial
	// example, it can double-write all logs.
	doubled := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, c)
	})

	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("single")
	logger.WithOptions(doubled).Info("doubled")
	// Output:
	// {"level":"info","msg":"single"}
	// {"level":"info","msg":"doubled"}
	// {"level":"info","msg":"doubled"}
}
