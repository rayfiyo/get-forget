import React, { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "./components/ui/card"; // 修正したパス
import { Brain, MessageSquare, RefreshCcw, AlertTriangle } from "lucide-react";

const GetForget = () => {
  const [messages, setMessages] = useState([]);
  const [memories, setMemories] = useState({});
  const [input, setInput] = useState("");
  const [forgettingProcess, setForgettingProcess] = useState(false);

  // メモリーの重要度を計算
  const calculateImportance = (memory) => {
    const age = Date.now() - memory.timestamp;
    const ageWeight = Math.exp(-age / (1000 * 60 * 60)); // 時間経過で重要度低下
    const useCount = memory.useCount || 1;
    return memory.initialImportance * ageWeight * Math.log(useCount + 1);
  };

  // 定期的な忘却プロセス
  useEffect(() => {
    const interval = setInterval(() => {
      setForgettingProcess(true);
      setMemories((prev) => {
        const newMemories = { ...prev };
        Object.keys(newMemories).forEach((key) => {
          const importance = calculateImportance(newMemories[key]);
          // ランダムな忘却の確率
          if (Math.random() > importance / 100) {
            delete newMemories[key];
            addMessage({
              type: "system",
              content: `記憶が曖昧になってきました: "${key}"`,
            });
          }
        });
        return newMemories;
      });
      setForgettingProcess(false);
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const addMessage = (message) => {
    setMessages((prev) => [
      ...prev,
      {
        ...message,
        id: Date.now(),
        timestamp: new Date().toISOString(),
      },
    ]);
  };

  const handleSend = () => {
    if (!input.trim()) return;

    // ユーザーメッセージの追加
    addMessage({
      type: "user",
      content: input,
    });

    // メモリーの作成/更新
    const keywords = input.split(" ");
    keywords.forEach((keyword) => {
      if (keyword.length > 3) {
        setMemories((prev) => ({
          ...prev,
          [keyword]: {
            content: input,
            timestamp: Date.now(),
            initialImportance: Math.random() * 50 + 50, // 50-100の重要度
            useCount: (prev[keyword]?.useCount || 0) + 1,
          },
        }));
      }
    });

    // AIの応答生成
    const rememberedContext = Object.keys(memories)
      .filter((key) => input.includes(key))
      .map((key) => memories[key].content);

    let response;
    if (rememberedContext.length > 0) {
      response = `以前の会話を思い出しました: "${rememberedContext[0]}"`;
    } else {
      response = "申し訳ありません。関連する記憶が曖昧です...";
    }

    // AI応答の追加
    addMessage({
      type: "ai",
      content: response,
    });

    setInput("");
  };

  return (
    <div className="w-full max-w-4xl p-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Brain className="w-6 h-6" />
            記憶の曖昧化チャットシステム
            {forgettingProcess && (
              <RefreshCcw className="w-4 h-4 animate-spin" />
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-96 overflow-y-auto mb-4 p-4 border rounded-lg">
            {messages.map((message) => (
              <div
                key={message.id}
                className={`mb-4 p-2 rounded-lg ${
                  message.type === "user"
                    ? "bg-blue-100 ml-auto max-w-[80%]"
                    : message.type === "system"
                    ? "bg-yellow-100 max-w-full"
                    : "bg-gray-100 mr-auto max-w-[80%]"
                }`}
              >
                {message.type === "system" && (
                  <AlertTriangle className="w-4 h-4 inline mr-2" />
                )}
                {message.content}
              </div>
            ))}
          </div>

          <div className="flex gap-2">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyPress={(e) => e.key === "Enter" && handleSend()}
              className="flex-1 p-2 border rounded-lg"
              placeholder="メッセージを入力..."
            />
            <button
              onClick={handleSend}
              className="p-2 bg-blue-500 text-white rounded-lg"
            >
              <MessageSquare className="w-6 h-6" />
            </button>
          </div>

          <div className="mt-4">
            <h3 className="font-semibold mb-2">現在の記憶状態:</h3>
            <div className="grid grid-cols-2 gap-2">
              {Object.entries(memories).map(([key, memory]) => (
                <div key={key} className="p-2 border rounded-lg">
                  <div className="font-medium">{key}</div>
                  <div className="text-sm text-gray-600">
                    重要度: {calculateImportance(memory).toFixed(2)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default GetForget;

/*
import logo from './logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
*/
