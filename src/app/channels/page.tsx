import TopBar from "../components/navigation/top-bar";
import search from "../components/svgs/search.svg";
import sort from "../components/svgs/sort.svg";
import add from "../components/svgs/add.svg";

export default function Page() {
  return (
    <main>
      <TopBar
        title="Channels"
        mobileIcons={[search, sort, add]}
        desktopIcons={[sort, add]}
      />
      <h1 className={`text-xl`}>Channels page</h1>
    </main>
  );
}
