import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import LoadingState from "../../components/LoadingState";
import ErrorState from "../../components/ErrorState";
import { fetchUserProfile } from "../../api/user";
import type { UserProfile } from "../../types/user";
import ProfileHeader from "./components/ProfileHeader";
import OverallRecordCard from "./components/OverallRecordCard";
import ColorPerformanceCard from "./components/ColorPerformanceCard";
import MatchHistorySection from "./components/MatchHistorySection";

export default function ProfilePage() {
  const { username } = useParams<{ username: string }>();
  return <ProfilePageView key={username} username={username} />;
}

function ProfilePageView({ username }: { username: string | undefined }) {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!username) return;

    fetchUserProfile(username)
      .then((res) => {
        if (res.status === "not-found") {
          setError(`User @${username} not found.`);
        } else {
          setProfile(res.data);
        }
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : "Failed to load user profile");
      })
      .finally(() => setLoading(false));
  }, [username, reloadKey]);

  if (loading) {
    return <LoadingState message="Loading user profile..." />;
  }

  if (error || !profile) {
    return (
      <ErrorState
        message={error || "Failed to load profile."}
        onRetry={() => {
          setLoading(true);
          setError(null);
          setReloadKey((k) => k + 1);
        }}
      />
    );
  }

  const { stats, user, matches } = profile;

  return (
    <div className="max-w-5xl mx-auto w-full px-4 py-6 flex flex-col gap-6">
      <ProfileHeader user={user} stats={stats} />

      <div className="grid md:grid-cols-3 gap-6">
        <OverallRecordCard stats={stats} />
        <ColorPerformanceCard color="white" record={stats.white} />
        <ColorPerformanceCard color="black" record={stats.black} />
      </div>

      <MatchHistorySection matches={matches} />
    </div>
  );
}
