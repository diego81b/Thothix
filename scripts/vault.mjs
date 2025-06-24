#!/usr/bin/env zx

/**
 * Thothix Vault Management Script
 * Unified script for vault operations: sync, init, and cleanup
 * Usage:
 *   zx vault.mjs --sync        # Sync secrets to Vault
 *   zx vault.mjs --init        # Initialize Vault + sync secrets
 *   zx vault.mjs --cleanup     # Clean temporary files
 */

import { $, chalk, echo, fs, path as zxPath } from 'zx';
import { readFileSync } from 'fs';
import { writeFile, unlink, readdir } from 'fs/promises';
import path from 'path';

// Configure for Windows compatibility
if (process.platform === 'win32') {
	$.shell = 'cmd.exe';
	$.prefix = '';
}

$.verbose = true;

const VAULT_MOUNT = 'thothix';

// Parse command line arguments
const isInitMode = process.argv.includes('--init');
const isSyncMode =
	process.argv.includes('--sync') ||
	(!isInitMode && !process.argv.includes('--cleanup'));
const isCleanupMode = process.argv.includes('--cleanup');

// Utility function for delays
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * Cleanup temporary files function
 */
async function performCleanup() {
	console.log('🧹 Thothix Temporary Files Cleanup');
	console.log('===================================');

	const tempPatterns = ['tmp-*-secrets.json', 'tmp-*-policy.hcl'];

	try {
		const files = await readdir('.');
		const tempFiles = [];

		for (const file of files) {
			for (const pattern of tempPatterns) {
				const regex = new RegExp(pattern.replace('*', '.*'));
				if (regex.test(file)) {
					tempFiles.push(file);
					break;
				}
			}
		}

		if (tempFiles.length === 0) {
			console.log('✅ No temporary files found - workspace is clean!');
			return;
		}

		console.log(`Found ${tempFiles.length} temporary file(s):`);
		tempFiles.forEach((file) => console.log(`  📄 ${file}`));

		console.log('\n🗑️  Cleaning up...');

		for (const file of tempFiles) {
			try {
				await unlink(file);
				console.log(`✅ Removed: ${file}`);
			} catch (error) {
				console.log(`❌ Failed to remove ${file}: ${error.message}`);
			}
		}

		console.log('\n🎉 Cleanup completed!');
	} catch (error) {
		console.error('❌ Cleanup failed:', error.message);
		process.exit(1);
	}
}

// Track temporary files for cleanup
const tempFiles = [];

// Cleanup function for temporary files
const cleanupTempFiles = async () => {
	for (const file of tempFiles) {
		try {
			await unlink(file);
			console.log(`🗑️ Cleaned up temporary file: ${file}`);
		} catch (error) {
			// File might already be deleted, ignore
		}
	}
	tempFiles.length = 0; // Clear the array
};

// Register cleanup handlers
process.on('exit', () => {
	// Synchronous cleanup on exit
	tempFiles.forEach((file) => {
		try {
			require('fs').unlinkSync(file);
		} catch (error) {
			// Ignore errors during exit cleanup
		}
	});
});

process.on('SIGINT', async () => {
	console.log('\n🛑 Script interrupted, cleaning up...');
	await cleanupTempFiles();
	process.exit(0);
});

process.on('SIGTERM', async () => {
	console.log('\n🛑 Script terminated, cleaning up...');
	await cleanupTempFiles();
	process.exit(0);
});

process.on('uncaughtException', async (error) => {
	console.error('\n❌ Uncaught exception:', error.message);
	await cleanupTempFiles();
	process.exit(1);
});

// Get environment
const environment =
	process.env.ENVIRONMENT || process.env.NODE_ENV || 'development';

// Handle cleanup mode first
if (isCleanupMode) {
	await performCleanup();
	process.exit(0);
}

// Display operation mode
console.log(
	`🔐 Thothix Vault ${isInitMode ? 'Init' : 'Sync'} - ${
		isInitMode ? 'Initializing' : 'Synchronizing'
	} .env secrets with Vault`,
);
console.log(`📦 Environment: ${environment}`);

// Read .env file
const envPath = path.join(process.cwd(), '.env');
let envContent;
try {
	envContent = readFileSync(envPath, 'utf8');
} catch (error) {
	console.error('❌ Cannot read .env file:', error.message);
	process.exit(1);
}

// Extract Vault configuration from .env
const vaultConfig = {};
envContent.split('\n').forEach((line) => {
	line = line.trim();
	if (line && !line.startsWith('#') && line.includes('VAULT_')) {
		const [key, ...valueParts] = line.split('=');
		if (key && valueParts.length > 0) {
			vaultConfig[key.trim()] = valueParts
				.join('=')
				.trim()
				.replace(/^['"]|['"]$/g, '');
		}
	}
});

const VAULT_ADDR = vaultConfig.VAULT_ADDR || 'http://localhost:8200';
const VAULT_TOKEN = vaultConfig.VAULT_ROOT_TOKEN || vaultConfig.VAULT_APP_TOKEN;

console.log('🔧 Vault Config:', {
	addr: VAULT_ADDR,
	mount: VAULT_MOUNT,
	hasToken: !!VAULT_TOKEN,
});

/**
 * Parse .env file into sections based on comment headers
 * Format: # :folder_name - Description
 * Only sections with :folder_name will be synced to Vault
 */
function parseEnvSections(envContent) {
	const sections = {};
	let currentSection = null;
	let currentSectionName = null;

	envContent.split('\n').forEach((line) => {
		line = line.trim();

		// Detect section headers: # :folder_name - Description
		if (line.startsWith('#')) {
			const match = line.match(/^#\s*:([a-zA-Z0-9_]+)\s*-\s*(.*)$/);

			if (match) {
				// Found a Vault section header
				const [, folderName, description] = match;
				currentSectionName = folderName.toLowerCase();

				currentSection = [];
				sections[currentSectionName] = {
					name: folderName,
					description: description.trim(),
					vars: currentSection,
				};

				console.log(`📁 Found Vault section: ${folderName} - ${description}`);
			} else {
				// Found a regular comment (not a Vault section) - reset tracking
				currentSection = null;
				currentSectionName = null;
				console.log(`🚫 Ignoring non-Vault section: ${line}`);
			}
		}
		// Parse variables ONLY if we're in a Vault section
		else if (currentSection !== null && line && !line.startsWith('#')) {
			const [key, ...valueParts] = line.split('=');
			if (key && valueParts.length > 0) {
				const value = valueParts
					.join('=')
					.trim()
					.replace(/^['"]|['"]$/g, '');

				console.log(
					`   📝 Adding to ${currentSectionName}: ${key.trim()} = ${value}`,
				);

				currentSection.push({
					key: key.trim(),
					value,
					vaultKey: key.trim().toLowerCase(), // Vault key naming
				});
			}
		}
		// Skip variables that are not in a Vault section
		else if (currentSection === null && line && !line.startsWith('#')) {
			console.log(
				`🚫 Skipping variable outside Vault sections: ${line.split('=')[0]}`,
			);
		}
	});

	return sections;
}

// Parse .env variables into sections
const envSections = parseEnvSections(envContent);
console.log('📋 Found sections to sync to Vault:', Object.keys(envSections));

// Check if vault container is running and wait for it if in init mode
async function waitForVault() {
	const maxRetries = isInitMode ? 30 : 3;
	const retryInterval = isInitMode ? 5000 : 1000;

	for (let i = 1; i <= maxRetries; i++) {
		try {
			await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault status"`;
			console.log('✅ Vault is accessible');
			return true;
		} catch (error) {
			if (i === maxRetries) {
				console.error(
					'❌ Cannot access Vault. Make sure Docker containers are running.',
				);
				console.error('   Try: docker compose up -d');
				console.error('   Vault error:', error.message);
				return false;
			}

			if (isInitMode) {
				console.log(`🔄 Waiting for Vault... (${i}/${maxRetries})`);
				await sleep(retryInterval);
			}
		}
	}
	return false;
}

async function initializeVault() {
	console.log('\n🏗️  Initializing Vault infrastructure...');

	// Check if secrets engine already exists
	try {
		const result =
			await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault secrets list"`;

		if (result.stdout.includes(`${VAULT_MOUNT}/`)) {
			console.log(`✅ KV secrets engine '${VAULT_MOUNT}/' already exists`);
		} else {
			console.log(`🔐 Enabling KV secrets engine at '${VAULT_MOUNT}/'...`);
			await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault secrets enable -path=${VAULT_MOUNT} kv-v2"`;
			console.log(`✅ KV secrets engine '${VAULT_MOUNT}/' created`);
		}
	} catch (error) {
		console.error('❌ Failed to setup secrets engine:', error.message);
		return false;
	}

	// Create policies
	console.log('\n🔐 Creating Vault policies...');

	// Define policies as simple strings without template literals
	const appPolicy = `# Read-only access to application secrets
path "${VAULT_MOUNT}/data/*" {
  capabilities = ["read"]
}

# Allow token renewal
path "auth/token/renew-self" {
  capabilities = ["update"]
}

# Allow token lookup
path "auth/token/lookup-self" {
  capabilities = ["read"]
}`;

	const readonlyPolicy = `# Read-only access for monitoring/debugging
path "${VAULT_MOUNT}/data/*" {
  capabilities = ["read"]
}`;

	try {
		// Create policy files locally and copy them to vault container
		const appPolicyContent = `# Read-only access to application secrets
path "${VAULT_MOUNT}/data/*" {
  capabilities = ["read"]
}

# Allow token renewal
path "auth/token/renew-self" {
  capabilities = ["update"]
}

# Allow token lookup
path "auth/token/lookup-self" {
  capabilities = ["read"]
}`;

		const readonlyPolicyContent = `# Read-only access for monitoring/debugging
path "${VAULT_MOUNT}/data/*" {
  capabilities = ["read"]
}`;

		// Write policy files locally
		const appPolicyFile = './tmp-app-policy.hcl';
		const readonlyPolicyFile = './tmp-readonly-policy.hcl';

		// Track temp files for cleanup
		tempFiles.push(appPolicyFile, readonlyPolicyFile);

		await writeFile(appPolicyFile, appPolicyContent);
		await writeFile(readonlyPolicyFile, readonlyPolicyContent);

		// Copy policies to vault container and apply them
		await $`docker compose cp ${appPolicyFile} vault:/tmp/thothix-app-policy.hcl`;
		await $`docker compose cp ${readonlyPolicyFile} vault:/tmp/thothix-readonly-policy.hcl`;

		// Apply policies
		await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault policy write thothix-app /tmp/thothix-app-policy.hcl"`;
		console.log('✅ Created thothix-app policy');

		await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault policy write thothix-readonly /tmp/thothix-readonly-policy.hcl"`;
		console.log('✅ Created thothix-readonly policy');

		// Clean up local temporary files
		await unlink(appPolicyFile);
		await unlink(readonlyPolicyFile);

		// Remove from tracking array
		tempFiles.splice(tempFiles.indexOf(appPolicyFile), 1);
		tempFiles.splice(tempFiles.indexOf(readonlyPolicyFile), 1);

		// Try to clean up remote temporary files (non-critical)
		try {
			await $`docker compose exec vault sh -c "rm -f /tmp/thothix-app-policy.hcl /tmp/thothix-readonly-policy.hcl"`;
		} catch (cleanupError) {
			console.log(
				'⚠️  Warning: Could not clean up temporary files in vault container (non-critical)',
			);
		}
	} catch (error) {
		console.error('❌ Failed to create policies:', error.message);
		return false;
	}

	// Create application tokens
	console.log('\n🔑 Creating application tokens...');

	try {
		const appTokenResult =
			await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault token create -policy=thothix-app -ttl=8760h -renewable=true -display-name='thothix-app-token-${environment}' -format=json"`;
		const appToken = JSON.parse(appTokenResult.stdout).auth.client_token;

		const readonlyTokenResult =
			await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault token create -policy=thothix-readonly -ttl=168h -renewable=true -display-name='thothix-readonly-token-${environment}' -format=json"`;
		const readonlyToken = JSON.parse(readonlyTokenResult.stdout).auth
			.client_token;

		console.log('\n✅ Vault initialized successfully!');
		console.log('\n📋 Configuration Summary:');
		console.log(`   Environment: ${environment}`);
		console.log(`   Mount Path: ${VAULT_MOUNT}`);
		console.log(`   Vault UI: http://localhost:8200`);
		console.log('\n🔑 Tokens created:');
		console.log(`   App Token: ${appToken}`);
		console.log(`   Readonly Token: ${readonlyToken}`);
		console.log('\n⚠️  IMPORTANT: Save these tokens securely!');
		console.log(`   Add VAULT_APP_TOKEN=${appToken} to your .env file`);

		return true;
	} catch (error) {
		console.error('❌ Failed to create tokens:', error.message);
		return false;
	}
}

// Check if vault container is running
if (!(await waitForVault())) {
	process.exit(1);
}

// Initialize Vault if in init mode
if (isInitMode) {
	if (!(await initializeVault())) {
		console.error('❌ Vault initialization failed');
		process.exit(1);
	}
}

// Sync each section to Vault
for (const [sectionName, section] of Object.entries(envSections)) {
	console.log(
		`\n🔐 Updating ${section.name} secrets (${section.description})...`,
	);
	try {
		// Build vault kv put command with individual key-value pairs
		const vaultPath = `${VAULT_MOUNT}/${sectionName}`;

		// Create a JSON object with all key-value pairs
		const secretData = {};
		section.vars.forEach((v) => {
			secretData[v.vaultKey] = v.value;
		});

		// Write the data as JSON to a temporary file and use vault kv put with @file
		const tempFile = `./tmp-${sectionName}-secrets.json`;

		// Track temp file for cleanup
		tempFiles.push(tempFile);

		await writeFile(tempFile, JSON.stringify(secretData, null, 2));

		// Copy file to vault container and use it
		await $`docker compose cp ${tempFile} vault:/tmp/${sectionName}-secrets.json`;
		await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault kv put ${vaultPath} @/tmp/${sectionName}-secrets.json"`;

		// Clean up local temp file
		await unlink(tempFile);

		// Remove from tracking array
		const fileIndex = tempFiles.indexOf(tempFile);
		if (fileIndex > -1) {
			tempFiles.splice(fileIndex, 1);
		}
		try {
			await $`docker compose exec vault sh -c "rm -f /tmp/${sectionName}-secrets.json"`;
		} catch (cleanupError) {
			// Ignore cleanup errors
		}

		console.log(`✅ ${section.name} secrets updated at ${vaultPath}`);
	} catch (error) {
		console.error(
			`❌ Failed to update ${section.name} secrets:`,
			error.message,
		);
	}
}

// Display current Vault contents
console.log('\n📋 Current Vault secrets:');
for (const [sectionName, section] of Object.entries(envSections)) {
	console.log(`\n🔐 ${section.name}:`);
	try {
		await $`docker compose exec vault sh -c "export VAULT_ADDR=${VAULT_ADDR} && export VAULT_TOKEN=${VAULT_TOKEN} && vault kv get ${VAULT_MOUNT}/${sectionName}"`;
	} catch (error) {
		console.error(`Cannot retrieve ${section.name} secrets`);
	}
}

console.log(
	`\n✅ Vault ${
		isInitMode ? 'initialization and sync' : 'synchronization'
	} completed!`,
);
console.log('\n📖 Next steps:');
console.log('   1. Check Vault UI at http://localhost:8200');
console.log('   2. Verify application can read secrets from Vault');
console.log('   3. Test with: npm run dev');

if (isInitMode) {
	console.log('\n🔍 Test secret retrieval:');
	console.log(`   vault kv get ${VAULT_MOUNT}/database`);
	console.log(`   vault kv get ${VAULT_MOUNT}/clerk`);
	console.log(`   vault kv get ${VAULT_MOUNT}/app`);
}

// Final cleanup of any remaining temp files
await cleanupTempFiles();
