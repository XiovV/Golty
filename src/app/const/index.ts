import { DropdownItem } from "@/types/dropdown";

export const API_URL = "http://localhost:8080";

export const videoExtensions: DropdownItem[] = [
  { name: "auto", value: "auto" },
  { name: "m4a", value: "m4a" },
  { name: "mp4", value: "mp4" },
  { name: "webm", value: "webm" },
  { name: "flv", value: "flv" },
  { name: "ogg", value: "ogg" },
  { name: "3gp", value: "3gp" },
];

export const videoResolutions: DropdownItem[] = [
  { name: "2160p", value: "2160p" },
  { name: "1440p", value: "1440p" },
  { name: "1080p", value: "1080p" },
  { name: "720p", value: "720p" },
  { name: "480p", value: "480p" },
  { name: "360p", value: "360p" },
  { name: "240p", value: "240p" },
  { name: "144p", value: "144p" },
];


export const audioExtensions: DropdownItem[] = [
  { name: "auto", value: "auto" },
  { name: "aac", value: "aac" },
  { name: "alac", value: "alac" },
  { name: "flac", value: "flac" },
  { name: "mp3", value: "mp3" },
  { name: "opus", value: "opus" },
  { name: "vorbois", value: "vorbois" },
  { name: "wav", value: "wav" },
];

export const audioQuality: DropdownItem[] = [
  { name: "High (10)", value: "10" },
  { name: "Medium (5)", value: "5" },
  { name: "Low (1)", value: "1" },
]