import { useNavigation } from "@remix-run/react";

export const PageLoading = () => {
  const navigation = useNavigation();

  if (navigation.state === "idle") return;
  return <div className="fixed top-0 left-0 h-1 z-50 page-loader" />;
};
