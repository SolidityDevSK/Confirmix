import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Theme } from "@radix-ui/themes";
import "@radix-ui/themes/styles.css";
import { BlockchainProvider } from "@/contexts/BlockchainContext";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Confirmix - Blockchain Explorer",
  description: "A modern blockchain explorer for Confirmix PoA blockchain",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Theme appearance="dark" accentColor="blue" radius="medium">
          <BlockchainProvider>
            {children}
          </BlockchainProvider>
        </Theme>
      </body>
    </html>
  );
}
