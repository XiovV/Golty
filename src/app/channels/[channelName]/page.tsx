import TopBar from "@/app/components/navigation/top-bar";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { LuRefreshCw } from "react-icons/lu";
import { IoMdSettings } from "react-icons/io";

import { FiTrash } from "react-icons/fi";

export default function Page({ params }: { params: { channelName: string } }) {
  return (
    <main>
      <TopBar
        title={params.channelName}
        mobileIcons={[IoSearch, LuRefreshCw, MdOutlineSort, IoMdSettings]}
        desktopIcons={[LuRefreshCw, MdOutlineSort, IoMdSettings, FiTrash]}
      />
    </main>
  );
}
