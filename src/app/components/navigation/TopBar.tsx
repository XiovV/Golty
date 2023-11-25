// "use client";

import Image from "next/image";
import profile from "../svgs/profile.svg";
import SearchBar from "./SearhBar";

interface TopBarProps {
  title: string;
  mobileIcons: any[];
  desktopIcons: any[];
}

export default function TopBar({
  title,
  mobileIcons,
  desktopIcons,
}: TopBarProps) {
  return (
    <div className="lg:bg-[#1D1F24] p-4 top w-full text-[#ffffff]">
      <div className="hidden lg:flex justify-between items-center">
        <SearchBar />

        <div className="flex gap-8">
          {desktopIcons.map((icon) => {
            const LinkIcon = icon;
            return <LinkIcon key={icon} className="h-auto w-6" />;
          })}

          <Image priority src={profile} alt="" className="h-auto w-8" />
        </div>
      </div>

      <div className="flex justify-between items-center">
        <p className="text-2xl text-white font-medium lg:hidden">{title}</p>

        <div className="flex gap-8 lg:hidden">
          {mobileIcons.map((icon) => {
            const LinkIcon = icon;
            return <LinkIcon key={icon} className="h-auto w-6" />;
          })}
        </div>
      </div>
    </div>
  );
}
