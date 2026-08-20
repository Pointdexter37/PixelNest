export const dynamic = "force-dynamic";

type Wallpaper = {
  id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  width: number;
  height: number;
};

type WallpaperResponse = {
  data: Wallpaper[];
};

async function getWallpapers(): Promise<Wallpaper[]> {
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  const response = await fetch(`${apiURL}/api/v1/wallpapers?limit=12`);

  if (!response.ok) {
    throw new Error("Unable to load wallpapers");
  }

  const payload: WallpaperResponse = await response.json();
  return payload.data;
}

export default async function HomePage() {
  const wallpapers = await getWallpapers();

  return (
    <main>
      <h1>PixNest</h1>
      <p>Discover high-quality wallpapers for your desktop.</p>
      <section className="wallpaper-grid" aria-label="Latest wallpapers">
        {wallpapers.map((wallpaper) => (
          <article className="wallpaper-card" key={wallpaper.id}>
            <img
              src={wallpaper.thumbnail_url}
              alt={wallpaper.title}
              width={wallpaper.width}
              height={wallpaper.height}
            />
            <h2>{wallpaper.title}</h2>
            <p>{wallpaper.description}</p>
          </article>
        ))}
      </section>
    </main>
  );
}
