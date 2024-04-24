import { useState, useEffect } from "react";
import { Card, Form } from "./components";

function App() {
  const [todos, setTodos] = useState([]);
  const [completedTodos, setCompletedTodos] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [sortOrder, setSortOrder] = useState("DESC");

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(
        `http://localhost:8080/api/todos?sort=${sortOrder}`,
      );
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const data = await response.json();
      setTodos(data.filter((todo) => !todo.done));
      setCompletedTodos(data.filter((todo) => todo.done));
    } catch (error) {
      setError("Failed to fetch todos: " + error.message);
    } finally {
      setIsLoading(false);
    }
  };

  const completeTodo = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/api/todos/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ done: true }),
      });
      if (!response.ok) {
        throw new Error("Failed to update todo");
      }
      fetchTodos();
    } catch (error) {
      console.error("Error completing todo:", error);
      setError("Failed to complete todo: " + error.message);
    }
  };

  const saveTodo = async (id, text) => {
    try {
      const response = await fetch(`http://localhost:8080/api/todos/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: text }),
      });
      console.log(response);
      if (!response.ok) {
        throw new Error("Failed to save todo");
      }
      fetchTodos();
    } catch (error) {
      console.error("Error saving todo:", error);
      setError("Failed to save todo: " + error.message);
    }
  };

  const toggleSortOrder = () => {
    const newSortOrder = sortOrder === "ASC" ? "DESC" : "ASC";
    setSortOrder(newSortOrder);
    fetchTodos();
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
              className="w-1/2 px-4 py-2 mb-2 bg-blue-500 text-white rounded hover:bg-blue-700"
            >
              Sort {sortOrder === "ASC" ? "Ascending" : "Descending"}
            </button>
          </div>
          {error && <p className="text-red-500">{error}</p>}
          {isLoading ? (
            <p>Loading todos...</p>
          ) : (
            <div className="my-4 p-2 flex flex-col gap-y-2 h-full overflow-y-auto border rounded-xl">
              {todos.map((todo) => (
                <Card
                  key={todo.id}
                  todo={todo}
                  onComplete={completeTodo}
                  onSave={saveTodo}
                />
              ))}
            </div>
          )}
        </div>
        <div className="flex flex-col pt-4 px-4 w-full sm:w-1/2 h-screen">
          <h2 className="text-xl text-gray-800 font-bold sm:text-3xl">
            Create Todo
          </h2>
          <Form fetchTodos={fetchTodos} />
          <h2 className="text-xl text-gray-800 font-bold sm:text-3xl">
            Completed todos
          </h2>
          {isLoading ? (
            <p>Loading completed todos...</p>
          ) : (
            <div className="mt-4 p-2 flex flex-col gap-y-2 h-1/2 overflow-y-auto border rounded-xl">
              {completedTodos.map((todo) => (
                <Card key={todo.id} todo={todo} />
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default App;
