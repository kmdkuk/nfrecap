import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface RecapStoryProps {
    children: React.ReactNode[];
    onComplete?: () => void;
}

export const RecapStory: React.FC<RecapStoryProps> = ({ children, onComplete }) => {
    const [currentIndex, setCurrentIndex] = useState(0);
    const totalSlides = React.Children.count(children);
    const slides = React.Children.toArray(children);

    const goToNext = () => {
        if (currentIndex < totalSlides - 1) {
            setCurrentIndex(prev => prev + 1);
        } else {
            onComplete?.();
        }
    };

    const goToPrev = () => {
        if (currentIndex > 0) {
            setCurrentIndex(prev => prev - 1);
        }
    };

    // Keyboard navigation
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'ArrowRight' || e.key === ' ') {
                goToNext();
            } else if (e.key === 'ArrowLeft') {
                goToPrev();
            }
        };
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [currentIndex, totalSlides]);

    return (
        <div className="w-full flex justify-center py-8">
            <div className="relative w-full max-w-[360px] aspect-[9/16] bg-black rounded-2xl overflow-hidden shadow-2xl">

                {/* Progress Bars */}
                <div className="absolute top-0 left-0 right-0 z-20 flex gap-1 p-2">
                    {Array.from({ length: totalSlides }).map((_, idx) => (
                        <div key={idx} className="h-1 flex-1 bg-white/30 rounded-full overflow-hidden">
                            <div
                                className={`h-full bg-white transition-all duration-300 ${idx < currentIndex ? 'w-full' :
                                        idx === currentIndex ? 'w-full' : 'w-0'
                                    }`}
                            />
                        </div>
                    ))}
                </div>

                {/* Navigation Zones (Tap Areas) */}
                <div className="absolute inset-0 z-10 flex">
                    <div className="w-1/3 h-full" onClick={goToPrev} />
                    <div className="w-2/3 h-full" onClick={goToNext} />
                </div>

                {/* Content */}
                <div className="relative w-full h-full">
                    {slides.map((child, idx) => (
                        <div
                            key={idx}
                            className={`absolute inset-0 transition-opacity duration-300 flex items-center justify-center ${idx === currentIndex ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'
                                }`}
                        >
                            {/* Render child directly but stripped of some styles if needed, 
                        or assume child is a RecapCard. 
                        We wrap it to ensure it fits the story container. */}
                            <div className="w-full h-full transform scale-100 origin-center">
                                {child}
                            </div>
                        </div>
                    ))}
                </div>

                {/* Optional: Navigation Buttons for desktop UX helper */}
                <button
                    onClick={(e) => { e.stopPropagation(); goToPrev(); }}
                    className={`absolute left-2 top-1/2 -translate-y-1/2 z-20 p-2 rounded-full bg-black/20 text-white backdrop-blur-sm transition-opacity hover:bg-black/40 ${currentIndex === 0 ? 'opacity-0 pointer-events-none' : 'opacity-100'}`}
                >
                    <ChevronLeft className="w-6 h-6" />
                </button>
                <button
                    onClick={(e) => { e.stopPropagation(); goToNext(); }}
                    className={`absolute right-2 top-1/2 -translate-y-1/2 z-20 p-2 rounded-full bg-black/20 text-white backdrop-blur-sm transition-opacity hover:bg-black/40 ${currentIndex === totalSlides - 1 ? 'opacity-0 pointer-events-none' : 'opacity-100'}`}
                >
                    <ChevronRight className="w-6 h-6" />
                </button>

            </div>
        </div>
    );
};
