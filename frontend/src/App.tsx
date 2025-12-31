import React, { useState } from 'react';
import './App.css';

interface UseMetric {
  Views: number;
  DurationMin: number;
}

interface Streak {
  Days: number;
  Start: string;
  End: string;
}

interface Gap {
  Days: number;
  Start: string;
  End: string;
}

interface GenreStat {
  Name: string;
  DurationMin: number;
  Views: number;
  Share: number;
}

interface Spike {
  Month: number;
  DurationMin: number;
}

interface TitleStat {
  Title: string;
  Type: string;
  DurationMin: number;
  Views: number;
}

interface SeriesStat {
  SeriesName: string;
  DurationMin: number;
  Views: number;
  SpanStart: string;
  SpanEnd: string;
}

interface UnresolvedItem {
  Title: string;
  Type: string;
  Views: number;
}

interface Stats {
  Year: number;
  GeneratedAt: string;
  SourceFile: string;
  TotalViews: number;
  TotalDurationMin: number;
  ActiveDays: number;
  TopStreaks: Streak[];
  MaxGap: Gap;
  MonthlyStats: Record<string, UseMetric>;
  WeekdayStats: Record<string, UseMetric>;
  GenreStats: GenreStat[];
  GenreMonthSpike: Record<string, Spike>;
  GenreSampleMovies: Record<string, string[]>;
  TopTitlesByDuration: TitleStat[];
  TopTitlesByViews: TitleStat[];
  TopSeriesByDuration: SeriesStat[];
  TopSeriesByViews: SeriesStat[];
  UnresolvedCount: number;
  UnresolvedList: UnresolvedItem[];
}

interface ApiResponse {
  recap: Stats;
}

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [data, setData] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    setLoading(true);
    setError(null);
    setData(null);

    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await fetch('/api/recap', {
        method: 'POST',
        body: formData,
      });
      if (!res.ok) {
        const txt = await res.text();
        throw new Error(txt || 'Upload failed');
      }
      const json: ApiResponse = await res.json();
      setData(json.recap);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatMin = (min: number) => {
    const h = Math.floor(min / 60);
    const m = min % 60;
    return `${h}h ${m}m`;
  };

  return (
    <div className="container">
      <header>
        <h1>Netflix Recap</h1>
        <p>Upload your NetflixViewingHistory.csv to see your stats.</p>
      </header>

      <div className="upload-section">
        <form onSubmit={handleSubmit}>
          <input type="file" accept=".csv" onChange={handleFileChange} />
          <button type="submit" disabled={!file || loading}>
            {loading ? 'Analyzing...' : 'Analyze'}
          </button>
        </form>
        {error && <div className="error">{error}</div>}
      </div>

      {data && (
        <div className="results fade-in">
          <h2>Stats for {data.Year}</h2>

          <div className="stats-grid">
            <div className="card">
              <h3>Total Time</h3>
              <p className="big-stat">{formatMin(data.TotalDurationMin)}</p>
              <p>{data.TotalViews} views</p>
            </div>
            <div className="card">
              <h3>Active Days</h3>
              <p className="big-stat">{data.ActiveDays}</p>
              <p>days watched</p>
            </div>
            <div className="card">
              <h3>Top Genre</h3>
              <p className="big-stat">{data.GenreStats[0]?.Name || '-'}</p>
              <p>{data.GenreStats[0] ? formatMin(data.GenreStats[0].DurationMin) : ''}</p>
            </div>
          </div>

          <div className="section">
            <h3>Top Genres</h3>
            <table>
              <thead>
                <tr>
                  <th>Genre</th>
                  <th>Duration</th>
                  <th>Share</th>
                </tr>
              </thead>
              <tbody>
                {data.GenreStats.slice(0, 10).map((g) => (
                  <tr key={g.Name}>
                    <td>{g.Name}</td>
                    <td>{formatMin(g.DurationMin)}</td>
                    <td>{g.Share.toFixed(1)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="section">
            <h3>Top Movies & TV (by Duration)</h3>
            <table>
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Type</th>
                  <th>Duration</th>
                  <th>Views</th>
                </tr>
              </thead>
              <tbody>
                {data.TopTitlesByDuration.slice(0, 10).map((t, i) => (
                  <tr key={i}>
                    <td>{t.Title}</td>
                    <td>{t.Type}</td>
                    <td>{formatMin(t.DurationMin)}</td>
                    <td>{t.Views}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
