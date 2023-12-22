"use client";

import { usePathname } from "next/navigation";
import NavItem from "./NavItem";
import GoltyLogo from "../svgs/GoltyLogo.svg";
import Image from "next/image";
import { FaListUl } from "react-icons/fa";
import { PiTelevisionFill } from "react-icons/pi";
import { FaRegPlayCircle } from "react-icons/fa";
import { MdVideoLibrary } from "react-icons/md";
import { HiUsers } from "react-icons/hi";
import { MdBugReport } from "react-icons/md";
import { FaBell } from "react-icons/fa";
import Link from "next/link";

const libraryLinks = [
  { name: "Channels", href: "/dashboard/channels", icon: PiTelevisionFill },
  { name: "Playlists", href: "/dashboard/playlists", icon: FaListUl },
  { name: "Videos", href: "/dashboard/videos", icon: FaRegPlayCircle },
];

const settingsLinks = [
  { name: "Libraries", href: "/libraries", icon: MdVideoLibrary },
  { name: "Users", href: "/users", icon: HiUsers },
  { name: "Logs", href: "/logs", icon: MdBugReport },
  { name: "Notifications", href: "/notifications", icon: FaBell },
];

export default function SideNav() {
  const pathname = usePathname();
  return (
    <div className="h-full px-6 py-4 bg-[#1D1F24]">
      <Link href="/" className="flex items-center gap-3 pb-3">
        <Image priority src={GoltyLogo} alt="Golty Logo" />
        <p className="text-white text-2xl font-bold">Golty</p>
      </Link>

      <p className="text-[#676D75] font-semibold pb-3">Library</p>

      {libraryLinks.map((link) => {
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

      <p className="text-[#676D75] font-semibold py-3">Settings</p>

      {settingsLinks.map((link) => {
        return (
          <NavItem
            key={link.name}
            link={link.href}
            name={link.name}
            isActive={link.href === pathname}
            icon={link.icon}
          />
        );
      })}
    </div>
  );
}
