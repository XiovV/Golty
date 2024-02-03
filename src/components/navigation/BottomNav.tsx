"use client";

import { usePathname } from "next/navigation";
import { FaListUl } from "react-icons/fa";
import { PiTelevisionFill } from "react-icons/pi";
import { FaRegPlayCircle } from "react-icons/fa";
import { IoMdSettings } from "react-icons/io";
import { FaUser } from "react-icons/fa";
import NavItem from "./NavItem";

const links = [
  { name: "Channels", href: "/dashboard/channels", icon: PiTelevisionFill },
  { name: "Playlists", href: "/dashboard/playlists", icon: FaListUl },
  { name: "Videos", href: "/dashboard/videos", icon: FaRegPlayCircle },
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
        const isActive = pathname.includes(link.href);
        return (
          <NavItem
            key={link.name}
            link={link.href}
            name={link.name}
            isActive={isActive}
            icon={link.icon}
          />
        );
      })}
    </div>
  );
}
