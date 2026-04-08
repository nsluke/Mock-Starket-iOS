'use client';

import { useState, useRef, type FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { TrendingUp, Mail, Phone, Eye, EyeOff, ArrowLeft } from 'lucide-react';
import { useAuthStore } from '@/stores/auth-store';
import { apiClient } from '@/lib/api-client';
import type { ConfirmationResult } from 'firebase/auth';

type AuthView = 'main' | 'email-signin' | 'email-signup' | 'phone';

export default function LoginPage() {
  const router = useRouter();
  const { setToken, setUser } = useAuthStore();
  const [view, setView] = useState<AuthView>('main');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Email form state
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  // Phone form state
  const [phoneNumber, setPhoneNumber] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [confirmResult, setConfirmResult] = useState<ConfirmationResult | null>(null);
  const recaptchaRef = useRef<HTMLDivElement>(null);

  async function finishLogin(token: string) {
    setToken(token);
    localStorage.setItem('mockstarket_token', token);
    document.cookie = `mockstarket_token=${token}; path=/; max-age=${60 * 60 * 24 * 30}`;

    try {
      const user = await apiClient.getMe();
      setUser(user);
    } catch {
      // User just registered — getMe may fail briefly, that's ok
    }

    router.replace('/market');
  }

  async function handleGuest() {
    setIsLoading(true);
    setError(null);
    try {
      const { signInAsGuest } = await import('@/lib/auth-service');
      const token = await signInAsGuest();
      await finishLogin(token);
    } catch (err: any) {
      setError(err.message || 'Guest login failed');
    } finally {
      setIsLoading(false);
    }
  }

  async function handleGoogle() {
    setIsLoading(true);
    setError(null);
    try {
      const { signInWithGoogle } = await import('@/lib/auth-service');
      const token = await signInWithGoogle();
      await finishLogin(token);
    } catch (err: any) {
      if (err.code === 'auth/popup-closed-by-user') return;
      setError(err.message || 'Google sign-in failed');
    } finally {
      setIsLoading(false);
    }
  }

  async function handleApple() {
    setIsLoading(true);
    setError(null);
    try {
      const { signInWithApple } = await import('@/lib/auth-service');
      const token = await signInWithApple();
      await finishLogin(token);
    } catch (err: any) {
      if (err.code === 'auth/popup-closed-by-user') return;
      setError(err.message || 'Apple sign-in failed');
    } finally {
      setIsLoading(false);
    }
  }

  async function handleEmailSignIn(e: FormEvent) {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      const { signInWithEmail } = await import('@/lib/auth-service');
      const token = await signInWithEmail(email, password);
      await finishLogin(token);
    } catch (err: any) {
      const msg = err.code === 'auth/invalid-credential'
        ? 'Invalid email or password'
        : err.code === 'auth/user-not-found'
        ? 'No account with this email'
        : err.message || 'Sign in failed';
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleEmailSignUp(e: FormEvent) {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      const { signUpWithEmail } = await import('@/lib/auth-service');
      const token = await signUpWithEmail(email, password, displayName);
      await finishLogin(token);
    } catch (err: any) {
      const msg = err.code === 'auth/email-already-in-use'
        ? 'An account with this email already exists'
        : err.code === 'auth/weak-password'
        ? 'Password must be at least 6 characters'
        : err.message || 'Sign up failed';
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleSendCode(e: FormEvent) {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      const { sendPhoneVerification } = await import('@/lib/auth-service');
      const result = await sendPhoneVerification(phoneNumber, 'recaptcha-container');
      setConfirmResult(result);
    } catch (err: any) {
      const msg = err.code === 'auth/invalid-phone-number'
        ? 'Enter a valid phone number (e.g. +1234567890)'
        : err.message || 'Failed to send code';
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleVerifyCode(e: FormEvent) {
    e.preventDefault();
    if (!confirmResult) return;
    setIsLoading(true);
    setError(null);
    try {
      const { verifyPhoneCode } = await import('@/lib/auth-service');
      const token = await verifyPhoneCode(confirmResult, verificationCode);
      await finishLogin(token);
    } catch (err: any) {
      setError(err.code === 'auth/invalid-verification-code'
        ? 'Invalid code. Please try again.'
        : err.message || 'Verification failed');
    } finally {
      setIsLoading(false);
    }
  }

  function goBack() {
    setView('main');
    setError(null);
    setConfirmResult(null);
    setVerificationCode('');
  }

  const inputClass = 'w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm';
  const primaryBtn = 'w-full rounded-lg bg-[#50E3C2] px-6 py-3 text-sm font-semibold text-black hover:bg-[#3BC4A7] transition-colors disabled:opacity-50 disabled:cursor-not-allowed';
  const secondaryBtn = 'w-full rounded-lg bg-[#21262D] px-6 py-3 text-sm font-semibold text-white hover:bg-[#30363D] transition-colors disabled:opacity-50';

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-[#0D1117] px-4">
      <div className="w-full max-w-sm space-y-8">
        {/* Logo */}
        <div className="text-center space-y-3">
          <TrendingUp className="w-14 h-14 text-[#50E3C2] mx-auto" />
          <h1 className="text-3xl font-bold tracking-tight text-white">Mock Starket</h1>
          <p className="text-sm text-[#8B949E]">Real stocks. Virtual money. Zero risk.</p>
        </div>

        {/* ---- Main View ---- */}
        {view === 'main' && (
          <div className="space-y-3">
            {/* Social logins */}
            <button onClick={handleGoogle} disabled={isLoading} className={secondaryBtn}>
              <span className="flex items-center justify-center gap-3">
                <svg className="w-5 h-5" viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
                Continue with Google
              </span>
            </button>

            <button onClick={handleApple} disabled={isLoading} className={secondaryBtn}>
              <span className="flex items-center justify-center gap-3">
                <svg className="w-5 h-5 fill-white" viewBox="0 0 24 24"><path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C4.24 16.5 4.89 10.57 8.7 10.3c1.23.07 2.08.72 2.8.75.99-.2 1.94-.77 3-.66 1.28.14 2.24.67 2.86 1.65-2.6 1.56-1.99 5.01.38 5.96-.46 1.24-.97 2.46-1.69 3.28zM12.05 10.23c-.13-2.15 1.67-3.99 3.69-4.17.26 2.43-2.21 4.27-3.69 4.17z"/></svg>
                Continue with Apple
              </span>
            </button>

            {/* Divider */}
            <div className="flex items-center gap-3 py-1">
              <div className="flex-1 h-px bg-[#30363D]" />
              <span className="text-xs text-[#6E7681]">or</span>
              <div className="flex-1 h-px bg-[#30363D]" />
            </div>

            {/* Email */}
            <button onClick={() => { setView('email-signin'); setError(null); }} disabled={isLoading} className={secondaryBtn}>
              <span className="flex items-center justify-center gap-3">
                <Mail className="w-5 h-5" />
                Continue with Email
              </span>
            </button>

            {/* Phone */}
            <button onClick={() => { setView('phone'); setError(null); }} disabled={isLoading} className={secondaryBtn}>
              <span className="flex items-center justify-center gap-3">
                <Phone className="w-5 h-5" />
                Continue with Phone
              </span>
            </button>

            {/* Divider */}
            <div className="flex items-center gap-3 py-1">
              <div className="flex-1 h-px bg-[#30363D]" />
              <span className="text-xs text-[#6E7681]">or</span>
              <div className="flex-1 h-px bg-[#30363D]" />
            </div>

            {/* Guest */}
            <button onClick={handleGuest} disabled={isLoading} className={primaryBtn}>
              {isLoading ? 'Loading...' : 'Try as Guest'}
            </button>
          </div>
        )}

        {/* ---- Email Sign In ---- */}
        {view === 'email-signin' && (
          <div className="space-y-4">
            <button onClick={goBack} className="flex items-center gap-1 text-sm text-[#8B949E] hover:text-white transition-colors">
              <ArrowLeft className="w-4 h-4" /> Back
            </button>
            <h2 className="text-lg font-semibold text-white">Sign in with email</h2>
            <form onSubmit={handleEmailSignIn} className="space-y-3">
              <input
                type="email"
                placeholder="Email address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className={inputClass}
                autoFocus
              />
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={6}
                  className={inputClass}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[#6E7681] hover:text-white"
                >
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              <button type="submit" disabled={isLoading} className={primaryBtn}>
                {isLoading ? 'Signing in...' : 'Sign In'}
              </button>
            </form>
            <p className="text-center text-sm text-[#8B949E]">
              No account?{' '}
              <button onClick={() => { setView('email-signup'); setError(null); }} className="text-[#50E3C2] hover:underline">
                Sign up
              </button>
            </p>
          </div>
        )}

        {/* ---- Email Sign Up ---- */}
        {view === 'email-signup' && (
          <div className="space-y-4">
            <button onClick={goBack} className="flex items-center gap-1 text-sm text-[#8B949E] hover:text-white transition-colors">
              <ArrowLeft className="w-4 h-4" /> Back
            </button>
            <h2 className="text-lg font-semibold text-white">Create an account</h2>
            <form onSubmit={handleEmailSignUp} className="space-y-3">
              <input
                type="text"
                placeholder="Display name"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                required
                maxLength={30}
                className={inputClass}
                autoFocus
              />
              <input
                type="email"
                placeholder="Email address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className={inputClass}
              />
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Password (min 6 characters)"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={6}
                  className={inputClass}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[#6E7681] hover:text-white"
                >
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              <button type="submit" disabled={isLoading} className={primaryBtn}>
                {isLoading ? 'Creating account...' : 'Create Account'}
              </button>
            </form>
            <p className="text-center text-sm text-[#8B949E]">
              Already have an account?{' '}
              <button onClick={() => { setView('email-signin'); setError(null); }} className="text-[#50E3C2] hover:underline">
                Sign in
              </button>
            </p>
          </div>
        )}

        {/* ---- Phone Sign In ---- */}
        {view === 'phone' && (
          <div className="space-y-4">
            <button onClick={goBack} className="flex items-center gap-1 text-sm text-[#8B949E] hover:text-white transition-colors">
              <ArrowLeft className="w-4 h-4" /> Back
            </button>
            <h2 className="text-lg font-semibold text-white">Sign in with phone</h2>

            {!confirmResult ? (
              <form onSubmit={handleSendCode} className="space-y-3">
                <input
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                  value={phoneNumber}
                  onChange={(e) => setPhoneNumber(e.target.value)}
                  required
                  className={inputClass}
                  autoFocus
                />
                <button type="submit" disabled={isLoading} className={primaryBtn}>
                  {isLoading ? 'Sending...' : 'Send Verification Code'}
                </button>
              </form>
            ) : (
              <form onSubmit={handleVerifyCode} className="space-y-3">
                <p className="text-sm text-[#8B949E]">
                  Code sent to <span className="text-white">{phoneNumber}</span>
                </p>
                <input
                  type="text"
                  placeholder="6-digit verification code"
                  value={verificationCode}
                  onChange={(e) => setVerificationCode(e.target.value)}
                  required
                  maxLength={6}
                  pattern="[0-9]{6}"
                  className={inputClass}
                  autoFocus
                />
                <button type="submit" disabled={isLoading} className={primaryBtn}>
                  {isLoading ? 'Verifying...' : 'Verify Code'}
                </button>
                <button
                  type="button"
                  onClick={() => { setConfirmResult(null); setVerificationCode(''); }}
                  className="w-full text-sm text-[#8B949E] hover:text-white transition-colors"
                >
                  Resend code
                </button>
              </form>
            )}
          </div>
        )}

        {/* Error */}
        {error && (
          <p className="text-center text-sm text-red-400 bg-red-400/10 rounded-lg px-4 py-2">{error}</p>
        )}

        {/* Footer */}
        <p className="text-center text-xs text-[#6E7681]">
          Start with $100,000 in virtual cash
        </p>

        {/* Invisible recaptcha container for phone auth */}
        <div id="recaptcha-container" ref={recaptchaRef} />
      </div>
    </div>
  );
}
