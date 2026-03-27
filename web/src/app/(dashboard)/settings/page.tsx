'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/stores/auth-store';

interface UserProfile {
  id: string;
  display_name: string;
  avatar_url: string | null;
  is_guest: boolean;
  created_at: string;
  login_streak: number;
  longest_streak: number;
}

export default function SettingsPage() {
  const router = useRouter();
  const { signOut } = useAuthStore();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [displayName, setDisplayName] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  useEffect(() => {
    async function load() {
      try {
        const user = await apiClient.getMe();
        setProfile(user);
        setDisplayName(user.display_name);
      } catch (err) {
        console.error('Failed to load profile:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  async function handleSave() {
    setSaving(true);
    setSaved(false);
    try {
      await apiClient.request('/api/v1/auth/me', {
        method: 'PUT',
        body: JSON.stringify({ display_name: displayName }),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (err) {
      console.error('Failed to update profile:', err);
    } finally {
      setSaving(false);
    }
  }

  function handleSignOut() {
    signOut();
    apiClient.setToken(null);
    router.push('/');
  }

  async function handleDeleteAccount() {
    setDeleting(true);
    try {
      await apiClient.request('/api/v1/auth/me', { method: 'DELETE' });
      signOut();
      apiClient.setToken(null);
      router.push('/');
    } catch (err) {
      console.error('Failed to delete account:', err);
      setDeleting(false);
    }
  }

  if (loading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading...</div>;
  }

  return (
    <div className="p-6 max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      {/* Profile */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
        <h2 className="font-semibold">Profile</h2>

        <div>
          <label className="block text-xs text-[#6E7681] mb-1.5">Display Name</label>
          <input
            type="text"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
          />
        </div>

        <button
          onClick={handleSave}
          disabled={saving || displayName === profile?.display_name}
          className="px-4 py-2 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
        >
          {saving ? 'Saving...' : saved ? 'Saved!' : 'Save Changes'}
        </button>
      </div>

      {/* Account Info */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-3">
        <h2 className="font-semibold">Account</h2>

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-[#6E7681]">Account Type</p>
            <p className="font-medium">{profile?.is_guest ? 'Guest' : 'Registered'}</p>
          </div>
          <div>
            <p className="text-[#6E7681]">Member Since</p>
            <p className="font-medium">
              {profile ? new Date(profile.created_at).toLocaleDateString() : '-'}
            </p>
          </div>
          <div>
            <p className="text-[#6E7681]">Login Streak</p>
            <p className="font-medium">{profile?.login_streak || 0} days</p>
          </div>
          <div>
            <p className="text-[#6E7681]">Longest Streak</p>
            <p className="font-medium">{profile?.longest_streak || 0} days</p>
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
        <h2 className="font-semibold">Actions</h2>

        <button
          onClick={handleSignOut}
          className="w-full py-2.5 rounded-lg border border-[#30363D] text-sm font-medium text-[#8B949E] hover:text-white hover:border-[#50E3C2] transition-colors"
        >
          Sign Out
        </button>

        {/* Danger Zone */}
        <div className="pt-4 border-t border-[#30363D]">
          {!showDeleteConfirm ? (
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="text-sm text-red-400 hover:text-red-300 transition-colors"
            >
              Delete Account
            </button>
          ) : (
            <div className="space-y-3">
              <p className="text-sm text-red-400">
                This will permanently delete your account and all data. This cannot be undone.
              </p>
              <div className="flex gap-3">
                <button
                  onClick={handleDeleteAccount}
                  disabled={deleting}
                  className="px-4 py-2 rounded-lg bg-red-500 text-white text-sm font-semibold hover:bg-red-600 transition-colors disabled:opacity-50"
                >
                  {deleting ? 'Deleting...' : 'Confirm Delete'}
                </button>
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  className="px-4 py-2 rounded-lg border border-[#30363D] text-sm text-[#8B949E] hover:text-white transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
