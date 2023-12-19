// "use client";

import Image from "next/image";
import profile from "../svgs/profile.svg";
import SearchBar from "./SearchBar";

interface TopBarProps {
  title: string;
  children: React.ReactNode;
}

export function TopBar({ children }: TopBarProps) {
  return (
    <div className="lg:bg-[#1D1F24] p-4 top w-full text-[#ffffff]">
      {children}
    </div>
  );
}

export function DesktopButtons({ children }: { children: React.ReactNode }) {
  return (
    <div className="hidden lg:flex justify-between">
      <SearchBar />

      <div className="flex gap-8 items-center">
        {children}
        <Image priority src={profile} alt="" className="h-auto w-8" />
      </div>
    </div>
  );
}

interface MobileButtonsProps {
  title: string;
  children: React.ReactNode;
}

export function MobileButtons({ title, children }: MobileButtonsProps) {
  return (
    <div className="flex justify-between items-center">
      <p className="text-2xl text-white font-medium lg:hidden">{title}</p>

      <div className="flex gap-8 lg:hidden">{children}</div>
    </div>
  );
}
