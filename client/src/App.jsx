import useTodos from "./hooks/useTodos";
import { Form, TodoList } from "./components";

function App() {
  const {
    todos,
    completedTodos,
    isLoading,
    error,
    sortOrder,
    setSortOrder,
    fetchTodos,
    completeTodo,
    saveTodo,
  } = useTodos();

  const toggleSortOrder = () => {
    const newSortOrder = sortOrder === "ASC" ? "DESC" : "ASC";
    setSortOrder(newSortOrder);
  };

  return (
    <>
      <div className="m-2 flex flex-row flex-wrap w-full h-screen">
        <div className="flex flex-col p-4 w-full sm:w-1/2 h-screen">
          <div className="flex flex-row justify-between">
            <h2 className="text-xl text-gray-800 font-bold sm:text-3xl">
              To do
            </h2>
            <button
              onClick={toggleSortOrder}
              className="inline-flex items-center gap-x-1 text-sm font-semibold rounded-lg border px-4 text-blue-600 hover:text-blue-800 disabled:opacity-50"
            >
              Sort {sortOrder === "ASC" ? "Descending" : "Ascending"}
            </button>
          </div>
          <TodoList
            todos={todos}
            isLoading={isLoading}
            error={error}
            onComplete={completeTodo}
            onSave={saveTodo}
          />
        </div>
        <div className="flex flex-col p-4 w-full sm:w-1/2 h-screen">
          <h2 className="text-xl text-gray-800 font-bold sm:text-3xl">
            Create Todo
          </h2>
          <Form fetchTodos={fetchTodos} />
          <h2 className="text-xl text-gray-800 font-bold sm:text-3xl">
            Completed todos
          </h2>
          <TodoList
            todos={completedTodos}
            isLoading={isLoading}
            error={error}
          />
        </div>
      </div>
    </>
  );
}

export default App;
