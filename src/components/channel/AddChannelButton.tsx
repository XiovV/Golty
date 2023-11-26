"use client";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { IconContext } from "react-icons";
import { FiPlus } from "react-icons/fi";
import { Label } from "../ui/label";
import { Input } from "../ui/input";

export default function AddChannelButton() {
  return (
    <Dialog>
      <DialogTrigger>
        <IconContext.Provider value={{ size: "1.5em" }}>
          <FiPlus />
        </IconContext.Provider>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <div className="flex flex-col gap-6">
            <DialogTitle>Add a channel</DialogTitle>
            <div className="grid w-full max-w-sm items-center gap-1.5">
              <Label htmlFor="email">Email</Label>
              <Input type="email" id="email" placeholder="Email" />
            </div>
          </div>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
