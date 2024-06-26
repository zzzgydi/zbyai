export const API_BASE_URL =
  {
    development: "http://localhost:14090/v1",
    test: "/v1",
    production: "/v1",
  }[process.env.NODE_ENV!] || "/v1";

export const metaDescription =
  "ZbyAI enhances your search with AI for deeper insights and more relevant results. Experience unparalleled accuracy and discover the web's wisdom.";
