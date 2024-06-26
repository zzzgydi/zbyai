import clsx from "clsx";
import { useEffect, useMemo, useRef, useState } from "react";
import { produce } from "immer";
import { ChevronLeft, ChevronRight, RefreshCw } from "lucide-react";
import { useThrottleFn } from "ahooks";
import { Button } from "~/components/ui/button";
import { Skeleton } from "~/components/ui/skeleton";
import { SourceCard } from "~/components/thread/source-card";
import { ApiClient } from "~/services/api.client";
import { rewiseModelName } from "~/utils/model";
import { InnerAnswer, InnerThink } from "./thread-run-inner";

interface Props {
  loading: boolean;
  threadId: string;
  item: ThreadRun;
  onRefresh?: () => void;
  onStartLoading?: () => void;
  onDone?: () => void;
}

export const ThreadRunItem = (props: Props) => {
  const { loading, threadId } = props;

  const runningRef = useRef(false); // 避免两次effect
  const retryCountRef = useRef(0);

  const [runItem, setRunItem] = useState(props.item);

  // current answer index
  const [curIndex, setCurIndex] = useState(
    () => (props.item.answer?.filter((i) => i.key === "main")?.length ?? 1) - 1
  );

  const { curAnswer, showAnswer, answerLength } = useMemo(() => {
    const mainAnswers = runItem.answer?.filter((i) => i.key === "main") ?? [];
    if (runItem.status >= 2 && mainAnswers.length === 0) {
      mainAnswers.push({
        id: Date.now().toString(),
        key: "main",
        status: 3,
      });
    }
    const curAnswer = mainAnswers.at(curIndex);
    const showAnswer = mainAnswers.length > 0;
    const answerLength = mainAnswers.length;
    return { curAnswer, showAnswer, answerLength };
  }, [runItem, curIndex]);

  const { run: handleStreamScroll } = useThrottleFn(
    () => {
      if (typeof window === "undefined") return;
      const top =
        document.body.scrollHeight || document.documentElement.scrollHeight;
      window.scrollTo({ top, behavior: "smooth" });
    },
    { wait: 100, trailing: true, leading: true }
  );

  const handleRewrite = async () => {
    if (loading || runningRef.current) return;

    props.onStartLoading?.();
    await ApiClient.rewriteThread(threadId, runItem.id);
    handleQuery(runItem.id, true);
  };

  const handleQuery = async (runId: number, rewrite = false) => {
    if (runningRef.current || !threadId) return;
    runningRef.current = true;
    const handleScroll = () => {
      if (rewrite) return;
      handleStreamScroll();
    };
    const handleRefresh = async () => {
      if (retryCountRef.current < 5) {
        await new Promise((r) => setTimeout(r, 1000));
        retryCountRef.current++;
        props.onRefresh?.();
      }
    };

    // 修改索引
    let rewriteChanged = false;

    const finish = await ApiClient.streamThread(threadId, runId, {
      onQuery: () => handleScroll(),
      onSetting(data) {
        setRunItem(
          produce((draft) => {
            draft.setting = data;
          })
        );
        handleScroll();
      },
      onSearch(data) {
        if (data.status === 1) {
          setRunItem(
            produce((draft) => {
              draft.search = [];
            })
          );
          return;
        }
        if (data.status !== 2) return;
        setRunItem(
          produce((draft) => {
            draft.search = (draft.search || []).concat(data.search || []);
          })
        );
        handleScroll();
      },
      onAnswer(data) {
        setRunItem(
          produce((draft) => {
            if (!draft.answer?.length) draft.answer = [];
            let ans = draft.answer.find((i) => i.id === data.id);
            if (!ans) {
              draft.answer.push({
                id: data.id,
                status: 0,
                key: "main",
                content: "",
              });
              ans = draft.answer.at(-1)!;
            }
            ans.status = data.status;
            if (data.model != null) {
              ans.model = data.model;
            }
            if (data.status < 2 && data.delta != null) {
              ans.content = (ans.content || "") + data.delta;
            }
            if (!rewriteChanged) {
              const target = draft.answer.length - 1;
              setTimeout(() => setCurIndex(target), 10);
              rewriteChanged = true;
            }
          })
        );
        handleScroll();
      },
      onError(error) {
        console.log(error);
        runningRef.current = false;
        handleRefresh();
      },
      onDone() {
        console.log("stream done");
        runningRef.current = false;
        handleScroll();
        props.onDone?.();
      },
    });

    if (!finish) {
      runningRef.current = false;
      handleRefresh();
    }
  };

  useEffect(() => {
    setRunItem(props.item);
    setCurIndex(
      (props.item.answer?.filter((i) => i.key === "main")?.length ?? 1) - 1
    );

    if (props.item.status < 2 || props.item.answer?.some((a) => a.status < 2)) {
      props.onStartLoading?.();
      handleQuery(props.item.id);
    }
  }, [props.item]);

  return (
    <div className="space-y-5">
      <h3
        className={clsx(
          "text-2xl font-regular pl-px whitespace-pre-wrap break-words [word-break:break-word]",
          runItem.query?.length > 200 && "text-xl"
        )}
      >
        {runItem.query || "fetch failed..."}
      </h3>

      {/* Think */}
      {(runItem.status < 2 || runItem.setting) && (
        <InnerThink setting={runItem.setting} />
      )}

      {/* Sources */}
      {runItem.search != null && (
        <div className="">
          <div className="flex items-center gap-2 mb-2">
            <div className="ml-1 bg-[var(--main-500)] w-2 h-4 rounded-sm" />
            <div className="text-xl font-regular">Sources</div>
          </div>

          {runItem.search.length > 0 ? (
            <div className="flex gap-2 overflow-x-auto pb-2 box-scrollbar">
              {runItem.search?.map((item, index) => (
                <SourceCard item={item} key={index} />
              ))}
            </div>
          ) : (
            <div className="flex gap-2 overflow-x-hidden pb-2 box-scrollbar">
              <Skeleton className="flex-none w-[200px] h-[118px] rounded-md" />
              <Skeleton className="flex-none w-[200px] h-[118px] rounded-md" />
              <Skeleton className="flex-none w-[200px] h-[118px] rounded-md" />
              <Skeleton className="flex-none w-[200px] h-[118px] rounded-md" />
              <Skeleton className="flex-none w-[200px] h-[118px] rounded-md" />
            </div>
          )}
        </div>
      )}

      {/* Answer */}
      {showAnswer && (
        <div className="">
          <div className="flex items-center gap-2 mb-2">
            <div className="ml-1 bg-[var(--main-500)] w-2 h-4 rounded-sm" />
            <div className="text-xl font-regular">Answer</div>
          </div>

          <div className="transition-all">
            <InnerAnswer answer={curAnswer} />
          </div>

          <div className="mt-1 flex items-center justify-between px-2">
            <div className="flex items-center text-sm gap-2">
              {curAnswer?.model && (
                <div className="text-sm text-muted-foreground">
                  {rewiseModelName(curAnswer.model)}
                </div>
              )}
            </div>

            <div className="flex items-center gap-2 text-muted-foreground">
              {(runItem.answer?.length || 0) > 1 && (
                <div className="flex items-center gap-1 text-sm">
                  <Button
                    size="mini"
                    variant="ghost"
                    className="h-auto p-1"
                    disabled={curIndex === 0}
                    onClick={() => setCurIndex((i) => i - 1)}
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </Button>
                  <div>
                    {curIndex + 1}/{answerLength}
                  </div>
                  <Button
                    size="mini"
                    variant="ghost"
                    className="h-auto p-1"
                    disabled={curIndex === answerLength - 1}
                    onClick={() => setCurIndex((i) => i + 1)}
                  >
                    <ChevronRight className="w-4 h-4" />
                  </Button>
                </div>
              )}
              <Button size="mini" variant="ghost" onClick={handleRewrite}>
                <RefreshCw className="w-4 h-4 mr-0.5" /> <span>Rewrite</span>
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
