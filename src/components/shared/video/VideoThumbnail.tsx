import Image from "next/image";

interface VideoThumbnailProps {
  thumbnailUrl: string;
  width?: number;
  height?: number;
}

export default function VideoThumbnail({
  thumbnailUrl,
  width,
  height,
}: VideoThumbnailProps) {
  return (
    <Image
      priority
      src={thumbnailUrl}
      height={height ? height : 0}
      width={width ? width : 350}
      alt={"video thumbnail"}
    />
  );
}
