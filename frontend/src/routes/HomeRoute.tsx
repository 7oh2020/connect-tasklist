import { List, Text, Title } from "@mantine/core";
import { FC } from "react";
import { NewTaskForm } from "../components/NewTaskForm";
import { TaskCard } from "../components/TaskCard";
import { useTaskList } from "../hooks/UseTaskList";

// タスクの一覧と操作ページ
export const HomeRoute: FC = () => {
  const { taskQuery } = useTaskList();

  return (
    <>
      <Title order={2}>Tasks</Title>
      <NewTaskForm />
      <List sx={{ listStyle: "none" }}>
        {taskQuery.isLoading && <Text>Loading...</Text>}
        {taskQuery.isError && <Text>{taskQuery.error.message}</Text>}
        {!taskQuery.isLoading &&
          !taskQuery.isError &&
          taskQuery.data.tasks.map((task) => (
            <TaskCard key={task.id} task={task} />
          ))}
      </List>
    </>
  );
};
