import clsx from "clsx";
import Link from "next/link";
import { ForwardRefExoticComponent } from "react";

interface NavItemProps {
  link: string;
  name: string;
  isActive: boolean;
  icon: any;
}

export default function NavItem({ link, name, isActive, icon }: NavItemProps) {
  const LinkIcon = icon;
  return (
    <Link
      href={link}
      className={clsx(
        "flex text-[#676D75] font-semibold text-lg py-3 px-4 rounded-[16px] items-center gap-3 hover:bg-[#292E37]",
        {
          "text-[#ffffff] font-bold bg-[#292E37]": isActive,
        }
      )}
    >
      <LinkIcon className="h-7 " />
      <p>{name}</p>
    </Link>
  );
}
