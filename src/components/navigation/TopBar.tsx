// "use client";

import Image from "next/image";
import profile from "../svgs/profile.svg";
import SearchBar from "./SearchBar";
import AddChannelButton from "../channel/AddChannelButton";

interface TopBarProps {
  title: string;
  mobileButtons: any[];
  desktopButtons: any[];
}

export default function TopBar({
  title,
  mobileButtons,
  desktopButtons,
}: TopBarProps) {
  return (
    <div className="lg:bg-[#1D1F24] p-4 top w-full text-[#ffffff]">
      <div className="hidden lg:flex justify-between items-center">
        <SearchBar />

        <div className="flex gap-8">
          {desktopButtons.map((button) => {
            const Button = button;
            return <Button key={button} />;
          })}

          <Image priority src={profile} alt="" className="h-auto w-8" />
        </div>
      </div>

      <div className="flex justify-between items-center">
        <p className="text-2xl text-white font-medium lg:hidden">{title}</p>

        <div className="flex gap-8 lg:hidden">
          {mobileButtons.map((button) => {
            const Button = button;
            return <Button key={button} />;
          })}
        </div>
      </div>
    </div>
  );
}
