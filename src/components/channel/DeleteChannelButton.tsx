"use client";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { IconContext } from "react-icons";
import { LuTrash } from "react-icons/lu";

interface DeleteChannelButtonProps {
  channelId: string;
}

export default function DeleteChannelButton({
  channelId,
}: DeleteChannelButtonProps) {
  return (
    <AlertDialog>
      <AlertDialogTrigger>
        <IconContext.Provider value={{ size: "1.5em" }}>
          <LuTrash />
        </IconContext.Provider>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Are you sure you want to delete this channel?
          </AlertDialogTitle>

          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete the
            channel and all the content related to it.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={() => console.log(channelId)}>
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
