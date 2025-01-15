# create models

```shell
aria2c "https://huggingface.co/lmstudio-community/Llama-3-Groq-8B-Tool-Use-GGUF/resolve/main/Llama-3-Groq-8B-Tool-Use-Q4_K_M.gguf" -o Llama-3-Groq-8B-Tool-Use-Q4_K_M.gguf
aria2c "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-f16.gguf" -o Llama-3.2-3B-Instruct-f16.gguf
aria2c "https://huggingface.co/bartowski/SmolLM2-1.7B-Instruct-GGUF/resolve/main/SmolLM2-1.7B-Instruct-f16.gguf" -o SmolLM2-1.7B-Instruct-f16.gguf

ollama create SmolLM2-1.7B-Instruct-f16 -f SmolLM2-1.7B-Instruct-f16.model
ollama create Llama-3.2-3B-Instruct-f16 -f Llama-3.2-3B-Instruct-f16-model
ollama create Llama-3-Groq-8B-Tool-Use-Q4 -f Llama-3-Groq-8B-Tool-Use-Q4_K_M-model
```
