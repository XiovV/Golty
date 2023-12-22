import {
  TopBar,
  DesktopButtons,
  MobileButtons,
} from "@/components/navigation/TopBar";
import AddChannelButton from "@/components/channel/buttons/AddChannelButton";
import ChannelList from "@/components/channel/ChannelList";

export default async function Home() {
  return (
    <main>
      <TopBar>
        <DesktopButtons>
          <AddChannelButton />
        </DesktopButtons>

        <MobileButtons title="Channels">
          <AddChannelButton />
        </MobileButtons>
      </TopBar>

      <div className="m-5">
        <h1 className="hidden lg:block text-white text-2xl font-bold">
          Channels
        </h1>

        <div className="mx-3 mt-5 lg:mx-1">
          <ChannelList />
          {/* <Suspense fallback={<ChannelListSkeleton />}>
            <ChannelList channelsResponse={channels} />
          </Suspense> */}
        </div>
      </div>
    </main>
  );
}
