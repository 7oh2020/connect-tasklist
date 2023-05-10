import { useForm } from "@mantine/form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createTask,
  getTaskList,
} from "../gen/task/v1/task-TaskService_connectquery";

type FormValue = {
  name: string;
};

// タスク追加処理のカスタムフック
export const useNewTask = () => {
  const client = useQueryClient();
  const createMutation = useMutation(createTask.useMutation());
  const form = useForm<FormValue>({
    initialValues: {
      name: "",
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    createMutation
      .mutateAsync(
        {
          name: values.name,
        },
        {
          onSuccess: () => {
            client
              .refetchQueries(getTaskList.getQueryKey())
              .catch((e) => console.error(e));
          },
        }
      )
      .catch((e) => console.error(e));
    form.reset();
  });

  return { form, createMutation, handleSubmit };
};
