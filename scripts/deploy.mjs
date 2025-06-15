#!/usr/bin/env zx

// Thothix Deployment Script - Cross-platform with Zx
// Usage: zx scripts/deploy.mjs [env] [command] [options]

import { $, argv, echo, fs, path } from 'zx'

// Configure for Windows compatibility
if (process.platform === 'win32') {
  $.shell = 'cmd.exe'
  $.prefix = ''
}

$.verbose = true

const env = argv._[0] || 'help'
const cmd = argv._[1] || 'up'
const opt = argv._[2]

if (env === 'help') {
  echo`
Usage: zx scripts/deploy.mjs [env] [command] [options]

Environments:
  dev      - Development environment (.env)
  staging  - Staging environment (.env.staging)
  prod     - Production environment (.env.prod) with Vault

Commands:
  up       - Start services
  down     - Stop services
  logs     - Show logs [service_name]
  status   - Show container status
  vault    - Vault commands (init, ui, status)

Examples:
  zx scripts/deploy.mjs dev up
  zx scripts/deploy.mjs prod logs backend
  zx scripts/deploy.mjs dev vault ui
  npm run dev              # Via NPM script
  npm run staging          # Via NPM script
`
  process.exit(0)
}

// Environment configuration
const configs = {
  dev: {
    envFile: '.env',
    composeFiles: ['-f', 'docker-compose.yml']
  },
  staging: {
    envFile: '.env.staging',
    composeFiles: ['-f', 'docker-compose.yml', '-f', 'docker-compose.staging.yml']
  },
  prod: {
    envFile: '.env.prod',
    composeFiles: ['-f', 'docker-compose.yml', '-f', 'docker-compose.prod.yml']
  }
}

const config = configs[env]
if (!config) {
  echo`❌ Invalid environment: ${env}. Use dev, staging, or prod`
  process.exit(1)
}

// Check environment file exists
const envFilePath = config.envFile
if (!await fs.pathExists(envFilePath)) {
  echo`❌ Environment file ${config.envFile} not found!`
  echo`📋 Please copy .env.example to ${config.envFile} and configure it`
  process.exit(1)
}

echo`🚀 Thothix Deployment - Environment: ${env}, Command: ${cmd}`

// Helper function for docker compose
const dockerCompose = async (args) => {
  const composeArgs = [...config.composeFiles, `--env-file=${config.envFile}`, ...args]
  return $`docker compose ${composeArgs}`
}

// Main command execution
switch (cmd) {
  case 'up':
    echo`📦 Starting ${env} environment...`
    await dockerCompose(['up', '-d', '--build'])
    echo`✅ ${env} environment started successfully`
    echo`🔍 Container status:`
    await dockerCompose(['ps'])
    break

  case 'down':
    echo`🛑 Stopping ${env} environment...`
    await dockerCompose(['down'])
    echo`✅ ${env} environment stopped`
    break

  case 'logs':
    if (opt) {
      echo`📋 Showing logs for service: ${opt}`
      await dockerCompose(['logs', '-f', opt])
    } else {
      echo`📋 Showing all logs...`
      await dockerCompose(['logs', '-f'])
    }
    break

  case 'status':
    echo`📊 Container status for ${env} environment:`
    await dockerCompose(['ps'])
    echo``
    echo`🔍 Resource usage:`
    await $`docker stats --no-stream --format "table {{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}"`
    break

  case 'vault':
    await handleVaultCommand(opt, dockerCompose, config.envFile)
    break

  default:
    echo`❌ Invalid command: ${cmd}. Use: up, down, logs, status, vault`
    process.exit(1)
}

// Vault command handler
async function handleVaultCommand(vaultCmd, dockerCompose, envFile) {
  switch (vaultCmd) {
    case 'init':
      echo`🔐 Initializing Vault...`
      await dockerCompose(['exec', 'vault', 'vault', 'operator', 'init'])
      break

    case 'ui':
      echo`🌐 Opening Vault UI...`
      // Read VAULT_ADDR from env file
      const envContent = await fs.readFile(envFile, 'utf8')
      const vaultAddrMatch = envContent.match(/VAULT_ADDR=(.+)/)
      const vaultAddr = vaultAddrMatch ? vaultAddrMatch[1].replace(/"/g, '') : 'http://localhost:8200'
      echo`Vault UI available at: ${vaultAddr}/ui`
      break

    case 'status':
      echo`🔍 Vault status:`
      await dockerCompose(['exec', 'vault', 'vault', 'status'])
      break

    default:
      echo`❌ Invalid vault command. Use: init, ui, status`
      process.exit(1)
  }
}

echo`🎉 Operation completed successfully!`
