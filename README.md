[![CircleCI](https://circleci.com/gh/devatherock/slack-webhook-facade.svg?style=svg)](https://circleci.com/gh/devatherock/slack-webhook-facade)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/slack-webhook-facade.svg)](https://hub.docker.com/r/devatherock/slack-webhook-facade/)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/slack-webhook-facade.svg?sort=date)](https://hub.docker.com/r/devatherock/slack-webhook-facade/)
[![Docker Image Layers](https://img.shields.io/microbadger/layers/devatherock/slack-webhook-facade.svg)](https://microbadger.com/images/devatherock/slack-webhook-facade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# slack-webhook-facade
A Slack webhook facade to post messages to other chat clients like [Zulip](https://zulipchat.com/)

## Usage
### Zulip
To post a slack webhook message to zulip, use the slack webhook URL in the below format:

```
{slackWebhookFacadeHost}/zulip/{base64(username:zulipApiKey)}?server={zulipHost}
```

#### Sample URL
```
https://slack-webhook-facade.herokuapp.com/zulip/Y2ktYm90QHp1bGlwY2hhdC5jb206eHl6?server=https://devatherock-chat.zulipchat.com
```

#### Sample slack payload to post to the facade
```json
{
  "text": "https://circleci.com/gh/devatherock/git-sync/66 by devatherock",
  "channel": "general",
  "attachments": [
    {
      "title": "Build completed",
      "text": "https://circleci.com/gh/devatherock/git-sync/66 by devatherock",
      "color": "#764FA5"
    }
  ]
}
```

#### Parameters
**Path parameters**
- **slackWebhookFacadeHost** - Host name of your `slack-webhook-facade` instance
- **base64(username:zulipApiKey)** - Base64 encoded value of the Zulip bot integration's username and API key, joined
together by a colon. Will be used as the `Basic` authorization header in the call to Zulip API. If the username is
`ci-bot@zulipchat.com` and the API key is `xyz`, the path variable will be what is in the sample URL

**Query parameters**
- **zulipHost** - Host name of your Zulip instance

**Payload parameters**
- **channel** - The Zulip `stream` to post the message to
- **title** - The `topic` name in Zulip
