import { KBDetail, NodeDetail, NodeListItem, WidgetInfo } from '@/assets/type';

interface ApiClientConfig {
  kb_id?: string;
  headers?: Record<string, string>;
  cache?: RequestCache;
  authToken?: string
  session?: string
}

interface Response<T> {
  data?: T;
  status: number;
  error?: string;
}

class ApiClient {
  private baseURL: string;

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || '';
  }

  // 服务端专用方法 - 带cookie的请求
  async request<T>(
    url: string,
    options: RequestInit = {},
    config: ApiClientConfig = {}
  ): Promise<Response<T>> {
    const { kb_id = '', headers = {}, cache, authToken, session } = config;

    const requestHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
      'x-kb-id': kb_id,
      ...headers,
    };

    if (session) requestHeaders['x-pw-session-id'] = session;
    if (authToken) requestHeaders['x-simple-auth-password'] = authToken

    const fullUrl = url.startsWith('http') ? url : `${this.baseURL}${url}`;
    try {
      const response = await fetch(fullUrl, {
        ...options,
        headers: requestHeaders,
        ...(cache && { cache }),
      });
      if (!response.ok) {
        return { status: response.status, error: `HTTP error! status: ${response.status}` };
      }
      const result = await response.json();
      return { data: result.data, status: response.status };
    } catch (error) {
      return { status: 500, error: error instanceof Error ? error.message : 'Unknown error' };
    }
  }

  // 服务端获取知识库信息
  async serverGetKBInfo(kb_id: string, authToken?: string): Promise<Response<KBDetail>> {
    return this.request(`/share/v1/app/web/info`, {
      method: 'GET',
    }, {
      kb_id,
      authToken,
      cache: 'no-store',
    });
  }

  async serverGetWidgetInfo(kb_id: string): Promise<Response<WidgetInfo>> {
    return this.request(`/share/v1/app/widget/info`, {
      method: 'GET',
    }, {
      kb_id,
      cache: 'no-store',
    });
  }

  // 服务端获取节点列表
  async serverGetNodeList(
    kb_id: string,
    authToken?: string,
  ): Promise<Response<NodeListItem[]>> {
    return this.request(`/share/v1/node/list`, {
      method: 'GET',
    }, {
      kb_id,
      authToken,
    });
  }
  // 服务端获取节点列表
  async clientGetNodeList(
    kb_id: string,
    authToken?: string,
  ): Promise<Response<NodeListItem[]>> {
    return this.request(window?.location.origin + `/client/v1/node/list`, {
      method: 'GET',
    }, {
      kb_id,
      authToken,
    });
  }

  // 服务端获取节点详情
  async serverGetNodeDetail(id: string, kb_id: string, authToken?: string): Promise<Response<NodeDetail>> {
    return this.request(`/share/v1/node/detail?id=${id}`, {
      method: 'GET',
    }, {
      kb_id,
      authToken,
    });
  }

  // 客服端请求
  async clientGetNodeDetail(id: string, kb_id: string, authToken?: string): Promise<Response<NodeDetail>> {
    return this.request(window?.location.origin + `/client/v1/node/detail?id=${id}`, {
      method: 'GET',
    }, {
      kb_id,
      authToken,
    });
  }

  // 服务端页面埋点
  async serviceStatPage(data: { node_id: string, scene: number, kb_id: string, authToken?: string, headers?: Record<string, string>, session: string }): Promise<Response<void>> {
    return this.request('/share/v1/stat/page', {
      method: 'POST',
      body: JSON.stringify({
        node_id: data.node_id,
        scene: data.scene
      }),
    }, {
      kb_id: data.kb_id,
      authToken: data.authToken,
      session: data.session,
      headers: data.headers,
    });
  }

  // 客服端页面埋点
  async clientStatPage(data: { node_id: string, scene: number, kb_id: string, authToken?: string }): Promise<Response<void>> {
    return this.request(window?.location.origin + '/client/v1/stat/page', {
      method: 'POST',
      body: JSON.stringify({
        node_id: data.node_id,
        scene: data.scene
      }),
    }, {
      kb_id: data.kb_id,
      authToken: data.authToken,
    });
  }

  // 客服端反馈
  async clientFeedback(data: {
    kb_id: string,
    authToken?: string,
    conversation_id: string,
    message_id: string,
    type: number,
    score: number,
    feedback_content?: string
  }): Promise<Response<void>> {
    return this.request(window?.location.origin + '/client/v1/chat/feedback', {
      method: 'POST',
      body: JSON.stringify({
        conversation_id: data.conversation_id,
        message_id: data.message_id,
        type: data.type,
        score: data.score,
        feedback_content: data.feedback_content,
      }),
    }, {
      kb_id: data.kb_id,
      authToken: data.authToken,
    });
  }
}

export const apiClient = new ApiClient();
