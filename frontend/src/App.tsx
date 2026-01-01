import React, { useState } from 'react';
import './App.css';
import type { ApiResponse, Stats } from './types';
import { UploadSection } from './components/UploadSection';
import { StatsSummary } from './components/StatsSummary';
import { GenreTable } from './components/GenreTable';
import { TitleTable } from './components/TitleTable';
import { RecapCard } from './components/features/recap/RecapCard';
import { RecapStory } from './components/features/recap/RecapStory';
import { Clock, Eye, Calendar, Film } from 'lucide-react';

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
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <header>
        <h1>Netflix Recap</h1>
        <p>Upload your NetflixViewingHistory.csv to see your stats.</p>
      </header>

      <UploadSection
        file={file}
        loading={loading}
        error={error}
        onFileChange={handleFileChange}
        onSubmit={handleSubmit}
      />

      {data && (
        <div className="results fade-in">
          <h2>Stats for {data.Year}</h2>

          <RecapStory>
            <RecapCard
              title="Total Views"
              value={data.TotalViews}
              category="Engagement"
              icon={Eye}
              variant="highlight"
              backgroundColor="#E50914"
            />
            <RecapCard
              title="Total Time"
              value={`${Math.round(data.TotalDurationMin / 60)}h`}
              category="Time Spent"
              icon={Clock}
            />
            <RecapCard
              title="Active Days"
              value={data.ActiveDays}
              category="Consistency"
              icon={Calendar}
            />
            {data.GenreStats.length > 0 && (
              <RecapCard
                title="Top Genre"
                value={data.GenreStats[0].Name}
                category="Favorite"
                icon={Film}
                variant="highlight"
                backgroundColor="#221f1f"
              />
            )}
          </RecapStory>

          <StatsSummary data={data} />

          <GenreTable stats={data.GenreStats} />

          <TitleTable
            title="Top Movies & TV (by Duration)"
            data={data.TopTitlesByDuration}
          />
        </div>
      )}
    </div>
  );
}

export default App;
