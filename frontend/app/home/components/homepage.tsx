"use client";
import React, { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";

type User = { username: string; fullname: string };
type Message = { from: string; to: string; content: string };

const Home = () => {
  const router = useRouter();
  const socketRef = useRef<WebSocket | null>(null);

  const [fullname, setFullname] = useState("");
  const [usernameInput, setUsernameInput] = useState("");
  const [searchedUser, setSearchedUser] = useState<User | null>(null);

  const [chatUsers, setChatUsers] = useState<User[]>([]);
  const [selectedChat, setSelectedChat] = useState<User | null>(null);

  const [messages, setMessages] = useState<Record<string, Message[]>>({});
  const [chatMessage, setChatMessage] = useState("");
  const [error, setError] = useState<string | null>(null);

  /* ---------------- RESTORE SESSION ---------------- */
  useEffect(() => {
    setFullname(sessionStorage.getItem("fullname") || "");
    const storedChats = sessionStorage.getItem("chat_users");
    const storedMessages = sessionStorage.getItem("chat_messages");
    const storedSelected = sessionStorage.getItem("selected_chat");

    if (storedChats) setChatUsers(JSON.parse(storedChats));
    if (storedMessages) setMessages(JSON.parse(storedMessages));
    if (storedSelected) setSelectedChat(JSON.parse(storedSelected));
  }, []);

  /* ---------------- WEBSOCKET CONNECT ---------------- */
  useEffect(() => {
    const token = sessionStorage.getItem("access_token");
    if (!token) return;

    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.onopen = () => console.log("WebSocket connected");

    ws.onmessage = (event) => {
      const msg: Message = JSON.parse(event.data);

      setMessages((prev) => {
        const updated = {
          ...prev,
          [msg.from]: [...(prev[msg.from] || []), msg],
        };
        sessionStorage.setItem("chat_messages", JSON.stringify(updated));
        return updated;
      });

      // Auto-add user to chat list if new
      setChatUsers((prev) => {
        if (prev.find((u) => u.username === msg.from)) return prev;
        const updated = [...prev, { username: msg.from, fullname: msg.from }];
        sessionStorage.setItem("chat_users", JSON.stringify(updated));
        return updated;
      });

      // Auto-open chatbox if user receives a message from a user not currently selected
      setSelectedChat((prev) => prev || { username: msg.from, fullname: msg.from });
      sessionStorage.setItem(
        "selected_chat",
        JSON.stringify({ username: msg.from, fullname: msg.from })
      );
    };

    ws.onclose = () => console.log("WebSocket disconnected");

    socketRef.current = ws;

    return () => ws.close();
  }, []);

  /* ---------------- AUTH FETCH ---------------- */
  async function protectedFetch(url: string) {
    let accessToken = sessionStorage.getItem("access_token");
    const refreshToken = sessionStorage.getItem("refresh_token");

    const res = await fetch(url, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });
    if (res.ok) return res.json();

    // Refresh token flow
    const refreshRes = await fetch("http://localhost:8080/refresh", {
      method: "POST",
      headers: { Authorization: `Bearer ${refreshToken}` },
    });
    if (!refreshRes.ok) {
      sessionStorage.clear();
      router.push("/");
      throw new Error("Session expired");
    }

    const data = await refreshRes.json();
    sessionStorage.setItem("access_token", data.access_token);

    // reconnect websocket
    socketRef.current?.close();
    socketRef.current = new WebSocket(
      `ws://localhost:8080/ws?token=${data.access_token}`
    );

    const retry = await fetch(url, {
      headers: { Authorization: `Bearer ${data.access_token}` },
    });
    return retry.json();
  }

  /* ---------------- SEARCH USER ---------------- */
  async function searchUser() {
    setError(null);
    setSearchedUser(null);

    const currentUsername = sessionStorage.getItem("username");
    if (!usernameInput || usernameInput === currentUsername) {
      setError("Invalid username");
      return;
    }

    try {
      const user = await protectedFetch(
        `http://localhost:8080/user/${usernameInput}`
      );
      setSearchedUser(user);
    } catch (err: any) {
      setError(err.message);
    }
  }

  /* ---------------- ADD CHAT ---------------- */
  function addChat(user: User) {
    if (!chatUsers.find((u) => u.username === user.username)) {
      const updated = [...chatUsers, user];
      setChatUsers(updated);
      sessionStorage.setItem("chat_users", JSON.stringify(updated));
    }

    setSelectedChat(user);
    sessionStorage.setItem("selected_chat", JSON.stringify(user));

    setSearchedUser(null);
    setUsernameInput("");
  }

  /* ---------------- SEND MESSAGE ---------------- */
  function sendMessage() {
    if (!selectedChat || !chatMessage.trim()) return;
    if (!socketRef.current) return;

    const msg: Message = {
      from: sessionStorage.getItem("username")!,
      to: selectedChat.username,
      content: chatMessage,
    };

    socketRef.current.send(JSON.stringify({ to: msg.to, content: msg.content }));

    setMessages((prev) => {
      const updated = {
        ...prev,
        [selectedChat.username]: [...(prev[selectedChat.username] || []), msg],
      };
      sessionStorage.setItem("chat_messages", JSON.stringify(updated));
      return updated;
    });

    setChatMessage("");
  }

  return (
    <div className="flex min-h-screen p-4 gap-6">
      {/* LEFT */}
      <div className="w-1/3 flex flex-col gap-4">
        <h1 className="text-2xl font-bold">Welcome, {fullname}!</h1>

        <div className="flex gap-2">
          <input
            value={usernameInput}
            onChange={(e) => setUsernameInput(e.target.value)}
            placeholder="Search username"
            className="border rounded px-2 py-1 w-full"
          />
          <button onClick={searchUser} className="bg-blue-500 text-white px-4 rounded">
            Search
          </button>
        </div>

        {searchedUser && (
          <div className="flex justify-between border p-2 rounded">
            <span>{searchedUser.fullname}</span>
            <button
              onClick={() => addChat(searchedUser)}
              className="bg-green-500 text-white px-3 rounded"
            >
              Chat
            </button>
          </div>
        )}

        {chatUsers.map((u) => (
          <div
            key={u.username}
            onClick={() => setSelectedChat(u)}
            className={`p-2 rounded cursor-pointer ${
              selectedChat?.username === u.username ? "bg-blue-100" : "bg-gray-100"
            }`}
          >
            {u.fullname}
          </div>
        ))}
      </div>

      {/* RIGHT */}
      <div className="flex-1 border rounded p-4 flex flex-col">
        {selectedChat ? (
          <>
            <h2 className="font-semibold mb-2">Chat with {selectedChat.fullname}</h2>

            <div className="flex-1 overflow-y-auto border p-2 mb-2 flex flex-col gap-1">
              {(messages[selectedChat.username] || []).map((m, i) => (
                <div
                  key={i}
                  className={`p-1 rounded max-w-xs ${
                    m.from === sessionStorage.getItem("username")
                      ? "self-end bg-blue-200"
                      : "self-start bg-gray-200"
                  }`}
                >
                  <div className="text-xs font-semibold">{m.from}</div>
                  <div>{m.content}</div>
                </div>
              ))}
            </div>

            <div className="flex gap-2">
              <input
                value={chatMessage}
                onChange={(e) => setChatMessage(e.target.value)}
                className="flex-1 border rounded px-2"
              />
              <button
                onClick={sendMessage}
                className="bg-blue-600 text-white px-4 rounded"
              >
                Send
              </button>
            </div>
          </>
        ) : (
          <p className="text-gray-400 text-center mt-20">Select a chat</p>
        )}
      </div>
    </div>
  );
};

export default Home;
