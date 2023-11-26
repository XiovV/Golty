import TopBar from "@/components/navigation/TopBar";
import { FiPlus } from "react-icons/fi";
import { IoSearch } from "react-icons/io5";
import { MdOutlineSort } from "react-icons/md";
import { LuRefreshCw } from "react-icons/lu";
import { IoMdSettings } from "react-icons/io";
import { Channel as IChannel } from "@/types/channel";

import { FiTrash } from "react-icons/fi";
import ChannelCard from "@/components/channel/ChannelCard";
import { Video } from "@/types/video";
import VideosList from "@/components/shared/video/VideosList";

async function fetchChannel(channelName: string): Promise<IChannel> {
  const res = await fetch(`${process.env.API_URL}/channels/${channelName}`, {
    cache: "no-store",
  });

  return await res.json();
}

async function fetchChannelVideos(channelName: string): Promise<Video[]> {
  const res = await fetch(
    `${process.env.VIDEOS_API_URL}/videos/channel/${channelName}`,
    {
      cache: "no-store",
    }
  );

  return res.json();
}

export default async function Page({
  params,
}: {
  params: { channelName: string };
}) {
  const channel = await fetchChannel(params.channelName);
  const channelVideosResponse = fetchChannelVideos(params.channelName);

  return (
    <main>
      <TopBar
        title={params.channelName}
        mobileIcons={[IoSearch, LuRefreshCw, MdOutlineSort, IoMdSettings]}
        desktopIcons={[LuRefreshCw, MdOutlineSort, IoMdSettings, FiTrash]}
      />

      <div className="flex flex-col gap-8 m-5">
        <ChannelCard
          avatarUrl={channel.avatarUrl}
          name={channel.name}
          totalVideos={channel.totalVideos}
          totalSize={channel.totalSize}
        />

        <div className="mx-auto lg:mx-5">
          {/* TODO: display videos list here */}
          <VideosList channelVideosResponse={channelVideosResponse} />
        </div>
      </div>
    </main>
  );
}
