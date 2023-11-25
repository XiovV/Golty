import clsx from "clsx";
import Link from "next/link";
import Image from "next/image";

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
        "flex flex-col justify-center items-center text-[#676D75] text-sm font-semibold lg:flex-row lg:justify-normal lg:text-lg lg:py-3 lg:px-4 lg:mb-1 lg:rounded-[16px] lg:gap-3 lg:hover:bg-[#292e37]",
        {
          "text-[#ffffff] lg:bg-[#292e37]": isActive,
        }
      )}
    >
      <LinkIcon className="h-5 w-auto lg:h-6" />
      <p>{name}</p>
    </Link>
  );
}
