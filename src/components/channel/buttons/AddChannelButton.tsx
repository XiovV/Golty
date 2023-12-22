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
import { Button } from "../../ui/button";
import { Checkbox as CheckboxShadcn } from "../../ui/checkbox";
import { Input } from "../../ui/input";
import { Label } from "../../ui/label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../../ui/select";
import { Switch as SwitchShadcn } from "../../ui/switch";
import ChannelInfoCard from "../cards/ChannelInfoCard";
import ChannelInfoCardSkeleton from "../cards/ChannelInfoCardSkeleton";
import { useAddChannel, useFetchChannelInfo } from "@/hooks/channel";
import { useState } from "react";
import { audioExtensions, resolutions, videoExtensions } from "@/app/const";

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
  const [downloadExtensions, setDownloadExtensions] =
    useState<string[]>(videoExtensions);
  const { loading, channelInfo, getChannelInfo } = useFetchChannelInfo();
  const { addChannel } = useAddChannel();

  return (
    <form
      onSubmit={(e) => addChannel(e, channelInfo!)}
      className="flex flex-col gap-6"
    >
      <div className="grid w-full max-w-sm items-center gap-2">
        <Label htmlFor="channelInput">Channel</Label>
        <Input
          type="text"
          id="channelInput"
          name="channelInput"
          placeholder="URL or Handle"
          onChange={(e) => getChannelInfo(e.target.value)}
        />
      </div>

      {loading && <ChannelInfoCardSkeleton />}

      {!loading && channelInfo && (
        <ChannelInfoCard
          avatarUrl={channelInfo.avatarUrl}
          name={channelInfo.uploader}
          handle={channelInfo.uploaderId}
          channelUrl={channelInfo.uploaderUrl}
        />
      )}

      <div className="flex flex-col gap-6">
        <div className="flex flex-row justify-between">
          <div className="flex flex-col gap-2 ">
            <div className="flex items-center space-x-2">
              <CheckboxShadcn
                id="video"
                defaultChecked
                name="video"
                disabled={!channelInfo}
                onCheckedChange={(checked) => {
                  const extensions = checked
                    ? videoExtensions
                    : audioExtensions;
                  setDownloadExtensions(extensions);
                }}
              />
              <label
                htmlFor={"video"}
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Video
              </label>
            </div>

            <div className="flex items-center space-x-2">
              <CheckboxShadcn
                id="audio"
                defaultChecked
                name="audio"
                disabled={!channelInfo}
              />
              <label
                htmlFor="audio"
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Audio
              </label>
            </div>
          </div>

          <div className="flex gap-2">
            <Dropdown
              items={downloadExtensions}
              name="format"
              label="Format"
              defaultValue="Auto"
              disabled={!channelInfo}
            />

            <Dropdown
              items={resolutions}
              name="resolution"
              label="Resolution"
              defaultValue="2160p"
              disabled={!channelInfo}
            />
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <Switch
            label="Automatically download new uploads"
            name="downloadAutomatically"
            disabled={!channelInfo}
            defaultChecked
          />
          <Switch
            label="Download the entire channel"
            name="downloadEntireChannel"
            disabled={!channelInfo}
            defaultChecked
          />
        </div>
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

interface DropdownProps {
  items: string[];
  name: string;
  label: string;
  defaultValue: string;
  disabled?: boolean;
}

function Dropdown({
  items,
  name,
  defaultValue,
  label,
  disabled,
}: DropdownProps) {
  return (
    <Select name={name} defaultValue={defaultValue} disabled={disabled}>
      <SelectTrigger className="w-[90px]">
        <SelectValue placeholder="Auto" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>{label}</SelectLabel>

          {items.map((item) => (
            <SelectItem key={item} value={item}>
              {item}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
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
