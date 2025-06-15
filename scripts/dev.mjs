#!/usr/bin/env zx

// Thothix Development Script - Cross-platform with Zx
// Usage: zx scripts/dev.mjs [action]

import { $, argv, echo, cd, path, fs } from 'zx'

// Configure for Windows compatibility
if (process.platform === 'win32') {
  $.shell = 'cmd.exe'
  $.prefix = ''
}

$.verbose = true

const action = argv._[0] || 'help'

echo`🔧 Thothix Development Script - Action: ${action}`

// Helper functions
const format = async () => {
  echo`📝 Formatting Go code...`
  await cd('./backend')
  await $`gofmt -w .`
  echo`✅ Formatting completed`
}

const lint = async () => {
  echo`🔍 Running golangci-lint...`
  await cd('./backend')
  try {
    await $`golangci-lint run --timeout=3m`
    echo`✅ Linting passed`
  } catch (error) {
    echo`❌ Linting failed`
    process.exit(1)
  }
}

const test = async () => {
  echo`🧪 Running tests...`
  await cd('./backend')
  try {
    await $`go test ./...`
    echo`✅ Tests passed`
  } catch (error) {
    echo`❌ Tests failed`
    process.exit(1)
  }
}

const preCommit = async () => {
  echo`🚀 Running pre-commit checks...`
  await format()

  echo`📋 Adding formatted files to git...`
  await cd('../')  // Return to project root
  await $`git add backend/`

  await lint()
  await test()
  echo`✅ Pre-commit checks completed`
}

const showHelp = () => {
  echo`
Usage: zx scripts/dev.mjs [action]

Actions:
  format      - Format Go code with gofmt
  lint        - Run golangci-lint
  test        - Run Go tests
  pre-commit  - Run format + git add + lint + test
  all         - Same as pre-commit

Examples:
  zx scripts/dev.mjs format
  zx scripts/dev.mjs pre-commit
  npm run format          # Via NPM script
  npm run pre-commit      # Via NPM script

Requirements:
  - Go installed
  - golangci-lint installed
  - Git repository
`
}

// Main execution
switch (action) {
  case 'format':
    await format()
    break
  case 'lint':
    await lint()
    break
  case 'test':
    await test()
    break
  case 'pre-commit':
  case 'all':
    await preCommit()
    break
  case 'help':
  default:
    showHelp()
    process.exit(action === 'help' ? 0 : 1)
}

echo`🎉 Script completed successfully!`
