#!/usr/bin/env zx

// Thothix Development Script - Cross-platform with Zx
// Usage: zx scripts/dev.mjs [action]

import { $, argv, echo, cd, path, fs } from 'zx';

// Configure for Windows compatibility
if (process.platform === 'win32') {
	$.shell = 'powershell.exe';
	$.prefix = '';
} else {
	$.shell = 'bash';
	$.prefix = 'set -euo pipefail;';
}

$.verbose = true;

const action = argv._[0] || 'help';

echo`🔧 Thothix Development Script - Action: ${action}`;

// Helper functions
const format = async () => {
	echo`📝 Formatting Go code...`;
	await cd('./backend');
	await $`gofmt -w .`;
	echo`✅ Formatting completed`;
};

const lint = async () => {
	echo`🔍 Running golangci-lint...`;
	await cd('./backend');
	try {
		await $`golangci-lint run --timeout=3m`;
		echo`✅ Linting passed`;
	} catch (error) {
		echo`❌ Linting failed`;
		process.exit(1);
	}
};

const test = async () => {
	// Simple debug detection: check for --debug flag
	const isDebug =
		argv.debug ||
		argv._.includes('--debug') ||
		process.argv.includes('--debug');

	if (isDebug) {
		echo`🔬 Running tests with maximum debug output...`;
	} else {
		echo`🧪 Running tests...`;
	}

	await cd('./backend');
	echo`🎨 Using gotestsum for colored output`;

	// Start a progress indicator
	const progressChars = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
	let progressIndex = 0;

	const progressInterval = setInterval(() => {
		process.stdout.write(`\r${progressChars[progressIndex]} Running tests... `);
		progressIndex = (progressIndex + 1) % progressChars.length;
	}, 100);

	try {
		// Enhanced environment variables for colored output across platforms
		process.env.FORCE_COLOR = '1';
		process.env.TERM = 'xterm-256color';
		process.env.COLORTERM = 'truecolor';
		delete process.env.NO_COLOR; // Remove NO_COLOR to avoid warnings

		// Windows-specific color support
		if (process.platform === 'win32') {
			process.env.ANSICON = '1';
			process.env.ConEmuANSI = 'ON';
		}

		// Use more informative gotestsum formats
		const gotestsumFormat = isDebug ? 'testdox' : 'pkgname-and-test-fails';

		if (isDebug) {
			// Skip -race on Windows due to CGO requirements
			if (process.platform === 'win32') {
				await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ./... -v -timeout=300s -count=1 -failfast -cover`;
			} else {
				await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ./... -v -timeout=300s -count=1 -failfast -race -cover`;
			}
		} else {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ./... -timeout=300s -count=1 -failfast`;
		}

		if (isDebug) {
			echo`✅ Tests passed with debug info`;
		} else {
			echo`✅ Tests passed`;
		}

		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write('\r✅ Tests completed successfully!           \n');
	} catch (error) {
		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write('\r❌ Tests failed!                          \n');

		echo`❌ Tests failed`;

		if (isDebug) {
			echo`📋 Detailed error output:`;
			console.error('STDOUT:', error.stdout);
			console.error('STDERR:', error.stderr);
			console.error('Exit code:', error.exitCode);
		} else {
			echo`📋 Error details:`;
			console.error(error.stdout);
			console.error(error.stderr);
		}

		process.exit(1);
	}
};

const preCommit = async () => {
	echo`🚀 Running pre-commit checks...`;
	await format();

	echo`📋 Adding formatted files to git...`;
	await cd('../'); // Return to project root
	await $`git add backend/`;

	await lint();
	await test();
	echo`✅ Pre-commit checks completed`;
};

const showHelp = () => {
	echo`
Usage: zx scripts/dev.mjs [action] [flags]

Actions:
  format      - Format Go code with gofmt
  lint        - Run golangci-lint
  test        - Run Go tests
              --debug  Run with maximum debug output and coverage
  pre-commit  - Run format + git add + lint + test
  all         - Same as pre-commit

Examples:
  zx scripts/dev.mjs format
  zx scripts/dev.mjs test
  zx scripts/dev.mjs test --debug
  zx scripts/dev.mjs pre-commit
  npm run format          # Via NPM script
  npm run test            # Via NPM script (normal)
  npm run test:debug      # Via NPM script (with debug output)

Requirements:
  - Go installed
  - gotestsum installed (go install gotest.tools/gotestsum@latest)
  - golangci-lint installed
  - Git repository
  - Docker (for test containers)
`;
};

// Main execution
switch (action) {
	case 'format':
		await format();
		break;
	case 'lint':
		await lint();
		break;
	case 'test':
		await test();
		break;
	case 'pre-commit':
	case 'all':
		await preCommit();
		break;
	case 'help':
	default:
		showHelp();
		process.exit(action === 'help' ? 0 : 1);
}

echo`🎉 Script completed successfully!`;
