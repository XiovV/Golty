import Image from "next/image";
import search from "../svgs/search.svg";
import { IoSearch } from "react-icons/io5";

export default function SearchBar() {
  return (
    <div className="flex items-center gap-5 text-">
      <IoSearch className="h-6 w-auto" />
      <input
        type="text"
        placeholder="Search"
        className="bg-transparent text-white"
      />
    </div>
  );
}
