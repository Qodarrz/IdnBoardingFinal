"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Image from "next/image";
import { GetBadges } from "@/helpers/GetBadges";

interface Badge {
  id: number;
  name: string;
  image_url: string;
  description: string;
  is_owned: boolean;
}

export default function BadgesPage() {
  const [badges, setBadges] = useState<Badge[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = typeof window !== "undefined" 
      ? localStorage.getItem("authtoken") 
      : null;

    if (!token) {
      setError("Token not found. Please log in again.");
      setIsLoading(false);
      return;
    }

    const fetchBadges = async () => {
      try {
        const badgesData = await GetBadges(token);
        setBadges(badgesData);
      } catch (error: unknown) {
        if (error instanceof Error) {
          setError(`Gagal mengambil data lencana: ${error.message}`);
        } else {
          setError("Terjadi kesalahan yang tidak diketahui");
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchBadges();
  }, []);

  const earnedBadges = badges.filter((badge) => badge.is_owned);
  const lockedBadges = badges.filter((badge) => !badge.is_owned);

  const renderBadgeGrid = (list: Badge[], isEarned: boolean) => (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
      {list.map((badge) => (
        <Card
          key={badge.id}
          className={`transition ${isEarned ? "opacity-100" : "opacity-50 grayscale"}`}
        >
          <CardHeader className="flex flex-col items-center text-center">
            <div className="w-16 h-16 relative mb-2">
              <Image
                src={badge.image_url}
                alt={badge.name}
                fill
                className="object-contain"
              />
            </div>
            <CardTitle className="text-lg">{badge.name}</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-center text-muted-foreground">
            {badge.description || "No description available"}
          </CardContent>
        </Card>
      ))}
    </div>
  );

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <div className="p-6 space-y-6">
      {/* Lencana yang Dimiliki */}
      <section>
        <h1 className="text-2xl font-bold mb-6">Lencana yang Dimiliki</h1>
        {earnedBadges.length > 0 ? (
          renderBadgeGrid(earnedBadges, true)
        ) : (
          <p className="text-muted-foreground">Belum ada lencana yang dimiliki.</p>
        )}
      </section>

      {/* Lencana yang Belum Dimiliki */}
      <section>
        <h1 className="text-2xl font-bold mb-6">Lencana yang Belum Dimiliki</h1>
        {lockedBadges.length > 0 ? (
          renderBadgeGrid(lockedBadges, false)
        ) : (
          <p className="text-muted-foreground">Semua lencana sudah didapatkan 🎉</p>
        )}
      </section>
    </div>
  );
}
