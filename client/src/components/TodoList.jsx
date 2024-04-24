import { Card } from ".";

export default function TodoList({
  todos,
  isLoading,
  error,
  onComplete,
  onSave,
}) {
  return (
    <>
      <div className="my-4 p-2 flex flex-col gap-y-2 h-full overflow-y-auto border rounded-xl">
        {error && <p className="text-red-500">{error}</p>}
        {isLoading ? (
          <p className="text-gray-600 text-center mt-4">Loading todos...</p>
        ) : todos.length > 0 ? (
          <>
            {todos.map((todo) => (
              <Card
                key={todo.id}
                todo={todo}
                onComplete={onComplete}
                onSave={onSave}
              />
            ))}
          </>
        ) : (
          <p className="text-gray-600 text-center mt-4">
            No todos to show. Add some new tasks!
          </p>
        )}
      </div>
    </>
  );
}
