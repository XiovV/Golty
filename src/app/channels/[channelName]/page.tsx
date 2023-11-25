import TopBar from "@/app/components/navigation/top-bar";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { LuRefreshCw } from "react-icons/lu";
import { IoMdSettings } from "react-icons/io";
import { Channel as IChannel } from "@/app/types/channel";

import { FiTrash } from "react-icons/fi";
import Channel from "@/app/components/channel";

async function fetchChannel(channelName: string): Promise<IChannel> {
  const res = await fetch(`${process.env.API_URL}/channels/${channelName}`, {
    cache: "no-store",
  });

  return await res.json();
}

export default async function Page({
  params,
}: {
  params: { channelName: string };
}) {
  const channel = await fetchChannel(params.channelName);

  return (
    <main>
      <TopBar
        title={params.channelName}
        mobileIcons={[IoSearch, LuRefreshCw, MdOutlineSort, IoMdSettings]}
        desktopIcons={[LuRefreshCw, MdOutlineSort, IoMdSettings, FiTrash]}
      />

      <div className="m-5">
        <Channel
          avatarUrl={channel.avatarUrl}
          name={channel.name}
          totalVideos={channel.totalVideos}
          totalSize={channel.totalSize}
        />
      </div>
    </main>
  );
}
