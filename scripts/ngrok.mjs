#!/usr/bin/env node

/**
 * Thothix - Ngrok Tunnel Script
 *
 * Starts ngrok HTTP tunnel for local webhook testing with Clerk
 * Automatically configures ngrok using environment variables from .env
 *
 * Usage:
 *   npm run ngrok              # Start tunnel on default port (30000)
 *   npm run ngrok -- 8080      # Start tunnel on custom port
 *   npm run ngrok -- --help    # Show help
 */

import { spawn } from 'child_process';
import { readFileSync, existsSync } from 'fs';
import { resolve } from 'path';

// Configuration
const DEFAULT_PORT = 30000;
const ENV_FILE = '.env';

// Colors for console output
const colors = {
	blue: '\x1b[34m',
	green: '\x1b[32m',
	yellow: '\x1b[33m',
	red: '\x1b[31m',
	reset: '\x1b[0m',
};

function log(message, color = 'reset') {
	console.log(`${colors[color]}${message}${colors.reset}`);
}

function showHelp() {
	log(
		'===================================================================',
		'blue',
	);
	log('🌐 Thothix - Ngrok Tunnel Script', 'blue');
	log(
		'===================================================================',
		'blue',
	);
	console.log();
	log('Description:', 'yellow');
	console.log(
		'  Starts ngrok HTTP tunnel for local webhook testing with Clerk',
	);
	console.log();
	log('Usage:', 'yellow');
	console.log(
		'  npm run ngrok              Start tunnel on default port (30000)',
	);
	console.log('  npm run ngrok -- 8080      Start tunnel on port 8080');
	console.log('  npm run ngrok -- --help    Show this help');
	console.log();
	log('Prerequisites:', 'yellow');
	console.log('  1. Install ngrok: https://ngrok.com/download');
	console.log(
		'  2. Get authtoken: https://dashboard.ngrok.com/get-started/your-authtoken',
	);
	console.log('  3. Set NGROK_AUTHTOKEN in .env file');
	console.log();
	log('Environment Variables (required in .env):', 'yellow');
	console.log('  NGROK_AUTHTOKEN     Your ngrok authentication token');
	console.log(
		'  NGROK_TUNNEL_URL    (optional) Static domain for paid accounts',
	);
	console.log();
	log('For webhook testing:', 'yellow');
	console.log('  1. Start this script to get tunnel URL');
	console.log('  2. If using free account: Copy HTTPS URL from ngrok output');
	console.log(
		'  3. If using paid account: Set NGROK_TUNNEL_URL in .env for static domain',
	);
	console.log('  4. Configure webhook in Clerk Dashboard:');
	console.log('     URL: https://YOUR_TUNNEL_URL/api/v1/auth/webhooks/clerk');
	console.log('     Events: user.created, user.updated, user.deleted');
	console.log();
}

function loadEnvFile() {
	if (!existsSync(ENV_FILE)) {
		log('❌ Error: .env file not found', 'red');
		log('💡 Run: cp .env.example .env', 'yellow');
		log('   Then configure NGROK_AUTHTOKEN in .env', 'yellow');
		process.exit(1);
	}

	const env = {};
	const envContent = readFileSync(ENV_FILE, 'utf8');

	for (const line of envContent.split('\n')) {
		const trimmed = line.trim();
		if (trimmed && !trimmed.startsWith('#')) {
			const [key, ...valueParts] = trimmed.split('=');
			if (key && valueParts.length > 0) {
				env[key] = valueParts.join('=');
			}
		}
	}

	return env;
}

function checkNgrokInstalled() {
	return new Promise((resolve) => {
		const check = spawn('ngrok', ['--version'], { stdio: 'pipe' });
		check.on('close', (code) => {
			resolve(code === 0);
		});
		check.on('error', () => {
			resolve(false);
		});
	});
}

async function configureNgrokAuthtoken(authtoken) {
	return new Promise((resolve, reject) => {
		log('🔑 Configuring ngrok authtoken...', 'yellow');

		const config = spawn('ngrok', ['config', 'add-authtoken', authtoken], {
			stdio: 'pipe',
		});

		config.on('close', (code) => {
			if (code === 0) {
				log('✅ Ngrok authtoken configured', 'green');
				resolve();
			} else {
				reject(new Error('Failed to configure ngrok authtoken'));
			}
		});

		config.on('error', (err) => {
			reject(err);
		});
	});
}

async function startNgrokTunnel(port, env) {
	return new Promise((resolve, reject) => {
		// Prepare ngrok command arguments
		const ngrokArgs = ['http', port.toString()];

		// Use static domain if configured (requires paid ngrok account)
		if (env.NGROK_TUNNEL_URL) {
			const domain = env.NGROK_TUNNEL_URL.replace('https://', '').replace(
				'http://',
				'',
			);
			ngrokArgs.push('--domain', domain);
			log(
				`🚀 Starting ngrok HTTP tunnel on port ${port} with domain ${domain}...`,
				'blue',
			);
			log(
				`🔗 Webhook URL: ${env.NGROK_TUNNEL_URL}/api/v1/auth/webhooks/clerk`,
				'yellow',
			);
		} else {
			log(`� Starting ngrok HTTP tunnel on port ${port}...`, 'blue');
			log(
				'�📝 Copy the HTTPS URL and update NGROK_TUNNEL_URL in .env',
				'yellow',
			);
			log(
				'🔗 Webhook URL will be: https://YOUR_TUNNEL_URL/api/v1/auth/webhooks/clerk',
				'yellow',
			);
		}

		console.log();
		log('Press Ctrl+C to stop the tunnel', 'yellow');
		log(
			'===================================================================',
			'blue',
		);
		console.log();

		const ngrok = spawn('ngrok', ngrokArgs, {
			stdio: 'inherit',
		});

		ngrok.on('close', (code) => {
			console.log();
			if (code === 0) {
				log('✅ Ngrok tunnel stopped', 'green');
			} else {
				log('❌ Ngrok tunnel exited with error', 'red');
			}
			resolve();
		});

		ngrok.on('error', (err) => {
			reject(err);
		});

		// Handle Ctrl+C
		process.on('SIGINT', () => {
			console.log();
			log('🛑 Stopping ngrok tunnel...', 'yellow');
			ngrok.kill('SIGINT');
		});
	});
}

async function main() {
	const args = process.argv.slice(2);

	// Handle help
	if (args.includes('--help') || args.includes('-h')) {
		showHelp();
		return;
	}

	// Parse port
	const port =
		args.length > 0 ? parseInt(args[0], 10) || DEFAULT_PORT : DEFAULT_PORT;

	log(
		'===================================================================',
		'blue',
	);
	log('🌐 Thothix - Starting Ngrok HTTP Tunnel', 'blue');
	log(
		'===================================================================',
		'blue',
	);
	console.log();

	try {
		// Load environment variables
		log('📋 Loading environment variables from .env...', 'yellow');
		const env = loadEnvFile();

		// Check NGROK_AUTHTOKEN
		if (!env.NGROK_AUTHTOKEN) {
			log('❌ Error: NGROK_AUTHTOKEN not found in .env', 'red');
			log(
				'💡 Get your authtoken from: https://dashboard.ngrok.com/get-started/your-authtoken',
				'yellow',
			);
			log('   Then add to .env: NGROK_AUTHTOKEN=your_token_here', 'yellow');
			process.exit(1);
		}

		// Check if ngrok is installed
		log('🔍 Checking ngrok installation...', 'yellow');
		const ngrokInstalled = await checkNgrokInstalled();

		if (!ngrokInstalled) {
			log('❌ Error: ngrok not found in PATH', 'red');
			log('💡 Install ngrok:', 'yellow');
			log('   - Download from: https://ngrok.com/download', 'yellow');
			log('   - Or via npm: npm install -g @ngrok/ngrok', 'yellow');
			log('   - Or via Chocolatey: choco install ngrok', 'yellow');
			process.exit(1);
		}

		log('✅ Environment variables loaded', 'green');
		log('✅ Ngrok found in PATH', 'green');
		console.log();

		// Configure ngrok authtoken
		await configureNgrokAuthtoken(env.NGROK_AUTHTOKEN);
		console.log();

		// Start ngrok tunnel
		await startNgrokTunnel(port, env);
	} catch (error) {
		console.log();
		log(`❌ Error: ${error.message}`, 'red');
		process.exit(1);
	}
}

// Run the script
main().catch((error) => {
	console.error('Unexpected error:', error);
	process.exit(1);
});
