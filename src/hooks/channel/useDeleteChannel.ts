import { useRouter } from "next/navigation";
import { useToast } from "@/components/ui/use-toast";
import { ErrorResponse } from "./types";
import { API_URL } from "@/app/const";

export const useDeleteChannel = () => {
  const router = useRouter();
  const { toast } = useToast();

  const deleteChannel = async (channelId: string) => {
    const res = await fetch(`${API_URL}/channels/${channelId}`, {
      method: "DELETE",
    });

    if (res.status != 204) {
      const errResponse: ErrorResponse = await res.json();

      toast({
        title: "Channel could not be deleted!",
        description: errResponse.message,
      });

      return;
    }

    router.push("/dashboard/channels");
  };

  return { deleteChannel };
};
