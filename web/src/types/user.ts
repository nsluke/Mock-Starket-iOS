export interface User {
  id: string;
  firebase_uid: string;
  display_name: string;
  avatar_url?: string;
  is_guest: boolean;
  created_at: string;
  login_streak: number;
  longest_streak: number;
}

export interface LeaderboardEntry {
  id: number;
  user_id: string;
  display_name: string;
  net_worth: string;
  total_return: string;
  rank: number;
  period: string;
}

export interface Achievement {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
}
