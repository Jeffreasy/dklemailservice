// WebSocket Test Script voor DKL Email Service
// Test de /ws/steps endpoint met verschillende scenarios

const WebSocket = require('ws');

// Test configuratie
const WS_URL = 'ws://localhost:8080/ws/steps';
const TEST_TIMEOUT = 5000;

// Kleuren voor console output
const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
    console.log(`${colors[color]}${message}${colors.reset}`);
}

function logSuccess(message) {
    log(`✓ ${message}`, 'green');
}

function logError(message) {
    log(`✗ ${message}`, 'red');
}

function logInfo(message) {
    log(`ℹ ${message}`, 'cyan');
}

function logWarning(message) {
    log(`⚠ ${message}`, 'yellow');
}

// Test 1: Basis WebSocket connectie (anoniem)
async function testBasicConnection() {
    return new Promise((resolve, reject) => {
        log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        log('TEST 1: Basis WebSocket Connectie (Anoniem)', 'blue');
        log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        
        const ws = new WebSocket(WS_URL);
        let welcomeReceived = false;

        const timeout = setTimeout(() => {
            ws.close();
            if (!welcomeReceived) {
                logError('Timeout: Geen welcome message ontvangen');
                reject(new Error('Timeout'));
            }
        }, TEST_TIMEOUT);

        ws.on('open', () => {
            logSuccess('WebSocket connectie geopend');
        });

        ws.on('message', (data) => {
            const message = JSON.parse(data.toString());
            logInfo(`Bericht ontvangen: ${JSON.stringify(message, null, 2)}`);
            
            if (message.type === 'welcome') {
                welcomeReceived = true;
                logSuccess('Welcome message ontvangen');
                logInfo(`Available channels: ${message.available_channels?.join(', ')}`);
                
                clearTimeout(timeout);
                ws.close();
                resolve();
            }
        });

        ws.on('error', (error) => {
            clearTimeout(timeout);
            logError(`WebSocket error: ${error.message}`);
            reject(error);
        });

        ws.on('close', () => {
            logInfo('WebSocket connectie gesloten');
        });
    });
}

// Test 2: Subscribe naar channels
async function testSubscription() {
    return new Promise((resolve, reject) => {
        log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        log('TEST 2: Subscribe naar Channels', 'blue');
        log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        
        const ws = new WebSocket(WS_URL);
        let welcomeReceived = false;
        let subscriptionSent = false;

        const timeout = setTimeout(() => {
            ws.close();
            logWarning('Timeout: Test voltooid (geen updates te testen)');
            resolve(); // Geen error, want geen updates is OK
        }, TEST_TIMEOUT);

        ws.on('open', () => {
            logSuccess('WebSocket connectie geopend');
        });

        ws.on('message', (data) => {
            const message = JSON.parse(data.toString());
            
            if (message.type === 'welcome') {
                welcomeReceived = true;
                logSuccess('Welcome message ontvangen');
                
                // Subscribe naar alle channels
                const subscribeMsg = {
                    type: 'subscribe',
                    channels: ['step_updates', 'total_updates', 'leaderboard_updates']
                };
                
                logInfo(`Subscribing naar channels: ${subscribeMsg.channels.join(', ')}`);
                ws.send(JSON.stringify(subscribeMsg));
                subscriptionSent = true;
                logSuccess('Subscribe message verzonden');
                
                // Wacht nog even op mogelijke updates, dan sluiten
                setTimeout(() => {
                    clearTimeout(timeout);
                    ws.close();
                    resolve();
                }, 2000);
            } else if (message.type === 'ping') {
                logInfo('Ping ontvangen (keep-alive)');
            } else if (message.type === 'pong') {
                logInfo('Pong ontvangen');
            } else {
                logSuccess(`Update ontvangen: ${message.type}`);
                logInfo(`Data: ${JSON.stringify(message, null, 2)}`);
            }
        });

        ws.on('error', (error) => {
            clearTimeout(timeout);
            logError(`WebSocket error: ${error.message}`);
            reject(error);
        });

        ws.on('close', () => {
            logInfo('WebSocket connectie gesloten');
        });
    });
}

// Test 3: Ping/Pong functionaliteit
async function testPingPong() {
    return new Promise((resolve, reject) => {
        log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        log('TEST 3: Ping/Pong Functionaliteit', 'blue');
        log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        
        const ws = new WebSocket(WS_URL);
        let pongReceived = false;

        const timeout = setTimeout(() => {
            ws.close();
            if (!pongReceived) {
                logError('Timeout: Geen pong ontvangen');
                reject(new Error('No pong received'));
            }
        }, TEST_TIMEOUT);

        ws.on('open', () => {
            logSuccess('WebSocket connectie geopend');
        });

        ws.on('message', (data) => {
            const message = JSON.parse(data.toString());
            
            if (message.type === 'welcome') {
                logSuccess('Welcome message ontvangen');
                
                // Stuur ping
                const pingMsg = { type: 'ping' };
                logInfo('Sending ping...');
                ws.send(JSON.stringify(pingMsg));
            } else if (message.type === 'pong') {
                pongReceived = true;
                logSuccess('Pong ontvangen!');
                logInfo(`Pong timestamp: ${new Date(message.timestamp * 1000).toISOString()}`);
                
                clearTimeout(timeout);
                ws.close();
                resolve();
            }
        });

        ws.on('error', (error) => {
            clearTimeout(timeout);
            logError(`WebSocket error: ${error.message}`);
            reject(error);
        });

        ws.on('close', () => {
            logInfo('WebSocket connectie gesloten');
        });
    });
}

// Test 4: Multiple simultane connecties
async function testMultipleConnections() {
    return new Promise((resolve, reject) => {
        log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        log('TEST 4: Multiple Simultane Connecties', 'blue');
        log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
        
        const connections = [];
        const numConnections = 3;
        let connectedCount = 0;
        let welcomeCount = 0;

        const timeout = setTimeout(() => {
            connections.forEach(ws => ws.close());
            if (welcomeCount !== numConnections) {
                logError(`Niet alle connecties succesvol: ${welcomeCount}/${numConnections}`);
                reject(new Error('Not all connections successful'));
            }
        }, TEST_TIMEOUT);

        for (let i = 0; i < numConnections; i++) {
            const ws = new WebSocket(WS_URL);
            
            ws.on('open', () => {
                connectedCount++;
                logSuccess(`Connectie ${i + 1}/${numConnections} geopend`);
            });

            ws.on('message', (data) => {
                const message = JSON.parse(data.toString());
                if (message.type === 'welcome') {
                    welcomeCount++;
                    logSuccess(`Connectie ${i + 1}/${numConnections} welcome ontvangen`);
                    
                    if (welcomeCount === numConnections) {
                        logSuccess(`Alle ${numConnections} connecties succesvol!`);
                        clearTimeout(timeout);
                        
                        // Sluit alle connecties
                        setTimeout(() => {
                            connections.forEach(ws => ws.close());
                            resolve();
                        }, 500);
                    }
                }
            });

            ws.on('error', (error) => {
                clearTimeout(timeout);
                logError(`Connectie ${i + 1} error: ${error.message}`);
                reject(error);
            });

            connections.push(ws);
        }
    });
}

// Test 5: WebSocket Stats Endpoint (admin)
async function testStatsEndpoint() {
    log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
    log('TEST 5: WebSocket Stats Endpoint', 'blue');
    log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'blue');
    
    try {
        const fetch = (await import('node-fetch')).default;
        const response = await fetch('http://localhost:8080/api/ws/stats');
        
        if (response.status === 401) {
            logWarning('Stats endpoint requires authentication (expected)');
            logInfo('Status: 401 Unauthorized (correct behavior)');
            return;
        }
        
        const data = await response.json();
        logSuccess('Stats endpoint bereikbaar');
        logInfo(`Stats: ${JSON.stringify(data, null, 2)}`);
    } catch (error) {
        logError(`Stats endpoint test failed: ${error.message}`);
        throw error;
    }
}

// Voer alle tests uit
async function runAllTests() {
    log('\n╔════════════════════════════════════════════════════════════╗', 'cyan');
    log('║        DKL Email Service - WebSocket Test Suite           ║', 'cyan');
    log('╚════════════════════════════════════════════════════════════╝', 'cyan');
    
    const tests = [
        { name: 'Basic Connection', fn: testBasicConnection },
        { name: 'Subscription', fn: testSubscription },
        { name: 'Ping/Pong', fn: testPingPong },
        { name: 'Multiple Connections', fn: testMultipleConnections },
        { name: 'Stats Endpoint', fn: testStatsEndpoint }
    ];

    let passed = 0;
    let failed = 0;

    for (const test of tests) {
        try {
            await test.fn();
            passed++;
            logSuccess(`${test.name} test PASSED`);
        } catch (error) {
            failed++;
            logError(`${test.name} test FAILED: ${error.message}`);
        }
    }

    // Samenvatting
    log('\n╔════════════════════════════════════════════════════════════╗', 'cyan');
    log('║                      Test Samenvatting                     ║', 'cyan');
    log('╚════════════════════════════════════════════════════════════╝', 'cyan');
    log(`\nTotaal tests: ${tests.length}`, 'blue');
    log(`Passed: ${passed}`, 'green');
    log(`Failed: ${failed}`, failed > 0 ? 'red' : 'green');
    
    if (failed === 0) {
        log('\n🎉 Alle tests zijn geslaagd!', 'green');
        process.exit(0);
    } else {
        log(`\n⚠️  ${failed} test(s) gefaald`, 'red');
        process.exit(1);
    }
}

// Start de tests
runAllTests().catch(error => {
    logError(`Test suite failed: ${error.message}`);
    process.exit(1);
});