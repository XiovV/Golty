import TopBar from "@/components/navigation/TopBar";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { LuRefreshCw } from "react-icons/lu";
import { IoMdSettings } from "react-icons/io";
import { Channel } from "@/types/channel";

import { FiTrash } from "react-icons/fi";
import ChannelCard from "@/components/channel/ChannelCard";
import VideosList from "@/components/shared/video/VideosList";

async function fetchChannel(channelHandle: string): Promise<Channel> {
  const res = await fetch(
    `http://localhost:8080/v1/channels/${channelHandle}`,
    {
      cache: "no-store",
    },
  );

  const channel: Channel = await res.json();

  return channel;
}

export default async function Page({
  params,
}: {
  params: { channelHandle: string };
}) {
  const channel = await fetchChannel(params.channelHandle);

  return (
    <main>
      <TopBar
        title={params.channelHandle}
        mobileButtons={[IoSearch, LuRefreshCw, MdOutlineSort, IoMdSettings]}
        desktopButtons={[LuRefreshCw, MdOutlineSort, IoMdSettings, FiTrash]}
      />

      <div className="flex flex-col gap-8 m-5">
        <ChannelCard
          avatarUrl={channel.avatarUrl}
          channelName={channel.channelName}
          channelHandle={channel.channelHandle}
          totalVideos={channel.totalVideos}
          totalSize={channel.totalSize}
        />

        <div className="mx-auto lg:mx-5 pb-16">
          <VideosList channelHandle={params.channelHandle} />
        </div>
      </div>
    </main>
  );
}
