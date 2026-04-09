import {
  signInAnonymously,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signInWithPopup,
  signInWithPhoneNumber,
  GoogleAuthProvider,
  OAuthProvider,
  RecaptchaVerifier,
  updateProfile,
  signOut as firebaseSignOut,
  type User,
  type ConfirmationResult,
} from 'firebase/auth';
import { auth } from './firebase';
import { apiClient } from './api-client';

// Providers
const googleProvider = new GoogleAuthProvider();
const appleProvider = new OAuthProvider('apple.com');
appleProvider.addScope('email');
appleProvider.addScope('name');

/**
 * After any Firebase sign-in, get the ID token and register with our backend.
 * Returns the Firebase ID token for use as the API bearer token.
 */
async function finalizeAuth(user: User, displayName?: string): Promise<string> {
  const token = await user.getIdToken();
  apiClient.setToken(token);

  // Register or update user in our backend
  const name = displayName || user.displayName || user.email?.split('@')[0] || 'Trader';
  const isGuest = user.isAnonymous;

  try {
    await apiClient.register(name, isGuest);
  } catch (e: any) {
    // User may already exist (409 conflict) — that's fine
    if (!e.message?.includes('409')) {
      throw e;
    }
  }

  return token;
}

/** Guest login via Firebase Anonymous Auth */
export async function signInAsGuest(): Promise<string> {
  const result = await signInAnonymously(auth);
  return finalizeAuth(result.user, 'Guest Trader');
}

/** Email + password sign in */
export async function signInWithEmail(email: string, password: string): Promise<string> {
  const result = await signInWithEmailAndPassword(auth, email, password);
  return finalizeAuth(result.user);
}

/** Email + password sign up */
export async function signUpWithEmail(email: string, password: string, displayName: string): Promise<string> {
  const result = await createUserWithEmailAndPassword(auth, email, password);
  if (displayName) {
    await updateProfile(result.user, { displayName });
  }
  return finalizeAuth(result.user, displayName);
}

/** Google sign in via popup */
export async function signInWithGoogle(): Promise<string> {
  const result = await signInWithPopup(auth, googleProvider);
  return finalizeAuth(result.user);
}

/** Apple sign in via popup */
export async function signInWithApple(): Promise<string> {
  const result = await signInWithPopup(auth, appleProvider);
  return finalizeAuth(result.user);
}

/** Phone number sign in — Step 1: send verification code */
let recaptchaVerifier: RecaptchaVerifier | null = null;

export async function sendPhoneVerification(
  phoneNumber: string,
  recaptchaContainerId: string
): Promise<ConfirmationResult> {
  // Clear previous verifier to avoid "already rendered" error
  if (recaptchaVerifier) {
    recaptchaVerifier.clear();
    recaptchaVerifier = null;
  }
  recaptchaVerifier = new RecaptchaVerifier(auth, recaptchaContainerId, {
    size: 'invisible',
  });
  return signInWithPhoneNumber(auth, phoneNumber, recaptchaVerifier);
}

/** Phone number sign in — Step 2: verify the code */
export async function verifyPhoneCode(
  confirmationResult: ConfirmationResult,
  code: string
): Promise<string> {
  const result = await confirmationResult.confirm(code);
  return finalizeAuth(result.user, result.user.phoneNumber || undefined);
}

/** Sign out of Firebase and clear API token */
export async function signOut(): Promise<void> {
  await firebaseSignOut(auth);
  apiClient.setToken(null);
}

/** Refresh the Firebase ID token (tokens expire after 1 hour) */
export async function refreshToken(): Promise<string | null> {
  const user = auth.currentUser;
  if (!user) return null;
  const token = await user.getIdToken(true);
  apiClient.setToken(token);
  return token;
}
