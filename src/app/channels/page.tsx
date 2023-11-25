import TopBar from "../components/navigation/top-bar";
import search from "../components/svgs/search.svg";
import sort from "../components/svgs/sort.svg";
import add from "../components/svgs/add.svg";
import { MdOutlineSort } from "react-icons/md";
import { IoSearch } from "react-icons/io5";
import { FiPlus } from "react-icons/fi";

export default function Page() {
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
