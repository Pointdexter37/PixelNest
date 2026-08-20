import { notFound } from "next/navigation";

export const dynamic = "force-dynamic";

type Wallpaper = {
  title: string;
  description: string;
  image_url: string;
  width: number;
  height: number;
  file_size: number;
  views: number;
  downloads: number;
};

async function getWallpaper(id: string): Promise<Wallpaper> {
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  const response = await fetch(`${apiURL}/api/v1/wallpapers/${id}`);

  if (response.status === 404) {
    notFound();
  }
  if (!response.ok) {
    throw new Error("Unable to load wallpaper");
  }

  const payload: { data: Wallpaper } = await response.json();
  return payload.data;
}

export default async function WallpaperPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const id = (await params).id;
  const wallpaper = await getWallpaper(id);
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  return (
    <main className="detail-page">
      <a href="/">Back to wallpapers</a>
      <img
        src={wallpaper.image_url}
        alt={wallpaper.title}
        width={wallpaper.width}
        height={wallpaper.height}
      />
      <h1>{wallpaper.title}</h1>
      <p>{wallpaper.description}</p>
      <p>
        {wallpaper.width} x {wallpaper.height} · {wallpaper.file_size} bytes ·{" "}
        {wallpaper.views} views ·{" "}
        {wallpaper.downloads} downloads
      </p>
      <a className="download-link" href={`${apiURL}/api/v1/wallpapers/${id}/download`}>
        Download wallpaper
      </a>
    </main>
  );
}
