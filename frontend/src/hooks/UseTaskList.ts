import { useQuery } from "@tanstack/react-query";
import { getTaskList } from "../gen/task/v1/task-TaskService_connectquery";

// タスク一覧処理のカスタムフック
export const useTaskList = () => {
  const taskQuery = useQuery(getTaskList.useQuery());

  return { taskQuery };
};
