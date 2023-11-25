import Image from "next/image";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import TopBar from "../components/navigation/top-bar";

export default function Home() {
  return (
    <main>
      <TopBar
        title="Channels"
        mobileIcons={[IoSearch, MdOutlineSort, FiPlus]}
        desktopIcons={[MdOutlineSort, FiPlus]}
      />
      <h1 className={`text-xl`}>Channels page</h1>
    </main>
  );
}
