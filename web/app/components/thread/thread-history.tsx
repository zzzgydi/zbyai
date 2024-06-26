import useSWR from "swr";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { Link } from "@remix-run/react";
import { ApiClient } from "~/services/api.client";
import { X } from "lucide-react";
import { useLockFn } from "ahooks";
import { toast } from "../ui/use-toast";

dayjs.extend(relativeTime);

export const ThreadHistory = ({ onClose }: { onClose: () => void }) => {
  const {
    data: threads,
    error,
    isLoading,
    mutate,
  } = useSWR("/api/thread/history", () => ApiClient.listThread());

  const handleDelete = useLockFn(async (id: string) => {
    try {
      await ApiClient.deleteThread(id);
      await mutate(
        threads?.filter((t) => t.id !== id),
        false
      );
    } catch (err: any) {
      toast({
        title: "Failed to delete thread 😔",
        description: err.message || err.toString(),
      });
    }
  });

  if (isLoading) {
    return (
      <div className="container text-muted-foreground pt-[10vh] pb-10 text-center">
        Loading...
      </div>
    );
  }

  // Sign in to view thread history...
  if (error) {
    return (
      <div className="relative container pb-10">
        <p className="text-muted-foreground pt-[10vh] pb-10 text-center">
          {error?.toString() || "Failed to loading threads."}
        </p>
      </div>
    );
  }

  return (
    <div className="flex-auto overflow-auto px-6 pt-2 pb-6">
      <ul className="space-y-3">
        {threads?.map((thread) => (
          <li key={thread.id}>
            <Link
              to={`/thread/${thread.id}`}
              className="block rounded-lg border px-3 py-2 text-sm transition-all hover:bg-accent relative group"
              onClick={() => onClose()}
            >
              <div className="space-y-0.5 pr-2">
                <h3
                  className="font-semibold truncate leading-5"
                  title={thread.title}
                >
                  {thread.title || "New Thread"}
                </h3>
                <p
                  className="flex-none font-normal text-muted-foreground leading-4"
                  title={thread.created_at}
                >
                  {dayjs(thread.created_at).fromNow()}
                </p>
              </div>

              <div className="hidden group-hover:block absolute right-2 top-2">
                <button
                  type="button"
                  className="rounded-sm ring-offset-background transition-all hover:scale-125 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    handleDelete(thread.id);
                  }}
                >
                  <X className="h-4 w-4" />
                  <span className="sr-only">Close</span>
                </button>
              </div>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
};
