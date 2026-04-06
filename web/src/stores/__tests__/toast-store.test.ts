import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useToastStore } from '../toast-store';

describe('useToastStore', () => {
  beforeEach(() => {
    useToastStore.setState({ toasts: [] });
    vi.useFakeTimers();
  });

  it('adds a toast', () => {
    useToastStore.getState().addToast('Test message', 'success');
    expect(useToastStore.getState().toasts).toHaveLength(1);
    expect(useToastStore.getState().toasts[0].message).toBe('Test message');
    expect(useToastStore.getState().toasts[0].type).toBe('success');
  });

  it('removes a toast by id', () => {
    useToastStore.getState().addToast('Toast 1', 'info');
    const id = useToastStore.getState().toasts[0].id;
    useToastStore.getState().removeToast(id);
    expect(useToastStore.getState().toasts).toHaveLength(0);
  });

  it('auto-removes toast after duration', () => {
    useToastStore.getState().addToast('Auto remove', 'warning', 2000);
    expect(useToastStore.getState().toasts).toHaveLength(1);

    vi.advanceTimersByTime(2000);
    expect(useToastStore.getState().toasts).toHaveLength(0);
  });

  it('supports multiple toasts', () => {
    useToastStore.getState().addToast('First', 'success');
    useToastStore.getState().addToast('Second', 'error');
    useToastStore.getState().addToast('Third', 'info');
    expect(useToastStore.getState().toasts).toHaveLength(3);
  });
});
