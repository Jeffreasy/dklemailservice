/**
 * Dashboard Component - Complete voorbeeld van Steps tracking met WebSocket
 * 
 * Features:
 * - Real-time step updates via WebSocket
 * - Leaderboard met live rankings
 * - Badge notifications
 * - Connection status indicator
 * - Manual step input
 * 
 * Usage:
 * ```tsx
 * import Dashboard from './Dashboard';
 * 
 * function App() {
 *   return <Dashboard userId="user-123" participantId="participant-456" />;
 * }
 * ```
 */

import React, { useState, useEffect } from 'react';
import { useStepsWebSocket } from './useStepsWebSocket';

interface DashboardProps {
  userId: string;
  participantId: string;
}

const Dashboard: React.FC<DashboardProps> = ({ userId, participantId }) => {
  const {
    connected,
    connectionState,
    latestUpdate,
    totalSteps,
    leaderboard,
    latestBadge,
    subscribe,
  } = useStepsWebSocket(userId, participantId, {
    debug: true,
    reconnectInterval: 1000,
    maxReconnectInterval: 30000,
  });

  const [stepInput, setStepInput] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [showBadgeNotification, setShowBadgeNotification] = useState<boolean>(false);

  // Subscribe to all channels on connect
  useEffect(() => {
    if (connected) {
      subscribe(['step_updates', 'total_updates', 'leaderboard_updates']);
    }
  }, [connected, subscribe]);

  // Show badge notification
  useEffect(() => {
    if (latestBadge) {
      setShowBadgeNotification(true);
      setTimeout(() => setShowBadgeNotification(false), 5000);
    }
  }, [latestBadge]);

  // Manual step update via REST API
  const updateSteps = async (delta: number) => {
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('/api/steps', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ steps: delta }),
      });

      if (!response.ok) {
        throw new Error('Failed to update steps');
      }

      // WebSocket will handle the UI update
      console.log('Steps updated successfully');
    } catch (error) {
      console.error('Error updating steps:', error);
      alert('Kon stappen niet bijwerken');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const delta = parseInt(stepInput, 10);
    if (isNaN(delta) || delta === 0) {
      alert('Voer een geldig aantal stappen in');
      return;
    }
    updateSteps(delta);
    setStepInput('');
  };

  return (
    <div className="dashboard">
      {/* Connection Status */}
      <div className="connection-status">
        <span className={`status-indicator ${connected ? 'online' : 'offline'}`}>
          {connected ? '🟢' : '🔴'}
        </span>
        <span className="status-text">
          {connectionState}
        </span>
      </div>

      {/* Badge Notification */}
      {showBadgeNotification && latestBadge && (
        <div className="badge-notification">
          <img src={latestBadge.badge_icon} alt={latestBadge.badge_name} />
          <div>
            <h3>Badge Verdiend! 🎉</h3>
            <p>{latestBadge.badge_name}</p>
            <p>+{latestBadge.points} punten</p>
          </div>
        </div>
      )}

      {/* Personal Stats */}
      <div className="personal-stats">
        <h2>Mijn Voortgang</h2>
        <div className="stats-grid">
          <div className="stat-card">
            <h3>Stappen</h3>
            <p className="stat-value">
              {latestUpdate?.steps.toLocaleString() || '0'}
            </p>
            {latestUpdate?.delta !== undefined && latestUpdate.delta !== 0 && (
              <p className="stat-delta">
                {latestUpdate.delta > 0 ? '+' : ''}
                {latestUpdate.delta.toLocaleString()}
              </p>
            )}
          </div>

          <div className="stat-card">
            <h3>Route</h3>
            <p className="stat-value">{latestUpdate?.route || 'N/A'}</p>
          </div>

          <div className="stat-card">
            <h3>Toegewezen Fonds</h3>
            <p className="stat-value">
              €{latestUpdate?.allocated_funds || 0}
            </p>
          </div>
        </div>
      </div>

      {/* Step Input Form */}
      <div className="step-input">
        <h3>Stappen Toevoegen</h3>
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <input
              type="number"
              value={stepInput}
              onChange={(e) => setStepInput(e.target.value)}
              placeholder="Aantal stappen"
              disabled={loading || !connected}
            />
            <button
              type="submit"
              disabled={loading || !connected}
              className="btn-primary"
            >
              {loading ? 'Bezig...' : 'Toevoegen'}
            </button>
          </div>
        </form>

        {/* Quick action buttons */}
        <div className="quick-actions">
          <button
            onClick={() => updateSteps(500)}
            disabled={loading || !connected}
            className="btn-secondary"
          >
            +500
          </button>
          <button
            onClick={() => updateSteps(1000)}
            disabled={loading || !connected}
            className="btn-secondary"
          >
            +1000
          </button>
          <button
            onClick={() => updateSteps(5000)}
            disabled={loading || !connected}
            className="btn-secondary"
          >
            +5000
          </button>
        </div>
      </div>

      {/* Global Stats */}
      <div className="global-stats">
        <h3>Totaal Stappen</h3>
        <p className="big-number">{totalSteps.toLocaleString()}</p>
      </div>

      {/* Leaderboard */}
      {leaderboard && (
        <div className="leaderboard">
          <h3>Top {leaderboard.top_n}</h3>
          <table>
            <thead>
              <tr>
                <th>Rank</th>
                <th>Naam</th>
                <th>Stappen</th>
                <th>Punten</th>
                <th>Totaal</th>
                <th>Route</th>
              </tr>
            </thead>
            <tbody>
              {leaderboard.entries.map((entry) => (
                <tr
                  key={entry.participant_id}
                  className={
                    entry.participant_id === participantId ? 'highlighted' : ''
                  }
                >
                  <td className="rank">
                    {entry.rank === 1 && '🥇'}
                    {entry.rank === 2 && '🥈'}
                    {entry.rank === 3 && '🥉'}
                    {entry.rank > 3 && entry.rank}
                  </td>
                  <td>{entry.naam}</td>
                  <td>{entry.steps.toLocaleString()}</td>
                  <td>{entry.achievement_points}</td>
                  <td><strong>{entry.total_score.toLocaleString()}</strong></td>
                  <td>{entry.route}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default Dashboard;

/**
 * Minimal CSS Voorbeeld (dashboard.css)
 */
const cssExample = `
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.status-indicator {
  font-size: 12px;
}

.badge-notification {
  position: fixed;
  top: 20px;
  right: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
  display: flex;
  gap: 15px;
  align-items: center;
  animation: slideIn 0.3s ease-out;
  z-index: 1000;
}

.badge-notification img {
  width: 60px;
  height: 60px;
}

@keyframes slideIn {
  from {
    transform: translateX(400px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin: 20px 0;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #667eea;
  margin: 10px 0;
}

.stat-delta {
  font-size: 18px;
  color: #4ade80;
}

.step-input {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin: 20px 0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.input-group input {
  flex: 1;
  padding: 12px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
}

.btn-primary {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #5568d3;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quick-actions {
  display: flex;
  gap: 10px;
}

.btn-secondary {
  background: #f3f4f6;
  border: 2px solid #e5e7eb;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e7eb;
  border-color: #667eea;
}

.big-number {
  font-size: 48px;
  font-weight: bold;
  color: #667eea;
  text-align: center;
  margin: 20px 0;
}

.leaderboard {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin: 20px 0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.leaderboard table {
  width: 100%;
  border-collapse: collapse;
}

.leaderboard th,
.leaderboard td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.leaderboard th {
  background: #f9fafb;
  font-weight: 600;
  color: #6b7280;
}

.leaderboard tr.highlighted {
  background: #fef3c7;
}

.leaderboard .rank {
  font-size: 20px;
  text-align: center;
}
`;

// Export CSS voor gebruik
export { cssExample };