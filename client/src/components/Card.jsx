import { useState } from "react";

export default function Card({ todo, onComplete, onSave }) {
  const [isEditing, setIsEditing] = useState(false);
  const [editText, setEditText] = useState(todo.text);

  const handleComplete = () => {
    onComplete(todo.id);
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSave = async () => {
    onSave(todo.id, editText);
    setIsEditing(false);
  };

  const dateOptions = {
    weekday: "long",
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "numeric",
    minute: "numeric",
  };

  return (
    <div className="flex flex-col bg-white border shadow-sm rounded-xl w-full">
      <div className="p-4 md:p-5">
        {isEditing ? (
          <textarea
            type="text"
            value={editText}
            onChange={(e) => setEditText(e.target.value)}
            className="w-full text-lg font-bold text-gray-800 p-2 rounded"
          />
        ) : (
          <h3
            className="text-lg font-bold text-gray-800 text-wrap p-2"
            onClick={!todo.done && handleEdit}
          >
            {todo.text}
          </h3>
        )}
        <p className="mt-2 text-gray-500 pl-2">
          Created on:{" "}
          {new Date(todo.createdAt).toLocaleString("en-US", dateOptions)}
        </p>
        <p className="text-gray-500 pl-2">
          Last updated:{" "}
          {new Date(todo.updatedAt).toLocaleString("en-US", dateOptions)}
        </p>
        {!todo.done && !isEditing && (
          <button
            onClick={handleComplete}
            className="mt-3 inline-flex items-center gap-x-1 text-sm font-semibold rounded-lg border px-4 py-2 text-blue-600 hover:text-blue-800 disabled:opacity-50"
          >
            Complete!
          </button>
        )}
        {isEditing && (
          <button
            onClick={handleSave}
            className="mt-3 ml-2 inline-flex items-center gap-x-1 text-sm font-semibold rounded-lg border px-4 py-2 text-green-600 hover:text-green-800"
          >
            Save
          </button>
        )}
        {isEditing && (
          <button
            onClick={() => setIsEditing(false)}
            className="mt-3 ml-2 inline-flex items-center gap-x-1 text-sm font-semibold rounded-lg border px-4 py-2 text-red-600 hover:text-red-800"
          >
            Cancel
          </button>
        )}
      </div>
      <div className="bg-gray-100 border-t rounded-b-xl py-3 px-4 md:py-4 md:px-5">
        <p className="mt-1 text-sm text-gray-500">
          Status: {todo.done ? "Completed" : "Pending"}
        </p>
      </div>
    </div>
  );
}
