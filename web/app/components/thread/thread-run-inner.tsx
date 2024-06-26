import { Brain, Search } from "lucide-react";
import { Skeleton } from "~/components/ui/skeleton";
import { MarkdownContent } from "~/components/thread/markdown";
import { rewiseModelName } from "~/utils/model";

export function InnerThink({ setting }: { setting?: ThreadSetting }) {
  if (!setting) {
    return <Skeleton className="flex-none w-full h-[45px] rounded-lg" />;
  }

  if (!setting.use_search || !setting.query_list) {
    return null;
  }

  return (
    <div>
      <div className="rounded-lg border py-3 px-3">
        <div className="flex flex-wrap gap-2 text-muted-foreground">
          {setting.model && (
            <div className="text-sm bg-muted py-px px-1 rounded flex items-center gap-0.5 overflow-hidden">
              <Brain className="w-4 h-4 flex-none" />
              <span>{rewiseModelName(setting.model)}</span>
            </div>
          )}

          {setting.query_list?.map((item, index) => (
            <a
              key={index}
              className="text-sm bg-muted py-px px-1 rounded flex items-center gap-0.5 overflow-hidden hover:underline"
              target="_blank"
              href={`https://www.google.com/search?q=${encodeURIComponent(
                item
              )}`}
            >
              <Search className="w-4 h-4 flex-none" />
              <span className="flex-auto truncate">{item}</span>
            </a>
          ))}
        </div>
      </div>
    </div>
  );
}

export function InnerAnswer({ answer }: { answer?: ThreadAnswer }) {
  const isRunning = !answer || answer.status < 2;

  // 正在运行且还没有首包内容
  if (isRunning && !answer?.content) {
    return <Skeleton className="flex-none w-full h-[45px] rounded-lg" />;
  }

  // 有内容，正在运行或者完成了的
  if (answer && answer.status <= 2 && answer?.content) {
    return (
      <div className="bg-muted py-3 px-4 rounded-lg">
        <MarkdownContent content={answer.content} />
      </div>
    );
  }

  // 有错误
  if (answer?.status === 3) {
    // 附加错误
    if (answer.content) {
      return (
        <div className="bg-muted py-3 px-4 rounded-lg">
          <MarkdownContent
            content={
              answer.content +
              `<br /><p class="text-destructive">Network Error...</p>`
            }
          />
        </div>
      );
    }

    return (
      <div className="bg-muted py-3 px-4 rounded-lg">
        <p className="text-destructive">Network Error, please try again...</p>
      </div>
    );
  }

  console.log(answer);
  return (
    <div className="bg-muted py-3 px-4 rounded-lg">
      <p className="text-destructive">Unhandled Error, please try again...</p>
    </div>
  );
}
