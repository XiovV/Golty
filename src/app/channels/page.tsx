import Image from "next/image";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { Suspense } from "react";
import TopBar from "@/components/navigation/TopBar";
import ChannelList from "@/components/channel/ChannelList";
import ChannelListSkeleton from "@/components/channel/ChannelCardSkeleton";

export default async function Home() {
  return (
    <main>
      <TopBar
        title="Channels"
        mobileIcons={[IoSearch, MdOutlineSort, FiPlus]}
        desktopIcons={[MdOutlineSort, FiPlus]}
      />

      <div className="m-5">
        <h1 className="hidden lg:block text-white text-2xl font-bold">
          Channels
        </h1>

        <div className="mx-3 mt-5 lg:mx-1">
          <Suspense fallback={<ChannelListSkeleton />}>
            <ChannelList />
          </Suspense>
        </div>
      </div>
    </main>
  );
}
