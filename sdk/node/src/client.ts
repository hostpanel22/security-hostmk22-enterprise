export interface HostMKClientOptions {
  apiKey: string;
  baseUrl: string;
  timeout?: number;
}

export interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
}

export class HostMKError extends Error {
  public readonly status: number;
  public readonly data: unknown;

  constructor(
    message: string,
    status: number,
    data?: unknown,
  ) {
    super(message);

    this.name = "HostMKError";
    this.status = status;
    this.data = data;
  }
}

export class HostMKClient {
  private readonly apiKey: string;

  private readonly baseUrl: string;

  private readonly timeout: number;

  constructor(
    options: HostMKClientOptions,
  ) {
    if (!options.apiKey) {
      throw new Error(
        "apiKey is required",
      );
    }

    if (!options.baseUrl) {
      throw new Error(
        "baseUrl is required",
      );
    }

    this.apiKey = options.apiKey;

    this.baseUrl = options.baseUrl.replace(
      /\/$/,
      "",
    );

    this.timeout =
      options.timeout ?? 30000;
  }

  public async get<T>(
    path: string,
  ): Promise<T> {
    return this.request<T>(
      path,
      {
        method: "GET",
      },
    );
  }

  public async post<T>(
    path: string,
    body?: unknown,
  ): Promise<T> {
    return this.request<T>(
      path,
      {
        method: "POST",
        body,
      },
    );
  }

  public async put<T>(
    path: string,
    body?: unknown,
  ): Promise<T> {
    return this.request<T>(
      path,
      {
        method: "PUT",
        body,
      },
    );
  }

  public async delete<T>(
    path: string,
  ): Promise<T> {
    return this.request<T>(
      path,
      {
        method: "DELETE",
      },
    );
  }

  private async request<T>(
    path: string,
    options: RequestOptions,
  ): Promise<T> {
    const controller =
      new AbortController();

    const timeout = setTimeout(
      () => controller.abort(),
      this.timeout,
    );

    try {
      const response =
        await fetch(
          `${this.baseUrl}${path}`,
          {
            method:
              options.method ??
              "GET",

            signal:
              controller.signal,

            headers: {
              Authorization:
                `Bearer ${this.apiKey}`,

              "Content-Type":
                "application/json",

              ...options.headers,
            },

            body:
              options.body !== undefined
                ? JSON.stringify(
                    options.body,
                  )
                : undefined,
          },
        );

      let data: unknown;

      const contentType =
        response.headers.get(
          "content-type",
        ) ?? "";

      if (
        contentType.includes(
          "application/json",
        )
      ) {
        data =
          await response.json();
      } else {
        data =
          await response.text();
      }

      if (!response.ok) {
        throw new HostMKError(
          response.statusText,
          response.status,
          data,
        );
      }

      return data as T;
    } finally {
      clearTimeout(
        timeout,
      );
    }
  }
}