import { useState } from "react";

export default function Form({ fetchTodos }) {
  const [text, setText] = useState("");

  const handleSubmit = async (event) => {
    event.preventDefault();

    const todo = { text, done: false };

    try {
      const response = await fetch("http://localhost:8080/api/todos", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(todo),
      });
      console.log(response);
      console.log(response.status === 201);

      if (!response.ok) {
        throw new Error("Failed to create todo");
      }

      setText("");
      fetchTodos();
    } catch (error) {
      console.error("Error creating todo:", error);
      alert("Failed to create todo");
    }
  };

  return (
    <div className="w-full my-4">
      <div className="mx-auto max-w-2xl">
        <div className="p-4 relative z-10 bg-white border rounded-xl md:p-10">
          <form onSubmit={handleSubmit}>
            <div className="mb-4 sm:mb-8">
              <textarea
                type="text"
                id="todo-text"
                className="py-3 px-4 block w-full border rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none"
                placeholder="What needs to be done?"
                value={text}
                onChange={(e) => setText(e.target.value)}
              />
            </div>

            <div className="mt-6 grid">
              <button
                type="submit"
                className="w-full py-3 px-4 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none"
              >
                Create
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
