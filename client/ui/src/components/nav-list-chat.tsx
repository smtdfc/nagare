"use client";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  MoreHorizontalIcon,
  LinkIcon,
  ArrowUpRightIcon,
  Trash2Icon,
} from "lucide-react";
import type { Session } from "@nagare-app/dtos";
import { Link, useLocation, useNavigate } from "@tanstack/react-router";
import { Spinner } from "#/components/ui/spinner";
import { ChatService, envConfig, Environment } from "@nagare-app/services";
import { useChat } from "#/hooks/use-chat";
import { showToastError } from "#/lib/toast";
import { ArchiveIcon, CopyIcon } from "lucide-react";

type NavListChatProps = {
  sessions: Session[];
  isLoading?: boolean;
};

export function NavListChat({ sessions, isLoading = false }: NavListChatProps) {
  const { isMobile } = useSidebar();
  const pathname = useLocation({ select: (location) => location.pathname });
  const navigate = useNavigate();
  const setSessions = useChat((state) => state.setSessions);

  const updateSession = (updatedSession: Session) => {
    setSessions(
      useChat
        .getState()
        .sessions.map((session) =>
          session.id === updatedSession.id ? updatedSession : session,
        ),
    );
  };

  const handleDuplicate = async (session: Session) => {
    try {
      const duplicatedSession = await ChatService.duplicateChatSession(
        session.id,
      );
      setSessions([duplicatedSession, ...useChat.getState().sessions]);
      await navigate({ to: `/chat/${duplicatedSession.id}` });
    } catch (error) {
      showToastError(error);
    }
  };

  const handleArchive = async (session: Session) => {
    try {
      const updatedSession = await ChatService.archiveChatSession(session.id, {
        isArchive: !session.isArchive,
      });
      updateSession(updatedSession);
    } catch (error) {
      showToastError(error);
    }
  };

  const handleDelete = async (session: Session) => {
    try {
      await ChatService.deleteChatSession(session.id);
      setSessions(
        useChat.getState().sessions.filter((item) => item.id !== session.id),
      );
      if (pathname === `/chat/${session.id}`) {
        await navigate({ to: "/chat/new" });
      }
    } catch (error) {
      showToastError(error);
    }
  };

  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarGroupLabel>History</SidebarGroupLabel>
      <SidebarMenu>
        {sessions.map((item) => (
          <SidebarMenuItem key={item.id}>
            <SidebarMenuButton
              isActive={pathname === `/chat/${item.id}`}
              render={
                <Link
                  to="/chat/$id"
                  params={{
                    id: item.id,
                  }}
                />
              }
            >
              <span>{item.title}</span>
            </SidebarMenuButton>
            <DropdownMenu>
              <DropdownMenuTrigger
                render={
                  <SidebarMenuAction
                    showOnHover
                    className="aria-expanded:bg-muted"
                  />
                }
              >
                <MoreHorizontalIcon />
                <span className="sr-only">More</span>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-56 rounded-lg"
                side={isMobile ? "bottom" : "right"}
                align={isMobile ? "end" : "start"}
              >
                <DropdownMenuGroup>
                  <DropdownMenuItem
                    onClick={() => {
                      navigator.clipboard.writeText(item.id);
                    }}
                  >
                    <LinkIcon className="text-muted-foreground" />
                    <span>Copy Session ID</span>
                  </DropdownMenuItem>

                  {envConfig.current === Environment.Web ? (
                    <DropdownMenuItem
                      onClick={() => window.open(`/chat/${item.id}`)}
                    >
                      <ArrowUpRightIcon className="text-muted-foreground" />
                      <span>Open in New Tab</span>
                    </DropdownMenuItem>
                  ) : (
                    ""
                  )}

                  <DropdownMenuItem onClick={() => void handleDuplicate(item)}>
                    <CopyIcon className="text-muted-foreground" />
                    <span>Duplicate</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => void handleArchive(item)}>
                    <ArchiveIcon className="text-muted-foreground" />
                    <span>{item.isArchive ? "Unarchive" : "Archive"}</span>
                  </DropdownMenuItem>

                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="text-red-600"
                    onClick={() => void handleDelete(item)}
                  >
                    <Trash2Icon />
                    <span>Delete</span>
                  </DropdownMenuItem>
                </DropdownMenuGroup>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        ))}
        {isLoading && (
          <SidebarMenuItem className="justify-center py-2">
            <Spinner />
          </SidebarMenuItem>
        )}
      </SidebarMenu>
    </SidebarGroup>
  );
}
