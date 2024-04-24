
import { useState, useEffect } from 'react';

const useTodos = (initialSortOrder = "DESC") => {
  const [todos, setTodos] = useState([]);
  const [completedTodos, setCompletedTodos] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [sortOrder, setSortOrder] = useState(initialSortOrder);

  useEffect(() => {
    fetchTodos();
  }, [sortOrder]);

  const fetchTodos = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`http://localhost:8080/api/todos?sort=${sortOrder}`);
      if (!response.ok) throw new Error("Network response was not ok");
      const data = await response.json();
      setTodos(data.filter(todo => !todo.done));
      setCompletedTodos(data.filter(todo => todo.done));
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
      if (!response.ok) throw new Error("Failed to update todo");
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
      if (!response.ok) throw new Error("Failed to save todo");
      fetchTodos();
    } catch (error) {
      console.error("Error saving todo:", error);
      setError("Failed to save todo: " + error.message);
    }
  };

  return { todos, completedTodos, isLoading, error, sortOrder, setSortOrder, fetchTodos, completeTodo, saveTodo };
};

export default useTodos;
