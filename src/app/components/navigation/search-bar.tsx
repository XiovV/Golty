import Image from "next/image";
import search from "../svgs/search.svg";

export default function SearchBar() {
  return (
    <div className="flex items-center gap-5">
      <Image priority src={search} alt="" className="h-auto w-6" />
      <input
        type="text"
        placeholder="Search"
        className="bg-transparent text-white"
      />
    </div>
  );
}
