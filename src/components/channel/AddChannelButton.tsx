"use client";
import {
  Dialog,
  DialogContent,
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
import { BaseSyntheticEvent, useRef, useState } from "react";
import ChannelInfoCard from "./ChannelInfoCard";
import ChannelInfoCardSkeleton from "./ChannelInfoCardSkeleton";
import { useDebouncedCallback } from "use-debounce";
import { useToast } from "../ui/use-toast";

interface ErrorResponse {
  message: string;
}

export default function AddChannelButton() {
  const { toast } = useToast();

  async function addChannel(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    const body = {
      channelUrl: formData.get("channelUrl"),
      downloadVideo: Boolean(formData.get("video")),
      downloadAudio: Boolean(formData.get("audio")),
      resolution: formData.get("resolution"),
      format: formData.get("format"),
      downloadAutomatically: Boolean(formData.get("downloadAutomatically")),
      downloadEntireChannel: Boolean(formData.get("downloadEntireChannel")),
    };

    const res = await fetch("http://localhost:8080/v1/channels", {
      method: "POST",
      body: JSON.stringify(body),
      headers: { "Content-Type": "application/json" },
    });

    if (res.status !== 201) {
      const err: ErrorResponse = await res.json();

      toast({
        title: "Unable to add the channel!",
        description: err.message,
      });

      return;
    }

    toast({
      title: "Channel added successfully!",
    });
  }

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
  const channelUrlRef = useRef<HTMLInputElement>(null);

  const getChannelInfo = useDebouncedCallback(async (channelUrl: string) => {
    if (!channelUrl || !channelUrl.includes("https://www.youtube.com/")) {
      return;
    }

    setChannelInfo(undefined);
    setLoading(true);
    const res = await fetch(
      `http://localhost:8080/v1/channels/info/${channelUrl}`,
      { cache: "no-cache" }
    );

    if (res.status !== 200) {
      setLoading(false);
      return;
    }

    const channelInfo: ChannelInfo = await res.json();
    setLoading(false);
    setChannelInfo(channelInfo);
  }, 500);

  return (
    <>
      <div className="grid w-full max-w-sm items-center gap-2">
        <Label htmlFor="channelUrl">Channel URL</Label>
        <Input
          type="channelUrl"
          id="channelUrl"
          name="channelUrl"
          placeholder="Channel URL"
          onChange={(e) => getChannelInfo(e.target.value)}
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
          {!channelInfo && "Add Channel"}

          {!loading && channelInfo && `Add ${channelInfo.uploader}`}
        </Button>
      </DialogFooter>
    </>
  );
}

interface SwitchProps {
  label: string;
  name: string;
  disabled?: boolean;
}

function Switch({ label, disabled, name }: SwitchProps) {
  return (
    <div className="flex items-center space-x-2">
      <SwitchShadcn id={label} name={name} disabled={disabled} />
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
      <Switch
        label="Automatically download new uploads"
        name="downloadAutomatically"
        disabled={disabled}
      />
      <Switch
        label="Download the entire channel"
        name="downloadEntireChannel"
        disabled={disabled}
      />
    </div>
  );
}

interface CheckboxGroupProps {
  disabled?: boolean;
}

function AddChannelFormCheckboxGroup({ disabled }: CheckboxGroupProps) {
  return (
    <div className="flex flex-col gap-2 ">
      <Checkbox label="Video" name="video" disabled={disabled} />
      <Checkbox label="Audio" name="audio" disabled={disabled} />
    </div>
  );
}

interface CheckboxProps {
  label: string;
  name: string;
  checked?: boolean;
  disabled?: boolean;
}

function Checkbox({ label, checked, disabled, name }: CheckboxProps) {
  return (
    <div className="flex items-center space-x-2">
      <CheckboxShadcn
        id={label}
        checked={checked}
        name={name}
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
