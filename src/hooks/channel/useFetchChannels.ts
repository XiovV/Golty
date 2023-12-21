import { API_URL } from "@/app/const";
import { useState } from "react";
import { Channel } from "./types";

export function useFetchChannels() {
  const [channels, setChannels] = useState<Channel[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    console.log("calling fetch channels");
    try {
      setLoading(true);
      const response = await fetch(`${API_URL}/channels`, {
        cache: "no-cache",
      });
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const result: Channel[] = await response.json();

      setLoading(false);
      setChannels(result);
    } catch (error) {
      console.log(error);
      setLoading(false);
    }
  };

  return { channels, loading, fetchData };
}
