"use client";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { IconContext } from "react-icons";
import { FiPlus } from "react-icons/fi";
import { Label } from "../ui/label";
import { Input } from "../ui/input";
import { Checkbox as CheckboxShadcn } from "../ui/checkbox";
import { Switch as SwitchShadcn } from "../ui/switch";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Button } from "../ui/button";

async function addChannel(e: React.FormEvent<HTMLFormElement>) {
  e.preventDefault();

  const formData = new FormData(e.currentTarget);

  const body = {
    channelUrl: formData.get("channelUrl"),
    downloadVideo: formData.get("video") == null ? false : true,
    resolution: formData.get("resolution"),
    format: formData.get("format"),
  };

  console.log(body);
}

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
          </div>
        </DialogHeader>

        <form onSubmit={addChannel} className="flex flex-col gap-6">
          <AddChannelForm />
        </form>
      </DialogContent>
    </Dialog>
  );
}

function AddChannelForm() {
  return (
    <>
      <div className="grid w-full max-w-sm items-center gap-2">
        <Label htmlFor="channelUrl">Channel URL</Label>
        <Input
          type="channelUrl"
          id="channelUrl"
          name="channelUrl"
          placeholder="Channel URL"
        />
      </div>

      <div className="flex flex-col gap-6">
        <div className="flex flex-row justify-between">
          <AddChannelFormCheckboxGroup />
          <AddChannelFormSelectGroup />
        </div>

        <AddChannelFormSwitchGroup />
      </div>

      <DialogFooter>
        <Button type="submit" variant="outline">
          Add Channel
        </Button>
      </DialogFooter>
    </>
  );
}

interface SwitchProps {
  label: string;
}

function Switch({ label }: SwitchProps) {
  return (
    <div className="flex items-center space-x-2">
      <SwitchShadcn id={label} name={label.toLowerCase()} />
      <Label htmlFor={label}>{label}</Label>
    </div>
  );
}

function AddChannelFormSwitchGroup() {
  return (
    <div className="flex flex-col gap-2">
      <Switch label="Automatically download new uploads" />
      <Switch label="Download the entire channel" />
    </div>
  );
}

function AddChannelFormCheckboxGroup() {
  return (
    <div className="flex flex-col gap-2 ">
      <Checkbox label="Video" />
      <Checkbox label="Audio" />
    </div>
  );
}

interface AddChannelFormCheckboxProps {
  label: string;
  checked?: boolean;
}

function Checkbox({ label, checked }: AddChannelFormCheckboxProps) {
  return (
    <div className="flex items-center space-x-2">
      <CheckboxShadcn id={label} checked={checked} name={label.toLowerCase()} />
      <label
        htmlFor={label}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
    </div>
  );
}

function AddChannelFormSelectGroup() {
  return (
    <div className="flex gap-2">
      <Select name="format" defaultValue="auto">
        <SelectTrigger className="w-[90px]">
          <SelectValue placeholder="Auto" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel>Format</SelectLabel>
            <SelectItem value="auto">Auto</SelectItem>
            <SelectItem value="m4a">m4a</SelectItem>
            <SelectItem value="mp3">mp3</SelectItem>
            <SelectItem value="mp4">mp4</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>

      <Select name="resolution" defaultValue="2160p">
        <SelectTrigger className="w-[90px]">
          <SelectValue placeholder="2160p" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel>Resolution</SelectLabel>
            <SelectItem value="2160p">2160p</SelectItem>
            <SelectItem value="1440p">1440p</SelectItem>
            <SelectItem value="1080p">1080p</SelectItem>
            <SelectItem value="720p">720p</SelectItem>
            <SelectItem value="480p">480p</SelectItem>
            <SelectItem value="360p">360p</SelectItem>
            <SelectItem value="240p">240p</SelectItem>
            <SelectItem value="144p">144p</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  );
}
