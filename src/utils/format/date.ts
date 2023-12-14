const secondsInAYear = 31536000;
const secondsInAMonth = 2629746;
const secondsInAWeek = 604800;
const secondsInADay = 86400;
const secondsInAnHour = 3600;

export function formatTimeAgo(timestamp: number) {
  const seconds = Math.floor((Date.now() - timestamp * 1000) / 1000);

  const years = Math.floor(seconds / secondsInAYear);
  if (years) {
    return formatTimeAgoString(years, "year");
  }

  const months = Math.floor(seconds / secondsInAMonth);
  if (months) {
    return formatTimeAgoString(months, "month");
  }

  const weeks = Math.floor(seconds / secondsInAWeek);
  if (weeks) {
    return formatTimeAgoString(weeks, "week");
  }

  const days = Math.floor(seconds / secondsInADay);
  if (days) {
    return formatTimeAgoString(days, "day");
  }

  const hours = Math.floor(seconds / secondsInAnHour);
  if (hours) {
    return formatTimeAgoString(hours, "hour");
  }

  const minutes = Math.floor(seconds / 60);
  if (minutes === 0) {
    return `Just now`;
  }

  if (minutes) {
    return formatTimeAgoString(minutes, "minute");
  }
}

function formatTimeAgoString(amount: number, timeUnit: string) {
  if (amount > 1) {
    return `${amount} ${timeUnit}s ago`;
  }

  return `${amount} ${timeUnit} ago`;
}
