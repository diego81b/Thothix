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
	// Ensure we start from the project root and go to backend
	const projectRoot = process.cwd();
	if (!projectRoot.endsWith('Thothix')) {
		// If we're already in backend, go back to project root
		if (projectRoot.endsWith('backend')) {
			await cd('../');
		}
	}
	await cd('./backend');
	await $`gofmt -w .`;
	echo`✅ Formatting completed`;
};

const lint = async () => {
	echo`🔍 Running golangci-lint...`;
	// Ensure we start from the project root and go to backend
	const projectRoot = process.cwd();
	if (!projectRoot.endsWith('Thothix')) {
		// If we're already in backend, go back to project root
		if (projectRoot.endsWith('backend')) {
			await cd('../');
		}
	}
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
		echo`🔬 Running unit tests with maximum debug output...`;
	} else {
		echo`🧪 Running unit tests...`;
	}

	// Ensure we start from the project root and go to backend
	const projectRoot = process.cwd();
	if (!projectRoot.endsWith('Thothix')) {
		// If we're already in backend, go back to project root
		if (projectRoot.endsWith('backend')) {
			await cd('../');
		}
	}
	await cd('./backend');
	echo`🎨 Using gotestsum for colored output`;

	// Start a progress indicator
	const progressChars = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
	let progressIndex = 0;

	const progressInterval = setInterval(() => {
		process.stdout.write(
			`\r${progressChars[progressIndex]} Running unit tests... `,
		);
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

		// Get test configuration
		const testConfig = getTestConfig('unit');
		if (!testConfig.patterns || testConfig.patterns.length === 0) {
			echo`⚠️  No test patterns found`;
			return; // Exit early if no tests found
		}

		if (isDebug) {
			// Skip -race on Windows due to CGO requirements
			if (process.platform === 'win32') {
				await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -v -timeout=120s -count=1 -failfast -cover -run=${testConfig.runPattern}`;
			} else {
				await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -v -timeout=120s -count=1 -failfast -race -cover -run=${testConfig.runPattern}`;
			}
		} else {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -timeout=120s -count=1 -failfast -run=${testConfig.runPattern}`;
		}

		if (isDebug) {
			echo`✅ Unit tests passed with debug info`;
		} else {
			echo`✅ Unit tests passed`;
		}

		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write(
			'\r✅ Unit tests completed successfully!           \n',
		);
	} catch (error) {
		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write('\r❌ Unit tests failed!                          \n');

		echo`❌ Unit tests failed`;

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

const testIntegration = async () => {
	const isDebug =
		argv.debug ||
		argv._.includes('--debug') ||
		process.argv.includes('--debug');

	if (isDebug) {
		echo`🔗 Running integration tests with maximum debug output...`;
	} else {
		echo`🔗 Running integration tests...`;
	}

	// Ensure we start from the project root and go to backend
	const projectRoot = process.cwd();
	if (!projectRoot.endsWith('Thothix')) {
		// If we're already in backend, go back to project root
		if (projectRoot.endsWith('backend')) {
			await cd('../');
		}
	}
	await cd('./backend');

	// Start a progress indicator
	const progressChars = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
	let progressIndex = 0;

	const progressInterval = setInterval(() => {
		process.stdout.write(
			`\r${progressChars[progressIndex]} Running integration tests... `,
		);
		progressIndex = (progressIndex + 1) % progressChars.length;
	}, 100);

	try {
		// Enhanced environment variables for testcontainers
		process.env.FORCE_COLOR = '1';
		process.env.TESTCONTAINERS_RYUK_DISABLED = 'true'; // Disable reaper for Windows compatibility

		const gotestsumFormat = isDebug ? 'testdox' : 'pkgname-and-test-fails';

		// Get test configuration
		const testConfig = getTestConfig('integration');
		if (!testConfig.patterns || testConfig.patterns.length === 0) {
			echo`⚠️  No integration test patterns found`;
			return; // Exit early if no tests found
		}

		if (isDebug) {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -v -timeout=300s -count=1 -failfast -cover -run=${testConfig.runPattern}`;
		} else {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -timeout=300s -count=1 -failfast -run=${testConfig.runPattern}`;
		}

		if (isDebug) {
			echo`✅ Integration tests passed with debug info`;
		} else {
			echo`✅ Integration tests passed`;
		}

		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write(
			'\r✅ Integration tests completed successfully!      \n',
		);
	} catch (error) {
		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write(
			'\r❌ Integration tests failed!                      \n',
		);

		echo`❌ Integration tests failed`;

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

const testE2E = async () => {
	const isDebug =
		argv.debug ||
		argv._.includes('--debug') ||
		process.argv.includes('--debug');

	if (isDebug) {
		echo`🌐 Running E2E tests with maximum debug output...`;
	} else {
		echo`🌐 Running E2E tests...`;
	}

	// Ensure we start from the project root and go to backend
	const projectRoot = process.cwd();
	if (!projectRoot.endsWith('Thothix')) {
		// If we're already in backend, go back to project root
		if (projectRoot.endsWith('backend')) {
			await cd('../');
		}
	}
	await cd('./backend');

	// Start a progress indicator
	const progressChars = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
	let progressIndex = 0;

	const progressInterval = setInterval(() => {
		process.stdout.write(
			`\r${progressChars[progressIndex]} Running E2E tests... `,
		);
		progressIndex = (progressIndex + 1) % progressChars.length;
	}, 100);

	try {
		// Enhanced environment variables for testcontainers and Docker
		process.env.FORCE_COLOR = '1';
		process.env.TESTCONTAINERS_RYUK_DISABLED = 'true'; // Disable reaper for Windows compatibility
		process.env.TESTCONTAINERS_LOGGING = isDebug ? 'true' : 'false';

		const gotestsumFormat = isDebug ? 'testdox' : 'pkgname-and-test-fails';

		// Get test configuration
		const testConfig = getTestConfig('e2e');
		if (!testConfig.patterns || testConfig.patterns.length === 0) {
			echo`⚠️  No E2E test patterns found`;
			return; // Exit early if no tests found
		}

		if (isDebug) {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -v -timeout=600s -count=1 -failfast -cover`;
		} else {
			await $`gotestsum --format=${gotestsumFormat} --format-hivis -- ${testConfig.patterns} -timeout=600s -count=1 -failfast`;
		}

		if (isDebug) {
			echo`✅ E2E tests passed with debug info`;
		} else {
			echo`✅ E2E tests passed`;
		}

		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write(
			'\r✅ E2E tests completed successfully!            \n',
		);
	} catch (error) {
		// Clear progress indicator
		clearInterval(progressInterval);
		process.stdout.write(
			'\r❌ E2E tests failed!                            \n',
		);

		echo`❌ E2E tests failed`;

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

const testAll = async () => {
	echo`🧪 Running all tests (unit + integration + e2e)...`;

	await test(); // Unit tests
	await testIntegration(); // Integration tests
	await testE2E(); // E2E tests

	echo`✅ All test suites completed successfully!`;
};

const preCommit = async () => {
	echo`🚀 Running pre-commit checks...`;
	await format();

	echo`📋 Adding formatted files to git...`;
	await cd('../'); // Return to project root
	await $`git add backend/`;

	await lint();
	await test(); // Unit tests only for pre-commit
	echo`✅ Pre-commit checks completed`;
};

const preCommitFull = async () => {
	echo`🚀 Running full pre-commit checks (all tests)...`;
	await format();

	echo`📋 Adding formatted files to git...`;
	await cd('../'); // Return to project root
	await $`git add backend/`;

	await lint();
	await testAll(); // All tests for full pre-commit
	echo`✅ Full pre-commit checks completed`;
};

const showHelp = () => {
	echo`
Usage: zx scripts/dev.mjs [action] [flags]

Actions:
  format           - Format Go code with gofmt
  lint             - Run golangci-lint
  test             - Run unit tests only (fast)
                    --debug  Run with maximum debug output and coverage
  test:integration - Run integration tests (requires Docker)
                    --debug  Run with maximum debug output and coverage
  test:e2e         - Run end-to-end tests (requires Docker)
                    --debug  Run with maximum debug output and coverage
  test:all         - Run all tests (unit + integration + e2e)
                    --debug  Run with maximum debug output and coverage
  pre-commit       - Run format + git add + lint + unit tests (fast)
  pre-commit:full  - Run format + git add + lint + all tests (comprehensive)
  all              - Same as pre-commit:full

Test Categories:
  🧪 Unit Tests       - Fast isolated tests (domain models, DTOs, mappers, handlers)
  🔗 Integration Tests - Database integration tests (service layer across all domains)
  🌐 E2E Tests        - Full pipeline tests (HTTP → DB, requires containers, all domains)

Examples:
  zx scripts/dev.mjs format
  zx scripts/dev.mjs test                    # Unit tests only
  zx scripts/dev.mjs test:integration        # Integration tests
  zx scripts/dev.mjs test:e2e                # E2E tests
  zx scripts/dev.mjs test:all                # All test types
  zx scripts/dev.mjs test --debug            # Unit tests with debug output
  zx scripts/dev.mjs test:e2e --debug        # E2E tests with debug output
  zx scripts/dev.mjs pre-commit              # Fast pre-commit (unit tests)
  zx scripts/dev.mjs pre-commit:full         # Full pre-commit (all tests)
  npm run format                             # Via NPM script
  npm run test                               # Via NPM script (unit tests)
  npm run test:integration                   # Via NPM script (integration)
  npm run test:e2e                           # Via NPM script (e2e)
  npm run test:all                           # Via NPM script (all tests)

Requirements:
  - Go installed
  - gotestsum installed (go install gotest.tools/gotestsum@latest)
  - golangci-lint installed
  - Git repository
  - Docker (for integration and e2e tests)
  - Docker containers: PostgreSQL (for testcontainers)
`;
};

/*
 * EXTENDING THE TEST SCRIPT FOR NEW DOMAINS
 *
 * When adding a new domain (e.g., 'projects', 'chats', 'messages'), follow these steps:
 *
 * 1. Add the domain name to the 'knownDomains' array in getDomainConfig()
 *    Example: knownDomains: ['users', 'projects', 'chats']
 *
 * 2. Update test patterns in getTestConfig() if the new domain has different naming conventions:
 *    - Add to runPattern for unit tests if the test suite name differs from pattern
 *    - Add to runPattern for integration tests if the test suite name differs from pattern
 *
 * 3. Create the standard directory structure in backend/internal/[domain]/:
 *    - domain/     (business logic, entities)
 *    - dto/        (data transfer objects)
 *    - mappers/    (entity ↔ DTO converters)
 *    - handlers/   (HTTP handlers)
 *    - service/    (business services, both unit and integration tests)
 *    - e2e/        (end-to-end tests)
 *
 * 4. Follow the test naming conventions:
 *    - Unit tests in service: TestSuiteNameTestSuite (e.g., ProjectServiceTestSuite)
 *    - Integration tests in service: TestSuiteNameIntegrationTestSuite (e.g., ProjectServiceIntegrationTestSuite)
 *    - Unit test files: *_test.go (without "integration" in the name)
 *    - Integration test files: *_integration_test.go
 *
 * The script automatically distinguishes between unit and integration tests using centralized configuration:
 * - All domain patterns are managed in getDomainConfig()
 * - All test patterns and regex filters are managed in getTestConfig()
 * - Easy to extend for new domains by updating just the central configuration
 *
 * Example domain structure:
 * backend/internal/projects/
 * ├── domain/
 * │   ├── project.go
 * │   └── project_test.go
 * ├── dto/
 * │   ├── project_dto.go
 * │   └── project_dto_test.go
 * ├── mappers/
 * │   ├── project_mapper.go
 * │   └── project_mapper_test.go
 * ├── handlers/
 * │   ├── project_handler.go
 * │   ├── project_handler_test.go
 * │   └── routes.go
 * ├── service/
 * │   ├── project_service.go
 * │   ├── project_service_test.go (unit tests - ProjectServiceTestSuite)
 * │   └── project_service_integration_test.go (integration tests - ProjectServiceIntegrationTestSuite)
 * └── e2e/
 *     └── project_end_to_end_test.go
 */

// Helper function to get domain-specific configurations
const getDomainConfig = () => {
	// Central configuration for all domains
	// Add new domains here when they're created
	return {
		knownDomains: ['users'],
		// Add domain-specific patterns here if needed in the future
	};
};

// Helper function to build test patterns and configurations
const getTestConfig = (testType) => {
	const { knownDomains } = getDomainConfig();

	switch (testType) {
		case 'unit':
			const unitPatterns = [];
			const unitServicePatterns = [];

			for (const domain of knownDomains) {
				// Non-service patterns
				unitPatterns.push(
					`./internal/${domain}/domain/...`,
					`./internal/${domain}/dto/...`,
					`./internal/${domain}/mappers/...`,
					`./internal/${domain}/handlers/...`,
				);
				// Service patterns (handled separately)
				unitServicePatterns.push(`./internal/${domain}/service`);
			}

			// Add shared components
			unitPatterns.push(
				'./internal/shared/dto/...',
				'./internal/shared/handlers/...',
				'./internal/shared/middleware/...',
			);

			return {
				patterns: [...unitPatterns, ...unitServicePatterns],
				runPattern:
					'UserServiceTestSuite|^Test[^I].*|^TestUser[^S].*|^TestGet.*|^TestCreate.*|^TestUpdate.*|^TestDelete.*|^TestNew.*',
			};

		case 'integration':
			const integrationPatterns = [];

			for (const domain of knownDomains) {
				integrationPatterns.push(`./internal/${domain}/service`);
			}

			return {
				patterns: integrationPatterns,
				runPattern: 'UserServiceIntegrationTestSuite|Integration',
			};

		case 'e2e':
			const e2ePatterns = [];

			for (const domain of knownDomains) {
				e2ePatterns.push(`./internal/${domain}/e2e/...`);
			}

			return {
				patterns: e2ePatterns,
				runPattern: null, // No filter needed for e2e
			};

		default:
			return { patterns: [], runPattern: null };
	}
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
	case 'test:integration':
		await testIntegration();
		break;
	case 'test:e2e':
		await testE2E();
		break;
	case 'test:all':
		await testAll();
		break;
	case 'pre-commit':
		await preCommit();
		break;
	case 'pre-commit:full':
	case 'all':
		await preCommitFull();
		break;
	case 'help':
	default:
		showHelp();
		process.exit(action === 'help' ? 0 : 1);
}

echo`🎉 Script completed successfully!`;
