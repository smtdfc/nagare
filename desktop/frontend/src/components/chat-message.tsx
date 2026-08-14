type ChatMessageProps = {
  index: string
  role: string
  content: string
  isErrMessage?: boolean
}

export default function ChatMessage({
  index,
  role,
  content,
  isErrMessage,
}: ChatMessageProps) {
  return (
    <div
      key={index}
      className={`flex gap-3 w-full ${
        role === 'user' ? 'justify-end' : 'justify-start'
      }`}
    >
      <div
        className={`max-w-[80%] md:max-w-[70%] p-4 rounded-2xl text-sm leading-relaxed whitespace-pre-wrap wrap-break-word shadow-sm ${
          role === 'user'
            ? 'bg-primary text-primary-foreground rounded-br-xs'
            : 'bg-card border text-card-foreground rounded-bl-xs'
        } ${isErrMessage ? 'bg-red-100 text-red-800' : ''} `}
      >
        {content}
      </div>
    </div>
  )
}
