import { Skeleton } from "../../ui/skeleton";

export default function ChannelInfoCardSkeleton() {
  return (
    <div className="flex gap-3">
      <Skeleton className="h-[85px] w-[85px] rounded-full" />

      <div className="flex flex-col justify-between">
        <div className="flex flex-col gap-3 mt-1">
          <Skeleton className="h-4 w-[130px]" />
          <Skeleton className="h-4 w-[80px]" />
        </div>
      </div>
    </div>
  );
}
