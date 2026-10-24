---
title: External Integration
nav_order: 3
has_toc: false
---

# External Integration

Use these interfaces to control DMXBOX from external software.
Configure group IDs and output destinations, then test connectivity with console output first.

| Interface | When to choose it | Details |
| --- | --- | --- |
| HTTP | You also need configuration, state, and operation logs | [HTTP API](http-api.md) |
| TCP | You need short, newline-delimited control commands | [TCP Protocol](tcp-protocol.md) |

TCP requires enabling the input module and restarting DMXBOX.
If external control is restricted, check [Browser only](operations-guide.md).
