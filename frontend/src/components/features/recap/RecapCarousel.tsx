import React, { useRef } from 'react';

interface RecapCarouselProps {
    children: React.ReactNode;
}

export const RecapCarousel: React.FC<RecapCarouselProps> = ({ children }) => {
    const scrollRef = useRef<HTMLDivElement>(null);

    // Optional: Function to handle scroll indication or buttons if needed later

    return (
        <div className="w-full flex justify-center py-8">
            <div className="relative w-full max-w-[340px]">
                <div
                    ref={scrollRef}
                    className="
                        flex gap-4 overflow-x-auto snap-x snap-mandatory 
                        px-10 pb-8
                        scrollbar-hide items-center
                        touch-pan-x overscroll-x-contain
                    "
                >
                    {children}
                </div>
                <p className="text-center text-xs text-gray-500 mt-2 opacity-50">
                    Swipe to explore →
                </p>
            </div>
        </div>
    );
};
