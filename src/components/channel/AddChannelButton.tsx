"use client";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useRef } from "react";
import { IconContext } from "react-icons";
import { FiPlus } from "react-icons/fi";
import { Button } from "../ui/button";
import { Checkbox as CheckboxShadcn } from "../ui/checkbox";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Switch as SwitchShadcn } from "../ui/switch";
import ChannelInfoCard from "./ChannelInfoCard";
import ChannelInfoCardSkeleton from "./ChannelInfoCardSkeleton";
import {
  useAddChannel,
  useFetchChannelInfo,
} from "@/hooks/channel/channelHooks";

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
          <DialogTitle>Add a channel</DialogTitle>
        </DialogHeader>

        <AddChannelForm />
      </DialogContent>
    </Dialog>
  );
}

function AddChannelForm() {
  const { loading, channelInfo, getChannelInfo } = useFetchChannelInfo();
  const { addChannel } = useAddChannel();
  const channelUrlRef = useRef<HTMLInputElement>(null);

  return (
    <form onSubmit={addChannel} className="flex flex-col gap-6">
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
    </form>
  );
}

interface SwitchProps {
  label: string;
  name: string;
  defaultChecked?: boolean;
  disabled?: boolean;
}

function Switch({ label, defaultChecked, disabled, name }: SwitchProps) {
  return (
    <div className="flex items-center space-x-2">
      <SwitchShadcn
        id={label}
        name={name}
        disabled={disabled}
        defaultChecked={defaultChecked}
      />
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
        defaultChecked
      />
      <Switch
        label="Download the entire channel"
        name="downloadEntireChannel"
        disabled={disabled}
        defaultChecked
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
      <Checkbox label="Video" name="video" disabled={disabled} defaultChecked />
      <Checkbox label="Audio" name="audio" disabled={disabled} defaultChecked />
    </div>
  );
}

interface CheckboxProps {
  label: string;
  name: string;
  defaultChecked?: boolean;
  disabled?: boolean;
}

function Checkbox({ label, defaultChecked, disabled, name }: CheckboxProps) {
  return (
    <div className="flex items-center space-x-2">
      <CheckboxShadcn
        id={label}
        defaultChecked={defaultChecked}
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
