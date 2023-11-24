"use client";

import {
  TvIcon,
  ListBulletIcon,
  PlayCircleIcon,
  Cog6ToothIcon,
  UserCircleIcon,
} from "@heroicons/react/24/outline";
import Link from "next/link";
import clsx from "clsx";
import { usePathname } from "next/navigation";

const links = [
  { name: "Channels", href: "/channels", icon: TvIcon },
  { name: "Playlists", href: "/playlists", icon: ListBulletIcon },
  { name: "Videos", href: "/", icon: PlayCircleIcon },
  { name: "Settings", href: "/settings", icon: Cog6ToothIcon },
  { name: "Admin", href: "/user", icon: UserCircleIcon },
];

export default function BottomNav() {
  const pathname = usePathname();

  return (
    <div
      className={`fixed flex flex-row bottom-0 w-full justify-evenly bg-[#1D1F24] gap-6 py-3 px-4`}
    >
      {links.map((link) => {
        const LinkIcon = link.icon;

        return (
          <Link key={link.name} href={link.href}>
            <div
              className={clsx(
                "flex flex-col justify-center text-[#676D75] text-sm",
                { "text-[#ffffff]": link.href === pathname }
              )}
            >
              <LinkIcon className="h-8" />
              <p>{link.name}</p>
            </div>
          </Link>
        );
      })}
    </div>
  );
}
