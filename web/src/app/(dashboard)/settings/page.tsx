'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { PageTransition } from '@/components/ui/PageTransition';
import { useCurrentUser, useUpdateProfile } from '@/hooks/use-user';
import { useAuthStore } from '@/stores/auth-store';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';
import { profileSchema, type ProfileFormValues } from '@/lib/schemas';
import { FormInput } from '@/components/ui/FormInput';

export default function SettingsPage() {
  const router = useRouter();
  const { signOut } = useAuthStore();
  const { data: profile, isLoading } = useCurrentUser();
  const updateProfile = useUpdateProfile();
  const { register, handleSubmit, reset, formState: { errors, isDirty } } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
  });

  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (profile) {
      reset({ display_name: profile.display_name });
    }
  }, [profile, reset]);

  function onSave(data: ProfileFormValues) {
    updateProfile.mutate(data.display_name);
  }

  function handleSignOut() {
    signOut();
    apiClient.setToken(null);
    document.cookie = 'mockstarket_token=; path=/; max-age=0';
    router.push('/');
  }

  async function handleDeleteAccount() {
    setDeleting(true);
    try {
      await apiClient.request('/api/v1/auth/me', { method: 'DELETE' });
      signOut();
      apiClient.setToken(null);
      document.cookie = 'mockstarket_token=; path=/; max-age=0';
      router.push('/');
    } catch {
      toast.error('Failed to delete account');
      setDeleting(false);
    }
  }

  if (isLoading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading...</div>;
  }

  return (
    <PageTransition>
    <div className="p-6 max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      {/* Profile */}
      <form onSubmit={handleSubmit(onSave)} className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
        <h2 className="font-semibold">Profile</h2>

        <FormInput
          label="Display Name"
          error={errors.display_name?.message}
          {...register('display_name')}
        />

        <button
          type="submit"
          disabled={updateProfile.isPending || !isDirty}
          className="px-4 py-2 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
        >
          {updateProfile.isPending ? 'Saving...' : 'Save Changes'}
        </button>
      </form>

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
    </PageTransition>
  );
}
