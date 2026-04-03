'use client';

import { useEffect, useRef } from 'react';
import { createChart, type IChartApi, ColorType } from 'lightweight-charts';

interface PortfolioChartProps {
  data: Array<{ time: string; value: number }>;
}

export default function PortfolioChart({ data }: PortfolioChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);

  useEffect(() => {
    if (!containerRef.current || data.length === 0) return;

    const chart = createChart(containerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: '#8B949E',
      },
      grid: {
        vertLines: { color: '#21262D' },
        horzLines: { color: '#21262D' },
      },
      width: containerRef.current.clientWidth,
      height: 250,
      timeScale: {
        borderColor: '#30363D',
        timeVisible: true,
      },
      rightPriceScale: {
        borderColor: '#30363D',
      },
      crosshair: {
        horzLine: { color: '#50E3C2', style: 3 },
        vertLine: { color: '#50E3C2', style: 3 },
      },
    });

    const series = chart.addAreaSeries({
      lineColor: '#50E3C2',
      topColor: 'rgba(80, 227, 194, 0.2)',
      bottomColor: 'rgba(80, 227, 194, 0)',
      lineWidth: 2,
    });

    const formatted = data.map((d) => ({
      time: d.time as any,
      value: d.value,
    }));

    series.setData(formatted);
    chart.timeScale().fitContent();
    chartRef.current = chart;

    const handleResize = () => {
      if (containerRef.current) {
        chart.applyOptions({ width: containerRef.current.clientWidth });
      }
    };
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chart.remove();
    };
  }, [data]);

  return <div ref={containerRef} />;
}
