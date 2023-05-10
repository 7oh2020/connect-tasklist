import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  completeTask,
  deleteTask,
  getTaskList,
  uncompleteTask,
} from "../gen/task/v1/task-TaskService_connectquery";

// タスク変更処理のカスタムフック
export const useChangeTask = (id: string) => {
  const client = useQueryClient();
  const completeMutation = useMutation(completeTask.useMutation());
  const uncompleteMutation = useMutation(uncompleteTask.useMutation());
  const deleteMutation = useMutation(deleteTask.useMutation());

  const handleComplete = () => {
    completeMutation
      .mutateAsync(
        { taskId: id },
        {
          onSuccess: () => {
            client
              .refetchQueries(getTaskList.getQueryKey())
              .catch((e) => console.error(e));
          },
        }
      )
      .catch((e) => console.error(e));
  };

  const handleUncomplete = () => {
    uncompleteMutation
      .mutateAsync(
        { taskId: id },
        {
          onSuccess: () => {
            client
              .refetchQueries(getTaskList.getQueryKey())
              .catch((e) => console.error(e));
          },
        }
      )
      .catch((e) => console.error(e));
  };

  const handleDelete = () => {
    deleteMutation
      .mutateAsync(
        { taskId: id },
        {
          onSuccess: () => {
            client
              .refetchQueries(getTaskList.getQueryKey())
              .catch((e) => console.error(e));
          },
        }
      )
      .catch((e) => console.error(e));
  };

  return { handleComplete, handleUncomplete, handleDelete };
};
