# How to use this module

Simply get the module using `go get`

```shell
go get github.com/amritsingh183/golangchainagents
```

## Langchain agents using localollama models

- [Create local model using ollama](examples/models/Readme.md)
- Checkout [ToolBox](pkg/toolbox/toolbox.go) and examples usage at mock of [Taxi Booking tool](examples/langchain.go) and [weather checking tool](examples/langchain.go) which allows to have a similar experience that you get while using tools in [crewai](https://docs.crewai.com/concepts/tools#using-crewai-tools).
- The [examples/models](docs/models) help you get tool-enabled LLM (credits: [ollama](https://ollama.com/))
- Langchain itslef is very powerful and easy to use, but this repo helps you to build tools and integrate them with lanchain with just few lines of code, see [ToolBox](pkg/toolbox/toolbox.go). All you need is just implement `GetsWorkDone interface`. Check the taxi and weather tools mentioned above.
- This repository contains code which enables you to make function/tool calls using small large language models (llms).
This helps you to build and try AI agents locally.
The code uses `github.com/tmc/langchaingo/llms/openai` to make calls to local ollama models.
- Since most of us have less than 8GB of GPU, we can use this code easily and still have fun with agents on our own machines.
