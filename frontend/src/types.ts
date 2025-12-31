export interface UseMetric {
    Views: number;
    DurationMin: number;
}

export interface Streak {
    Days: number;
    Start: string;
    End: string;
}

export interface Gap {
    Days: number;
    Start: string;
    End: string;
}

export interface GenreStat {
    Name: string;
    DurationMin: number;
    Views: number;
    Share: number;
}

export interface Spike {
    Month: number;
    DurationMin: number;
    Views: number;
}

export interface TitleStat {
    Title: string;
    Type: string;
    DurationMin: number;
    Views: number;
}

export interface SeriesStat {
    SeriesName: string;
    DurationMin: number;
    Views: number;
    SpanStart: string;
    SpanEnd: string;
}

export interface UnresolvedItem {
    Title: string;
    Type: string;
    Views: number;
}

export interface Stats {
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

export interface ApiResponse {
    recap: Stats;
}
