import React from 'react';
import type { GenreStat } from '../types';
import { formatMinLong } from '../utils/format';

interface GenreTableProps {
    stats: GenreStat[];
}

export const GenreTable: React.FC<GenreTableProps> = ({ stats }) => {
    return (
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
                    {stats.slice(0, 10).map((g) => (
                        <tr key={g.Name}>
                            <td>{g.Name}</td>
                            <td>{formatMinLong(g.DurationMin)}</td>
                            <td>{g.Share.toFixed(1)}%</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};
