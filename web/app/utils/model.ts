export function rewiseModelName(model: string): string {
  if (!model) return "";
  if (model.includes("/")) {
    return model.split("/").pop() || "";
  }

  if (model.startsWith("gpt-3.5-turbo")) {
    return "gpt-3.5-turbo";
  }

  return model;
}
