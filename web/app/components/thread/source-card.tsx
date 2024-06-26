import { SquareArrowOutUpRight } from "lucide-react";

interface Props {
  item: SearchItem;
}

export const SourceCard = (props: Props) => {
  const { item } = props;

  const domain = getDomainFromUrl(item.link);

  return (
    <div className="flex-none w-[200px] overflow-hidden p-3 rounded-lg border text-left text-sm transition-all hover:bg-accent">
      <a href={item.link} target="_blank" className="group">
        <div className="flex items-center gap-1 truncate mb-1">
          <img
            src={`https://s2.googleusercontent.com/s2/favicons?domain=${domain}&sz=16`}
            alt={domain}
            className="rounded"
          />
          <span className="text-xs text-muted-foreground mr-1 truncate">
            {domain}
          </span>
          <SquareArrowOutUpRight className="ml-auto flex-none w-3 h-3 text-muted-foreground" />
        </div>

        <div className="mb-1 text-sm font-semibold tracking-tight truncate group-hover:underline">
          {item.title}
        </div>
      </a>

      <p
        className="line-clamp-2 text-xs text-muted-foreground"
        title={item.snippet}
      >
        {item.snippet}
      </p>
    </div>
  );
};

function getDomainFromUrl(url: string) {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname; // 获取域名
  } catch (e) {
    console.error("invalid url:", e);
    return "";
  }
}
