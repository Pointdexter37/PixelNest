import Link from "next/link";

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
  meta: { page: number; limit: number; total: number };
};

type Category = { id: number; name: string };

async function getWallpapers(
  query: string,
  categoryID: string,
  page: number,
): Promise<WallpaperResponse> {
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  const search = query ? `&q=${encodeURIComponent(query)}` : "";
  const category = categoryID ? `&category_id=${categoryID}` : "";
  const response = await fetch(
    `${apiURL}/api/v1/wallpapers?page=${page}&limit=12${search}${category}`,
  );

  if (!response.ok) {
    throw new Error("Unable to load wallpapers");
  }

  const payload: WallpaperResponse = await response.json();
  return payload;
}

async function getCategories(): Promise<Category[]> {
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  const response = await fetch(`${apiURL}/api/v1/categories`);
  if (!response.ok) throw new Error("Unable to load categories");
  const payload: { data: Category[] } = await response.json();
  return payload.data;
}

export default async function HomePage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string; category_id?: string; page?: string }>;
}) {
  const params = await searchParams;
  const query = params.q ?? "";
  const categoryID = params.category_id ?? "";
  const page = Math.max(1, Number(params.page ?? "1") || 1);
  const [wallpaperResponse, categories] = await Promise.all([
    getWallpapers(query, categoryID, page),
    getCategories(),
  ]);
  const totalPages = Math.ceil(wallpaperResponse.meta.total / wallpaperResponse.meta.limit);
  const pageLink = (nextPage: number) => {
    const search = new URLSearchParams();
    if (query) search.set("q", query);
    if (categoryID) search.set("category_id", categoryID);
    search.set("page", String(nextPage));
    return `/?${search.toString()}`;
  };

  return (
    <main>
      <h1>PixelNest</h1>
      <p>Discover high-quality wallpapers for your desktop.</p>
      <form className="search-form" role="search">
        <input
          name="q"
          type="search"
          placeholder="Search wallpapers"
          defaultValue={query}
          aria-label="Search wallpapers"
        />
        <button type="submit">Search</button>
        <select name="category_id" defaultValue={categoryID} aria-label="Filter by category">
          <option value="">All categories</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </form>
      {query && <p>Search results for &quot;{query}&quot;</p>}
      <section className="wallpaper-grid" aria-label="Latest wallpapers">
        {wallpaperResponse.data.map((wallpaper) => (
          <article className="wallpaper-card" key={wallpaper.id}>
            <img
              src={wallpaper.thumbnail_url}
              alt={wallpaper.title}
              width={wallpaper.width}
              height={wallpaper.height}
            />
            <h2>{wallpaper.title}</h2>
            <p>{wallpaper.description}</p>
            <Link href={`/wallpapers/${wallpaper.id}`}>View wallpaper</Link>
          </article>
        ))}
      </section>
      {totalPages > 1 && (
        <nav className="pagination" aria-label="Pagination">
          {page > 1 && <Link href={pageLink(page - 1)}>Previous</Link>}
          <span>
            Page {page} of {totalPages}
          </span>
          {page < totalPages && <Link href={pageLink(page + 1)}>Next</Link>}
        </nav>
      )}
    </main>
  );
}
