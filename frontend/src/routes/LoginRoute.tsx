import {
  Alert,
  Box,
  Button,
  Group,
  PasswordInput,
  TextInput,
  Title,
} from "@mantine/core";
import { FC } from "react";
import { LoginResponse } from "../gen/user/v1/user_pb";
import { useLoginUser } from "../hooks/UseLoginUser";

type Props = {
  onSubmit: (user: LoginResponse) => void;
};

// ログインページ
export const LoginRoute: FC<Props> = ({ onSubmit }) => {
  const { form, loginMutation, handleSubmit } = useLoginUser(onSubmit);

  return (
    <>
      <Title order={2}>Login</Title>
      <Box>
        <form onSubmit={handleSubmit}>
          {loginMutation.isError && (
            <Alert title="Error" color="red">
              {loginMutation.error.message}
            </Alert>
          )}
          <Group mt="md">
            <TextInput
              {...form.getInputProps("email")}
              label="Email"
              autoFocus={true}
              placeholder="Email"
            />
          </Group>
          <Group mt="md">
            <PasswordInput
              {...form.getInputProps("pass")}
              label="Pass"
              placeholder="password"
            />
          </Group>
          <Group mt="md">
            <Button type="submit" loading={loginMutation.isLoading}>
              Login
            </Button>
          </Group>
        </form>
      </Box>
    </>
  );
};
