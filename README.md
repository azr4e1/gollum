# goLLuM

Small library for LLM completion. It provides a single client to interface to multiple LLM providers.

## LLM providers supported:

- OpenAI:
  - [X] Completion
  - [X] TTS
  - [ ] Embeddings
- Ollama:
  - [x] Completion
  - [ ] Embeddings
- Gemini:
  - [ ] Completion
  - [ ] Embeddings
- Claude
  - [ ] Completion
  - [ ] Embeddings


## Completion

```go
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "time"

  g "github.com/azr4e1/gollum"
)

func main() {
  // create new client
  cOllama, err := g.NewClient(g.WithProvider(g.OLLAMA), g.WithApiBase("http://localhost:11434"))
  if err != nil {
    panic(err)
  }
  cOllama.Timeout = 200 * time.Second // default is 30 for http request

  // create a chat
  chat := g.NewChat(o.SystemMessage("You are a Yakuza member. Act like it! Do not use emoji, and be very straightforward and to the point."), o.UserMessage("What is the difference between Proof of Work and Proof of Stake in Blockchain? What is your opinion on this? Which one is better?"))

  _, res, err := cOllama.Complete(o.WithModel("gemma2:2b"), o.WithMessages(chat.History()))
  if err != nil {
    panic(err)
  }

  b := new(bytes.Buffer)
  enc := json.NewEncoder(b)
  enc.SetIndent("", "  ")
  err = enc.Encode(res)
  if err != nil {
  	panic(err)
  }
  j, err := io.ReadAll(b)
  if err != nil {
  	panic(err)
  }
}
```

Response:

```json
{
  "Id": "",
  "Object": "",
  "Created": 38,
  "Model": "gemma2:2b",
  "Choices": [
    {
      "Index": 0,
      "Message": {
        "Role": "assistant",
        "Content": "Proof of work (PoW) requires massive computing power to solve complex equations.  This ensures the network's integrity, as only legitimate participants can contribute.\n\nProof of stake (PoS) uses a system where validators lock up their cryptocurrency as collateral. The more you have, the higher your chance of being selected to create and verify transactions. \n\nMy opinion? PoW is too cumbersome, energy-hungry. It's outdated.  PoS offers efficiency and less environmental impact. This leads to faster network speeds and lower costs. \n\nBetter? PoS is a cleaner solution, more sustainable in the long term. \n"
      },
      "FinishReason": "done"
    }
  ],
  "Usage": {
    "PromptTokens": 66,
    "CompletionTokens": 131,
    "TotalTokens": 197,
    "CompletionTokensDetails": null
  },
  "Error": {
    "Message": "",
    "Type": ""
  },
  "StatusCode": 200
}
```

Streaming is also supported. You need to provide a streaming function that will operate on the stream

```go
cOllama.EnableStream(func(req g.CompletionResponse) error {
  messages := req.Messages()
  if len(messages) == 1 {
    fmt.Print(messages[0])
  }
  return nil
})
_, _, err := cOllama.Complete(o.WithModel("gemma2:2b"), o.WithMessages(chat.History()))
if err != nil {
	panic(err)
}
```

Result:

```md
You want to know about Proof of Work and Proof of Stake, huh?

* **Proof of Work:**  Think of it like a competition. Miner have to solve complex math problems before they can add a new block to the chain. Takes tons of energy and resources.  Old school.
* **Proof of Stake:**  Instead of racing, you just need to own some coins.  The more you own, the bigger chance you get to propose a new block.  More sustainable.

Personally? I think PoS is cleaner, faster, cheaper. It's the future. Less headache for everyone.  But every coin gotta have its own method.
```

the text will appear as a stream on your terminal.

## Text To Speech

Currently only openai is supported

```go
package main

import (
  "os"
  "time"

  g "github.com/azr4e1/gollum"
)

func main() {
  apiKey := os.Getenv("OPENAI_API_KEY")
  client, err := g.NewClient(g.WithAPIKey(apiKey), g.WithProvider(g.OPENAI))
  if err != nil {
    panic(err)
  }
  client.Timeout = 200 * time.Second
  _, res, err := client.TextToSpeech(g.WithTTSInput("My name is Ken Takakura, but my friends call me Okarun."), g.WithTTSVoice("onyx"), g.WithTTSModel("tts-1-hd"))
  if err != nil {
    panic(err)
  }

  audio := res.Audio
  audioFile, err := os.Create("speech.mp3")
  defer audioFile.Close()
  if err != nil {
    panic(err)
  }
  _, err = audioFile.Write(audio)
  if err != nil {
    panic(err)
  }
}
```

