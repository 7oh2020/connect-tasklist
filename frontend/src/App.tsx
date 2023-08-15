import { MantineProvider, Text, Title } from "@mantine/core";
import { FC, useState } from "react";
import { Client } from "./Client";
import { LoginResponse } from "./gen/auth/v1/auth_pb";
import { HomeRoute } from "./routes/HomeRoute";
import { LoginRoute } from "./routes/LoginRoute";

export const App: FC = () => {
  const [user, setUser] = useState<LoginResponse | null>(null);
  const handleSubmit = (res: LoginResponse) => {
    setUser(res);
  };

  return (
    <Client baseUrl={"http://localhost:8080"} token={user?.token}>
      <MantineProvider withGlobalStyles withNormalizeCSS>
        <header>
          <Title order={1}>connect-list</Title>
          {user != null && <Text>Logged in</Text>}
        </header>
        <main>
          {user != null ? (
            <HomeRoute />
          ) : (
            <LoginRoute onSubmit={handleSubmit} />
          )}
        </main>

        <footer>
          <Text mt="lg">Footer Links</Text>
        </footer>
      </MantineProvider>
    </Client>
  );
};
