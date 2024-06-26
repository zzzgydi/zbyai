import { useLockFn } from "ahooks";
import { SupabaseClient } from "@supabase/supabase-js";
import { FcGoogle } from "react-icons/fc";
import { FaGithub } from "react-icons/fa";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import Logo from "~/assets/images/icon-logo.svg?react";

interface Props {
  supabase: SupabaseClient;
}

export const SigninDialog = ({ supabase }: Props) => {
  const handleGoogle = useLockFn(async () => {
    await supabase.auth.signInWithOAuth({
      provider: "google",
      options: {
        redirectTo: `${window.location.origin}`,
      },
    });
  });

  const handleGithub = useLockFn(async () => {
    await supabase.auth.signInWithOAuth({
      provider: "github",
      options: {
        redirectTo: `${window.location.origin}`,
      },
    });
  });

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <span>Sign in</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-[400px] pt-8">
        <DialogHeader className="">
          <div className="flex items-center gap-0.5 mb-4">
            <Logo className="w-10 h-10" />
          </div>

          <DialogTitle>Sign in</DialogTitle>
          <DialogDescription>
            Welcome back! Sign in to continue.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-3 mt-4">
          <Button variant="secondary" className="gap-1" onClick={handleGoogle}>
            <FcGoogle className="w-5 h-5" />
            <span>Sign in with Google</span>
          </Button>

          <Button variant="secondary" className="gap-1" onClick={handleGithub}>
            <FaGithub className="w-4 h-4" />
            <span>Sign in with Github</span>
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
