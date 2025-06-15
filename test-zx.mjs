#!/usr/bin/env zx

import { $, echo } from 'zx'

// Configure for Windows
$.shell = 'cmd.exe'
$.prefix = ''

echo`🧪 Testing Zx functionality...`
await $`echo Hello from Zx!`
echo`✅ Test completed`
