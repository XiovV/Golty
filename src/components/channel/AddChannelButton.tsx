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
import { BaseSyntheticEvent, useEffect, useRef, useState } from "react";
import { Heading1 } from "lucide-react";
import { Channel } from "diagnostics_channel";
import { unescape } from "querystring";
import ChannelCard from "./ChannelCard";
import ChannelInfoCard from "./ChannelInfoCard";
import ChannelInfoCardSkeleton from "./ChannelInfoCardSkeleton";

async function addChannel(e: React.FormEvent<HTMLFormElement>) {
  e.preventDefault();

  const formData = new FormData(e.currentTarget);

  const body = {
    channelUrl: formData.get("channelUrl"),
    downloadVideo: Boolean(formData.get("video")),
    downloadAudio: Boolean(formData.get("audio")),
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

interface ChannelInfo {
  uploader_id: string;
  uploader: string;
  avatar: Avatar;
}

interface Avatar {
  url: string;
}

function AddChannelForm() {
  const [loading, setLoading] = useState(false);
  const [channelInfo, setChannelInfo] = useState<ChannelInfo>();
  const [disableInputs, setDisableInputs] = useState(true);
  const channelUrlRef = useRef<HTMLInputElement>(null);

  async function getChannelInfo(e: BaseSyntheticEvent) {
    const channelUrl = e.target.value;

    if (!channelUrl) {
      return;
    }

    setChannelInfo(undefined);
    setLoading(true);
    const res = await fetch(
      `http://localhost:8080/v1/channels/info/${channelUrl}`,
      { cache: "no-cache" }
    );

    const channelInfo: ChannelInfo = await res.json();
    setLoading(false);
    setChannelInfo(channelInfo);
  }

  return (
    <>
      <div className="grid w-full max-w-sm items-center gap-2">
        <Label htmlFor="channelUrl">Channel URL</Label>
        <Input
          type="channelUrl"
          id="channelUrl"
          name="channelUrl"
          placeholder="Channel URL"
          onBlur={getChannelInfo}
          ref={channelUrlRef}
        />
      </div>

      {loading && <ChannelInfoCardSkeleton />}

      {!loading && channelInfo && (
        <ChannelInfoCard
          avatarUrl={channelInfo.avatar.url}
          name={channelInfo.uploader}
          handle={channelInfo.uploader_id}
          channelUrl={channelUrlRef.current?.value!}
        />
      )}

      <div className="flex flex-col gap-6">
        <div className="flex flex-row justify-between">
          <AddChannelFormCheckboxGroup disabled={!channelInfo} />
          <AddChannelFormSelectGroup disabled={!channelInfo} />
        </div>

        <AddChannelFormSwitchGroup disabled={!channelInfo} />
      </div>

      <DialogFooter>
        <Button type="submit" variant="outline" disabled={!channelInfo}>
          Add Channel
        </Button>
      </DialogFooter>
    </>
  );
}

interface SwitchProps {
  label: string;
  disabled?: boolean;
}

function Switch({ label, disabled }: SwitchProps) {
  return (
    <div className="flex items-center space-x-2">
      <SwitchShadcn id={label} name={label.toLowerCase()} disabled={disabled} />
      <Label htmlFor={label}>{label}</Label>
    </div>
  );
}

interface SwitchGroupProps {
  disabled: boolean;
}

function AddChannelFormSwitchGroup({ disabled }: SwitchGroupProps) {
  return (
    <div className="flex flex-col gap-2">
      <Switch label="Automatically download new uploads" disabled={disabled} />
      <Switch label="Download the entire channel" disabled={disabled} />
    </div>
  );
}

interface CheckboxGroupProps {
  disabled?: boolean;
}

function AddChannelFormCheckboxGroup({ disabled }: CheckboxGroupProps) {
  return (
    <div className="flex flex-col gap-2 ">
      <Checkbox label="Video" disabled={disabled} />
      <Checkbox label="Audio" disabled={disabled} />
    </div>
  );
}

interface CheckboxProps {
  label: string;
  checked?: boolean;
  disabled?: boolean;
}

function Checkbox({ label, checked, disabled }: CheckboxProps) {
  return (
    <div className="flex items-center space-x-2">
      <CheckboxShadcn
        id={label}
        checked={checked}
        name={label.toLowerCase()}
        disabled={disabled}
      />
      <label
        htmlFor={label}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
    </div>
  );
}

interface SelectGroupProps {
  disabled?: boolean;
}

function AddChannelFormSelectGroup({ disabled }: SelectGroupProps) {
  return (
    <div className="flex gap-2">
      <Select name="format" defaultValue="auto" disabled={disabled}>
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

      <Select name="resolution" defaultValue="2160p" disabled={disabled}>
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
