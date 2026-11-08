import { ChatService } from "@nagare-app/services";
import type { Session } from "@nagare-app/dtos";
import { useChat } from "#/hooks/use-chat.ts";
import { ChatModelSelector } from "#/components/chat-model-selector";
import { ChatSessionActionsMenu } from "#/components/chat-session-actions-menu";

type ChatNavActionsProps = {
  session: Session;
  onSessionUpdated: (session: Session) => void;
  onSessionDeleted: () => void;
  onSessionDuplicated: (session: Session) => void;
  onCompact?: () => void;
};

export function ChatNavActions({
  session,
  onSessionUpdated,
  onSessionDeleted,
  onSessionDuplicated,
  onCompact,
}: ChatNavActionsProps) {
  const setInvokeID = useChat((c) => c.setInvokeID);
  const setIsProcessing = useChat((c) => c.setIsProcessing);

  const handleArchive = async () => {
    const updatedSession = await ChatService.archiveChatSession(session.id, {
      isArchive: !session.isArchive,
    });
    onSessionUpdated(updatedSession);
  };

  const handleDuplicate = async () => {
    const duplicatedSession = await ChatService.duplicateChatSession(
      session.id,
    );
    onSessionDuplicated(duplicatedSession);
  };

  const handleDelete = async () => {
    await ChatService.deleteChatSession(session.id);
    onSessionDeleted();
  };

  const handleCompact = async () => {
    const res = await ChatService.compactConversation(session.id);
    if (res?.invokeID) {
      setInvokeID(res.invokeID);
      setIsProcessing(true);
    }
    onCompact?.();
  };

  return (
    <div className="flex items-center gap-2 text-sm">
      <ChatModelSelector
        session={session}
        onSessionUpdated={onSessionUpdated}
      />
      <ChatSessionActionsMenu
        session={session}
        onArchive={handleArchive}
        onDuplicate={handleDuplicate}
        onDelete={handleDelete}
        onCompact={handleCompact}
      />
    </div>
  );
}
