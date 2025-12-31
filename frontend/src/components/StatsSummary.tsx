import React from 'react';
import type { Stats } from '../types';
import { formatMinLong } from '../utils/format';

interface StatsSummaryProps {
    data: Stats;
}

export const StatsSummary: React.FC<StatsSummaryProps> = ({ data }) => {
    const topGenreName = data.GenreStats[0]?.Name || '-';
    const topGenreDur = data.GenreStats[0]
        ? formatMinLong(data.GenreStats[0].DurationMin)
        : '';

    return (
        <div className="stats-grid">
            <div className="card">
                <h3>Total Time</h3>
                <p className="big-stat">{formatMinLong(data.TotalDurationMin)}</p>
                <p>{data.TotalViews} views</p>
            </div>
            <div className="card">
                <h3>Active Days</h3>
                <p className="big-stat">{data.ActiveDays}</p>
                <p>days watched</p>
            </div>
            <div className="card">
                <h3>Top Genre</h3>
                <p className="big-stat">{topGenreName}</p>
                <p>{topGenreDur}</p>
            </div>
        </div>
    );
};
