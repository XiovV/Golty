import { API_URL } from "@/app/const";
import { useState, useEffect } from "react";

export interface Channel {
  id: string;
  channelName: string;
  channelHandle: string;
  channelUrl: string;
  avatarUrl: string;
  totalVideos: number;
  totalSize: number;
  state: string;
}

export function fetchChannels() {
  const [channels, setChannels] = useState<Channel[]>();
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_URL}/channels`);
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const result: Channel[] = await response.json();
      setChannels(result);
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  return { channels, loading, fetchData };
}
