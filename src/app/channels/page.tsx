import Image from "next/image";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import TopBar from "../components/navigation/top-bar";
import Channel from "../components/channel";
import { Channel as IChannel } from "../types/channel";
import ChannelList from "../components/channel-list";

async function fetchChannels(): Promise<IChannel[]> {
  const res = await fetch(`${process.env.API_URL}/channels`, {
    cache: "no-store",
  });

  return await res.json();
}

export default async function Home() {
  const channels = await fetchChannels();

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
          <ChannelList channels={channels} />
        </div>
      </div>
    </main>
  );
}
