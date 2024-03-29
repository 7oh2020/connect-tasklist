import { useMutation } from "@tanstack/react-query";
import { login } from "../gen/auth/v1/auth-AuthService_connectquery";
import { LoginResponse } from "../gen/auth/v1/auth_pb";
import { useForm } from "@mantine/form";

type FormValue = {
  email: string;
  pass: string;
};

// ログイン処理のカスタムフック
export const useLoginUser = (onSubmit: (user: LoginResponse) => void) => {
  const loginMutation = useMutation(login.useMutation());
  const form = useForm<FormValue>({
    initialValues: {
      email: "",
      pass: "",
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    loginMutation
      .mutateAsync({
        email: values.email,
        password: values.pass,
      })
      .then((res) => onSubmit(res))
      .catch((e) => console.error(e));
  });

  return { form, loginMutation, handleSubmit };
};
