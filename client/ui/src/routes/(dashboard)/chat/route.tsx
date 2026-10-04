import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { ChatNavActions } from "#/components/chat-nav-actions";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from "#/components/ui/breadcrumb.tsx";
import { Separator } from "#/components/ui/separator.tsx";
import { SidebarTrigger } from "#/components/ui/sidebar.tsx";
import ChatInput from "#/components/chat-input.tsx";
import type { ChatData } from "#/types/chat";
import { useChat } from "#/hooks/use-chat.ts";
import { ChatService } from "@nagare-app/services";
import { MessageType, Role } from "@nagare-app/messages";
import { useState } from "react";

export const Route = createFileRoute("/(dashboard)/chat")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const [chatData] = useState<ChatData>({
    text: "",
  });
  const isProcessing = useChat((c) => c.isProcessing);
  const chatSession = useChat((c) => c.chatSession);
  const setChatSession = useChat((c) => c.setChatSession);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setIsConnected = useChat((c) => c.setIsConnected);
  const sessions = useChat((c) => c.sessions);
  const setSessions = useChat((c) => c.setSessions);

  const handleSessionUpdated = (
    updatedSession: NonNullable<typeof chatSession>,
  ) => {
    setChatSession(updatedSession);
    setSessions(
      sessions.map((session) =>
        session.id === updatedSession.id ? updatedSession : session,
      ),
    );
  };

  const handleSessionDeleted = () => {
    if (!chatSession) return;
    setSessions(sessions.filter((session) => session.id !== chatSession.id));
    void navigate({ to: "/chat/new" });
  };

  const handleSessionDuplicated = (
    duplicatedSession: NonNullable<typeof chatSession>,
  ) => {
    setSessions([duplicatedSession, ...sessions]);
    void navigate({ to: `/chat/${duplicatedSession.id}` });
  };

  const onSend = async (data: ChatData) => {
    setPendingMessages([
      {
        id: crypto.randomUUID(),
        role: Role.USER,
        type: MessageType.TextMessageType,
        invokeID: "",
        content: data.text,
      },
    ]);

    if (!chatSession) {
      const session = await ChatService.createChatSession(data.text);
      setIsConnected(false);
      setSessions([...sessions, session]);
      await navigate({
        to: `/chat/${session.id}`,
      });
      return;
    }
  };

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <header className="flex h-14 shrink-0 items-center gap-2 border-b bg-background sticky top-0 z-10">
        <div className="flex flex-1 items-center gap-2 px-3">
          <SidebarTrigger />
          <Separator orientation="vertical" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbPage className="line-clamp-1">
                  {chatSession ? chatSession.title : "Chat"}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
        {chatSession ? (
          <div className="ml-auto px-3">
            <ChatNavActions
              session={chatSession}
              onSessionUpdated={handleSessionUpdated}
              onSessionDeleted={handleSessionDeleted}
              onSessionDuplicated={handleSessionDuplicated}
            />
          </div>
        ) : (
          ""
        )}
      </header>

      <div className="flex-1 min-h-0 overflow-hidden flex flex-col">
        <div className="flex-1 min-h-0 overflow-hidden px-4 py-10 flex flex-col gap-4">
          <Outlet />
        </div>
        {chatSession?.isArchive ? (
          "Chat has been archived"
        ) : (
          <ChatInput
            onSend={onSend}
            currentChatData={chatData}
            disabled={isProcessing}
          />
        )}
      </div>
    </div>
  );
}
