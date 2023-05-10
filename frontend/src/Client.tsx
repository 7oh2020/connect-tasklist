import { TransportProvider } from "@bufbuild/connect-query";
import { Interceptor, createConnectTransport } from "@bufbuild/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { FC, ReactNode } from "react";

type Props = {
  baseUrl: string;
  token?: string;
  children: ReactNode;
};

// connect-queryのセットアップ
export const Client: FC<Props> = ({ baseUrl, token, children }) => {
  const authInterceptor: Interceptor = (next) => async (req) => {
    if (token != null) {
      // リクエストヘッダーにトークンをセットする
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return await next(req);
  };

  const transport = createConnectTransport({
    baseUrl,
    interceptors: [authInterceptor],
  });
  const client = new QueryClient();

  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={client}>{children}</QueryClientProvider>
    </TransportProvider>
  );
};
