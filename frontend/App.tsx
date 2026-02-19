import React, { useState, useCallback, useEffect } from "react";

type Game = {
  name: string;
  lastPlayed: string;
  minutesPlayed: number;
  appId: string;
  iconUrl: string;
};

function parseGames(raw: unknown): Game[] {
  if (!raw || !Array.isArray(raw)) return [];
  return raw.map((g) => {
    if (typeof g === "string") {
      try {
        return JSON.parse(g) as Game;
      } catch {
        return null;
      }
    }
    return (g as Game) ?? null;
  }).filter(Boolean) as Game[];
}

export default function App() {
  const [authLoading, setAuthLoading] = useState(true);
  const [user, setUser] = useState<{ steamId: string } | null>(null);
  const [loading, setLoading] = useState(false);
  const [games, setGames] = useState<Game[]>([]);
  const [revealed, setRevealed] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/auth/me", { credentials: "include" })
      .then((res) => {
        if (res.ok) return res.json();
        setUser(null);
        return null;
      })
      .then((data) => {
        if (data?.steamId) setUser({ steamId: data.steamId });
      })
      .finally(() => setAuthLoading(false));
  }, []);

  const fetchAndReveal = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const res = await fetch("/api/games/library", { credentials: "include" });
      if (res.status === 401) {
        setUser(null);
        setError("Session expired. Sign in again.");
        return;
      }
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `HTTP ${res.status}`);
      }
      const data = (await res.json()) as { games?: unknown };
      const list = parseGames(data.games);
      setGames(list);
      setRevealed(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  if (authLoading) {
    return (
      <div className="app">
        <div className="cta-wrap">
          <span className="cta cta--muted">Loading…</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="app">
        <div className="cta-wrap">
          <a href="/api/auth/steam/login" className="cta">
            Sign in with Steam
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <div className={`cta-wrap ${revealed ? "revealed" : ""}`}>
        <button
          type="button"
          className={`cta ${revealed ? "revealed" : ""}`}
          onClick={fetchAndReveal}
          disabled={loading}
        >
          {loading ? "Loading…" : revealed ? "Check again" : "Check your playtime"}
        </button>
      </div>

      <div className={`content ${revealed ? "visible" : ""}`}>
        {revealed && (
          <p className="auth-row">
            <a href="/api/auth/logout" className="auth-link">Sign out</a>
          </p>
        )}
        {error && (
          <p className="error">{error}</p>
        )}
        {revealed && !error && (
          <>
            <h2>Library</h2>
            {games.length === 0 ? (
              <p className="loading">No games in library.</p>
            ) : (
              <ul className="library">
                {games.map((game) => (
                  <li key={game.appId} className="game-card">
                    {game.iconUrl ? (
                      <img
                        src={game.iconUrl}
                        alt=""
                        className="game-icon"
                        width={48}
                        height={48}
                      />
                    ) : (
                      <div className="game-icon" aria-hidden />
                    )}
                    <div className="game-info">
                      <p className="game-name">{game.name}</p>
                      <p className="game-meta">
                        {game.minutesPlayed >= 0 && `${Math.round(game.minutesPlayed / 60)}h played`}
                        {game.lastPlayed && ` · Last: ${new Date(game.lastPlayed).toLocaleDateString()}`}
                      </p>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </>
        )}
      </div>
    </div>
  );
}
