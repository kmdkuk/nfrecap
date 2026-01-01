import React from 'react';
import type { LucideIcon } from 'lucide-react';

interface RecapCardProps {
    title: string;
    value: string | number;
    category: string;
    icon?: LucideIcon;
    variant?: 'default' | 'highlight';
    backgroundImageUrl?: string; // Optional background image
    backgroundColor?: string; // Optional background color override
}

export const RecapCard: React.FC<RecapCardProps> = ({
    title,
    value,
    category,
    icon: Icon,
    variant = 'default',
    backgroundImageUrl,
    backgroundColor,
}) => {
    const isHighlight = variant === 'highlight';

    return (
        <div
            className={`
        relative overflow-hidden flex-shrink-0
        w-full h-full
        ${isHighlight ? 'bg-indigo-600 text-white' : 'bg-white text-gray-900'}
      `}
            style={{
                backgroundColor: backgroundColor,
                // Remove border radius if fitting into story container precisely, 
                // or keep it if design demands. Story container is already rounded.
                borderRadius: '0px',
            }}
        >
            {/* Background Image Layer */}
            {backgroundImageUrl && (
                <div className="absolute inset-0 z-0">
                    <img
                        src={backgroundImageUrl}
                        alt=""
                        className="w-full h-full object-cover opacity-30"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                </div>
            )}

            {/* Content Layer */}
            <div className="relative z-10 h-full flex flex-col justify-between p-6 sm:p-8">
                {/* Header */}
                <div className="flex items-start justify-between">
                    <span className={`
            text-xs font-semibold uppercase tracking-wider px-2 py-1 rounded-full
            ${isHighlight || backgroundImageUrl ? 'bg-white/20 text-white backdrop-blur-sm' : 'bg-gray-100 text-gray-600'}
          `}>
                        {category}
                    </span>
                    {Icon && <Icon className={`w-6 h-6 ${isHighlight || backgroundImageUrl ? 'text-white' : 'text-gray-400'}`} />}
                </div>

                {/* Main Value */}
                <div className="flex flex-col gap-2">
                    <h3 className={`text-lg font-medium opacity-90 ${isHighlight || backgroundImageUrl ? 'text-white' : 'text-gray-600'}`}>
                        {title}
                    </h3>
                    <p className={`text-4xl sm:text-5xl font-bold tracking-tight ${isHighlight || backgroundImageUrl ? 'text-white' : 'text-gray-900'}`}>
                        {value}
                    </p>
                </div>

                {/* Decorative elements or bottom info could go here */}
                <div className="h-2 w-12 rounded-full bg-current opacity-20" />
            </div>
        </div>
    );
};
