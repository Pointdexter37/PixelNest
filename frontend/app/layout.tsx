import type { Metadata } from "next";

import "./globals.css";

export const metadata: Metadata = {
  title: "PixNest",
  description: "Discover high-quality wallpapers for your desktop.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}

