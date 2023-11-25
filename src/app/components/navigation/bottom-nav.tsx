"use client";

import Link from "next/link";
import clsx from "clsx";
import { usePathname } from "next/navigation";
import { FaListUl } from "react-icons/fa";
import { PiTelevisionFill } from "react-icons/pi";
import { FaRegPlayCircle } from "react-icons/fa";
import { IoMdSettings } from "react-icons/io";
import { FaUser } from "react-icons/fa";

const links = [
  { name: "Channels", href: "/channels", icon: PiTelevisionFill },
  { name: "Playlists", href: "/playlists", icon: FaListUl },
  { name: "Videos", href: "/", icon: FaRegPlayCircle },
  { name: "Settings", href: "/settings", icon: IoMdSettings },
  { name: "Admin", href: "/user", icon: FaUser },
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
                "flex flex-col justify-center items-center text-[#676D75] text-sm",
                { "text-[#ffffff]": link.href === pathname }
              )}
            >
              <LinkIcon className="h-8 w-auto" />
              <p>{link.name}</p>
            </div>
          </Link>
        );
      })}
    </div>
  );
}
