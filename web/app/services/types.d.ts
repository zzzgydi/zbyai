interface Thread {
  id: string;
  user_id: string;
  visible: boolean;
  title?: string;
  options?: ThreadOption;
  created_at: string;
  updated_at: string;
}

interface ThreadDetail {
  id: string;
  title: string;
  history: ThreadRun[];
  current: number;
  status: number;
}

interface ThreadRun {
  id: number;
  status: number;
  query: string;
  setting?: ThreadSetting;
  answer?: ThreadAnswer[];
  search?: ThreadSearch[];
  created_at: string;
  updated_at: string;
}

interface ThreadSetting {
  model?: string;
  query_list?: string[];
  use_search?: boolean;
}

interface ThreadAnswer {
  id: string;
  run_id?: number;
  key: string;
  status: number;
  model?: string;
  content?: string;
  created_at?: string;
}

interface ThreadOption {
  model: string;
}

interface ThreadSearch {
  id: string;
  run_id: number;
  // status: number;
  title: string;
  link: string;
  snippet: string;
  page?: string;
}

interface SearchResult {
  items?: SearchItem[];
  related?: SearchRelated[];
}

interface SearchItem {
  title: string;
  link: string;
  snippet: string;
  sitelinks?: SerperSiteLink[];
}

interface SearchRelated {
  title?: string;
  link?: string;
  snippet?: string;
}

interface SerperSiteLink {
  title: string;
  link: string;
}
