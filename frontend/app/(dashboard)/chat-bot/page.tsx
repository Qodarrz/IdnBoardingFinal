"use client"

import { useState } from "react"
import { IconSend, IconRobot } from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import Image from "next/image"
import { FormatedChatbot } from "@/components/FormatedChatbot"

type Message = {
  id: number
  sender: "user" | "bot"
  text: string
  time: string
}

export default function ChatBotPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 2,
      sender: "bot",
      text: "Halo selamat datang di chat-bot GreenFlow. silakan bertanya yang berkaitan dengan penggunaan karbon ya!",
      time: "02:22 AM",
    },
  ])

  const [input, setInput] = useState("")
  const [loading, setLoading] = useState(false)

  const handleSend = async () => {
    if (!input.trim()) return

    const newMessage: Message = {
      id: Date.now(),
      sender: "user",
      text: input,
      time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
    }

    setMessages((prev) => [...prev, newMessage])
    setInput("")
    setLoading(true)

    try {
      const token = localStorage.getItem("authtoken")

      const res = await fetch("https://backend-phi-murex-10.vercel.app/api/gemini/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: token ? `Bearer ${token}` : "",
        },
        body: JSON.stringify({
          prompt: input,
        }),
      })

      const data = await res.json()
      console.log(data)

      const botReply: Message = {
        id: Date.now(),
        sender: "bot",
        text: data?.result || "Maaf, terjadi kesalahan.",
        time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
      }

      setMessages((prev) => [...prev, botReply])
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          id: Date.now(),
          sender: "bot",
          text: "⚠️ Gagal mengambil jawaban dari server.",
          time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
        },
      ])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col h-full bg-gray-50">
      <Card className="flex-1 flex flex-col rounded-none">
        <CardContent className="flex-1 overflow-y-auto p-4 space-y-4">
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex items-start gap-2 ${
                msg.sender === "user"
                  ? "justify-end pl-10 md:pl-0"
                  : "pr-10 md:pr-0 justify-start"
              }`}
            >
              {msg.sender === "bot" && (
                <div className="w-8 h-8 aspect-square rounded-full bg-green-100 flex items-center justify-center">
                  <IconRobot size={18} className="text-green-600" />
                </div>
              )}

              <div
                className={`max-w-md p-3 rounded-lg text-sm ${
                  msg.sender === "user"
                    ? "bg-primary text-white"
                    : "bg-white border text-gray-800"
                }`}
              >
                {/* <p>{msg.text}</p> */}
                <FormatedChatbot text={msg.text ?? ''} />
                <span className="text-xs text-gray-400 mt-1 block">{msg.time}</span>
              </div>

              {msg.sender === "user" && (
                <div className="w-8 h-8 aspect-square rounded-full bg-blue-100 flex items-center justify-center">
                  <Image src="/images/profile.png" alt="profile" width={100} height={100} />
                </div>
              )}
            </div>
          ))}

          {loading && (
            <div className="flex items-start gap-2 pr-10 md:pr-0 justify-start">
              <div className="w-8 h-8 aspect-square rounded-full bg-green-100 flex items-center justify-center">
                <IconRobot size={18} className="text-green-600" />
              </div>
              <div className="max-w-md p-3 rounded-lg text-sm bg-white border text-gray-400 italic">
                Bot sedang mengetik...
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Input */}
      <div className="border-t bg-white p-3 flex items-center gap-2">
        <Input
          placeholder="Tulis pesan..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
        />
        <Button onClick={handleSend} size="icon" disabled={loading}>
          <IconSend size={18} />
        </Button>
      </div>
    </div>
  )
}
