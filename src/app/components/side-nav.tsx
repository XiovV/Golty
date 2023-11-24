"use client";

import {
  ListBulletIcon,
  PlayCircleIcon,
  TvIcon,
  DocumentDuplicateIcon,
  UsersIcon,
  BugAntIcon,
  BellAlertIcon,
} from "@heroicons/react/24/outline";
import clsx from "clsx";
import Link from "next/link";
import { usePathname } from "next/navigation";

const libraryLinks = [
  { name: "Channels", href: "/channels", icon: TvIcon },
  { name: "Playlists", href: "/playlists", icon: ListBulletIcon },
  { name: "Videos", href: "/", icon: PlayCircleIcon },
];

const settingsLinks = [
  { name: "Libraries", href: "/libraries", icon: DocumentDuplicateIcon },
  { name: "Users", href: "/users", icon: UsersIcon },
  { name: "Logs", href: "/logs", icon: BugAntIcon },
  { name: "Notifications", href: "/notifications", icon: BellAlertIcon },
];

export default function SideNav() {
  const pathname = usePathname();
  return (
    <div className="h-full px-6 py-4 bg-[#1D1F24]">
      <p className="text-[#676D75] font-semibold pb-3">Library</p>

      {libraryLinks.map((link) => {
        const LinkIcon = link.icon;

        return (
          <Link
            href={link.href}
            className={clsx(
              "flex text-[#676D75] text-lg font-semibold py-3 px-4 rounded-[16px] items-center gap-3 hover:bg-[#292E37]",
              {
                "text-[#ffffff] font-bold bg-[#292E37]": link.href === pathname,
              }
            )}
          >
            <LinkIcon className="h-7 " />
            <p>{link.name}</p>
          </Link>
        );
      })}

      <p className="text-[#676D75] font-semibold py-3">Settings</p>
      {settingsLinks.map((link) => {
        const LinkIcon = link.icon;

        return (
          <Link
            href={link.href}
            className={clsx(
              "flex text-[#676D75] font-semibold text-lg py-3 px-4 rounded-[16px] items-center gap-3 hover:bg-[#292E37]",
              {
                "text-[#ffffff] font-bold bg-[#292E37]": link.href === pathname,
              }
            )}
          >
            <LinkIcon className="h-7 " />
            <p>{link.name}</p>
          </Link>
        );
      })}
    </div>
  );
}
