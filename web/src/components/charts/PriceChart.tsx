'use client';

import { useEffect, useRef } from 'react';
import { createChart, ColorType, type IChartApi, type ISeriesApi, type Time } from 'lightweight-charts';

interface PriceChartProps {
  data: Array<{
    time: string;
    open: number;
    high: number;
    low: number;
    close: number;
  }>;
  height?: number;
}

export default function PriceChart({ data, height = 300 }: PriceChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<'Candlestick'> | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const chart = createChart(containerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: '#8B949E',
        fontSize: 11,
      },
      grid: {
        vertLines: { color: '#21262D' },
        horzLines: { color: '#21262D' },
      },
      width: containerRef.current.clientWidth,
      height,
      crosshair: {
        vertLine: { color: '#50E3C2', width: 1, style: 2, labelBackgroundColor: '#50E3C2' },
        horzLine: { color: '#50E3C2', width: 1, style: 2, labelBackgroundColor: '#50E3C2' },
      },
      timeScale: {
        borderColor: '#30363D',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#30363D',
      },
    });

    const series = chart.addCandlestickSeries({
      upColor: '#4ADE80',
      downColor: '#F87171',
      borderDownColor: '#F87171',
      borderUpColor: '#4ADE80',
      wickDownColor: '#F87171',
      wickUpColor: '#4ADE80',
    });

    chartRef.current = chart;
    seriesRef.current = series;

    const handleResize = () => {
      if (containerRef.current) {
        chart.applyOptions({ width: containerRef.current.clientWidth });
      }
    };
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chart.remove();
      chartRef.current = null;
      seriesRef.current = null;
    };
  }, [height]);

  useEffect(() => {
    if (!seriesRef.current || data.length === 0) return;

    const formatted = data
      .map((d) => ({
        time: (new Date(d.time).getTime() / 1000) as Time,
        open: d.open,
        high: d.high,
        low: d.low,
        close: d.close,
      }))
      .sort((a, b) => (a.time as number) - (b.time as number));

    seriesRef.current.setData(formatted);
    chartRef.current?.timeScale().fitContent();
  }, [data]);

  return <div ref={containerRef} />;
}
