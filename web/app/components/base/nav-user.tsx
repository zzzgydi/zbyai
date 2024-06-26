import { Session, SupabaseClient } from "@supabase/supabase-js";
import { Suspense, useEffect, useState } from "react";
import { Skeleton } from "../ui/skeleton";
import { Await, useFetcher, useRevalidator } from "@remix-run/react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "~/components/ui/sheet";
import { SigninDialog } from "./signin-dialog";
import { ThreadHistory } from "../thread/thread-history";

interface Props {
  session?: Promise<Session | null>;
  supabase: SupabaseClient;
}

export const NavUser = (props: Props) => {
  const { session, supabase } = props;

  const fetcher = useFetcher();
  const { revalidate } = useRevalidator();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange(async (event, newSession) => {
      const serverSession = await session;
      console.log(
        "auth change",
        event,
        newSession,
        newSession?.access_token === serverSession?.access_token
      );
      if (
        event !== "INITIAL_SESSION" &&
        newSession?.access_token !== serverSession?.access_token
      ) {
        console.log("revalidate");
        // server and client are out of sync.
        revalidate();
      }
    });

    return () => {
      subscription.unsubscribe();
    };
  }, [session, supabase, revalidate]);

  const handleSignOut = () => {
    if (fetcher.state !== "idle") return;
    supabase.auth.signOut();
    fetcher.submit({}, { method: "POST", action: "/auth/signout" });
  };

  return (
    <Suspense fallback={<Skeleton className="w-7 h-7 rounded-full" />}>
      <Await resolve={session}>
        {(session) => {
          const user = session?.user;

          return !user ? (
            <SigninDialog supabase={supabase} />
          ) : (
            <>
              <DropdownMenu>
                <DropdownMenuTrigger className="outline-none">
                  <Avatar className="w-7 h-7">
                    <AvatarImage src={user.user_metadata?.avatar_url} />
                    <AvatarFallback className="select-none">
                      {(user.user_metadata?.name || "User")
                        .at(0)
                        ?.toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuLabel>
                    {user.user_metadata?.name || "User"}
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => setOpen(true)}>
                    History
                  </DropdownMenuItem>

                  <DropdownMenuItem onClick={handleSignOut}>
                    Sign Out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>

              <Sheet open={open} onOpenChange={setOpen}>
                <SheetContent className="overflow-hidden flex flex-col p-0">
                  <SheetHeader className="flex-none px-6 pt-6">
                    <SheetTitle>Thread History</SheetTitle>
                  </SheetHeader>

                  <ThreadHistory onClose={() => setOpen(false)} />
                </SheetContent>
              </Sheet>
            </>
          );
        }}
      </Await>
    </Suspense>
  );
};
