# OpenAI compatibility

> **Note:** OpenAI compatibility is experimental and is subject to major adjustments including breaking changes. For fully-featured access to the AiAdmin API, see the AiAdmin [Python library](https://github.com/AllSage/AiAdmin-python), [JavaScript library](https://github.com/AllSage/AiAdmin-js) and [REST API](https://github.com/AllSage/AiAdmin/blob/main/docs/api.md).

AiAdmin provides experimental compatibility with parts of the [OpenAI API](https://platform.openai.com/docs/api-reference) to help connect existing applications to AiAdmin.

## Usage

### OpenAI Python library

```python
from openai import OpenAI

client = OpenAI(
    base_url='http://localhost:11434/v1/',

    # required but ignored
    api_key='AiAdmin',
)

chat_completion = client.chat.completions.create(
    messages=[
        {
            'role': 'user',
            'content': 'Say this is a test',
        }
    ],
    model='llama3',
)
```

### OpenAI JavaScript library

```javascript
import OpenAI from 'openai'

const openai = new OpenAI({
  baseURL: 'http://localhost:11434/v1/',

  // required but ignored
  apiKey: 'AiAdmin',
})

const chatCompletion = await openai.chat.completions.create({
  messages: [{ role: 'user', content: 'Say this is a test' }],
  model: 'llama3',
})
```

### `curl`

```
curl http://localhost:11434/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
        "model": "llama3",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "Hello!"
            }
        ]
    }'

```

## Endpoints

### `/v1/chat/completions`

#### Supported features

- [x] Chat completions
- [x] Streaming
- [x] JSON mode
- [x] Reproducible outputs
- [ ] Vision
- [ ] Function calling
- [ ] Logprobs

#### Supported request fields

- [x] `model`
- [x] `messages`
  - [x] Text `content`
  - [ ] Array of `content` parts
- [x] `frequency_penalty`
- [x] `presence_penalty`
- [x] `response_format`
- [x] `seed`
- [x] `stop`
- [x] `stream`
- [x] `temperature`
- [x] `top_p`
- [x] `max_tokens`
- [ ] `logit_bias`
- [ ] `tools`
- [ ] `tool_choice`
- [ ] `user`
- [ ] `n`

## Models

Before using a model, pull it locally `AiAdmin pull`:

```shell
AiAdmin pull llama3
```

### Default model names

For tooling that relies on default OpenAI model names such as `gpt-3.5-turbo`, use `AiAdmin cp` to copy an existing model name to a temporary name:

```
AiAdmin cp llama3 gpt-3.5-turbo
```

Afterwards, this new model name can be specified the `model` field:

```shell
curl http://localhost:11434/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
        "model": "gpt-3.5-turbo",
        "messages": [
            {
                "role": "user",
                "content": "Hello!"
            }
        ]
    }'
```
