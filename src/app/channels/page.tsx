import Image from "next/image";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import TopBar from "../components/navigation/top-bar";
import Channel from "../components/channel";

interface Channel {
  avatarUrl: string;
  name: string;
  totalVideos: number;
  totalSize: string;
}

async function fetchChannels(): Promise<Channel[]> {
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
          <div className="flex flex-col gap-6 lg:flex-row lg:flex-wrap lg:gap-x-12">
            {channels.map((channel) => {
              return (
                <Channel
                  avatarUrl={channel.avatarUrl}
                  name={channel.name}
                  totalVideos={channel.totalVideos}
                  totalSize={channel.totalSize}
                />
              );
            })}
          </div>
        </div>
      </div>
    </main>
  );
}
