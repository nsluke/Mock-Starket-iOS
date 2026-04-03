import { useState, useMemo } from 'react';

type SortDirection = 'asc' | 'desc';

export function useSort<T>(data: T[], defaultKey: keyof T, defaultDir: SortDirection = 'asc') {
  const [sortKey, setSortKey] = useState<keyof T>(defaultKey);
  const [sortDirection, setSortDirection] = useState<SortDirection>(defaultDir);

  const sorted = useMemo(() => {
    return [...data].sort((a, b) => {
      const aVal = a[sortKey];
      const bVal = b[sortKey];

      let cmp: number;
      if (typeof aVal === 'string' && typeof bVal === 'string') {
        // Try numeric comparison for string numbers
        const aNum = parseFloat(aVal);
        const bNum = parseFloat(bVal);
        if (!isNaN(aNum) && !isNaN(bNum)) {
          cmp = aNum - bNum;
        } else {
          cmp = aVal.localeCompare(bVal);
        }
      } else if (typeof aVal === 'number' && typeof bVal === 'number') {
        cmp = aVal - bVal;
      } else {
        cmp = String(aVal).localeCompare(String(bVal));
      }

      return sortDirection === 'asc' ? cmp : -cmp;
    });
  }, [data, sortKey, sortDirection]);

  function onSort(key: keyof T) {
    if (key === sortKey) {
      setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortKey(key);
      setSortDirection('asc');
    }
  }

  return { sorted, sortKey, sortDirection, onSort };
}
