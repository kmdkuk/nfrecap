import React from 'react';
import type { TitleStat } from '../types';
import { formatMinLong } from '../utils/format';

interface TitleTableProps {
    title: string;
    data: TitleStat[];
}

export const TitleTable: React.FC<TitleTableProps> = ({ title, data }) => {
    return (
        <div className="section">
            <h3>{title}</h3>
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
                    {data.slice(0, 10).map((t, i) => (
                        <tr key={i}>
                            <td>{t.Title}</td>
                            <td>{t.Type}</td>
                            <td>{formatMinLong(t.DurationMin)}</td>
                            <td>{t.Views}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};
