export const formatMin = (min: number): string => {
    const h = Math.floor(min / 60);
    const m = min % 60;
    return `${h}h${m}m`; // Simplified format or matching previous `${h}h ${m}m`
};

export const formatMinLong = (min: number): string => {
    const h = Math.floor(min / 60);
    const m = min % 60;
    return `${h}h ${m}m`;
};
