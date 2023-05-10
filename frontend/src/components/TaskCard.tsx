import { Button, Divider, Group, List, Text } from "@mantine/core";
import { FC } from "react";
import { Task } from "../gen/task/v1/task_pb";
import { useChangeTask } from "../hooks/UseChangeTask";

type Props = {
  task: Task;
};

// タスク一覧のアイテム
export const TaskCard: FC<Props> = ({ task }) => {
  const { handleComplete, handleUncomplete, handleDelete } = useChangeTask(
    task.id
  );

  return (
    <List.Item my="md">
      <Group>
        <Text>{task.name}</Text>
        <Text>{task.updatedAt?.toDate().toLocaleDateString()}</Text>
      </Group>
      <Group>
        {task.isCompleted ? (
          <Button onClick={handleUncomplete}>未完了に戻す</Button>
        ) : (
          <Button onClick={handleComplete}>完了にする</Button>
        )}

        <Button onClick={handleDelete} color="red">
          削除
        </Button>
      </Group>
      <Divider />
    </List.Item>
  );
};
