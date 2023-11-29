import Image from "next/image";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { Suspense } from "react";
import TopBar from "@/components/navigation/TopBar";
import ChannelList from "@/components/channel/ChannelList";
import ChannelListSkeleton from "@/components/channel/ChannelCardSkeleton";
import { Channel as IChannel } from "@/types/channel";
import AddChannelButton from "@/components/channel/AddChannelButton";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

async function fetchChannels(): Promise<IChannel[]> {
  const res = await fetch(`${process.env.API_URL}/channels`, {
    cache: "no-store",
  });

  return res.json();
}

export default async function Home() {
  const channels = fetchChannels();

  return (
    <main>
      <TopBar
        title="Channels"
        mobileButtons={[IoSearch, MdOutlineSort, AddChannelButton]}
        desktopButtons={[MdOutlineSort, AddChannelButton]}
      />

      <div className="m-5">
        <h1 className="hidden lg:block text-white text-2xl font-bold">
          Channels
        </h1>

        <div className="mx-3 mt-5 lg:mx-1">
          <Suspense fallback={<ChannelListSkeleton />}>
            <ChannelList channelsResponse={channels} />
          </Suspense>
        </div>
      </div>
    </main>
  );
}
