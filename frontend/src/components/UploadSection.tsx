import React from 'react';

interface UploadSectionProps {
    file: File | null;
    loading: boolean;
    error: string | null;
    onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSubmit: (e: React.FormEvent) => void;
}

export const UploadSection: React.FC<UploadSectionProps> = ({
    file,
    loading,
    error,
    onFileChange,
    onSubmit,
}) => {
    return (
        <div className="upload-section">
            <form onSubmit={onSubmit}>
                <input type="file" accept=".csv" onChange={onFileChange} />
                <button type="submit" disabled={!file || loading}>
                    {loading ? 'Analyzing...' : 'Analyze'}
                </button>
            </form>
            {error && <div className="error">{error}</div>}
        </div>
    );
};
