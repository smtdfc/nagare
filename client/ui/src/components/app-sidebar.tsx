"use client";

import * as React from "react";

import { NavListChat } from "#/components/nav-list-chat.tsx";
import { NavMain } from "@/components/nav-main";
import { NavSecondary } from "@/components/nav-secondary";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import {
  TerminalIcon,
  SearchIcon,
  Settings2Icon,
  MessageCircleQuestionIcon,
  PlugIcon,
  MessageCircle,
} from "lucide-react";
import { useChat } from "#/hooks/use-chat.ts";
import { useCallback, useEffect, useRef, useState } from "react";
import { ChatService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";

const SESSION_PAGE_SIZE = 50;

const data = {
  teams: [
    {
      name: "Acme Inc",
      logo: <TerminalIcon />,
      plan: "Enterprise",
    },
  ],
  navMain: [
    {
      title: "New chat",
      url: "/chat/new",
      icon: <MessageCircle />,
    },
    {
      title: "Search",
      url: "#",
      icon: <SearchIcon />,
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      url: "#",
      icon: <Settings2Icon />,
    },
    {
      title: "Plugins",
      url: "#",
      icon: <PlugIcon />,
    },
    {
      title: "Help",
      url: "#",
      icon: <MessageCircleQuestionIcon />,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const sessions = useChat((c) => c.sessions);
  const setSessions = useChat((c) => c.setSessions);
  const [isLoadingSessions, setIsLoadingSessions] = useState(false);
  const sessionOffsetRef = useRef(0);
  const isLoadingSessionsRef = useRef(false);
  const reachedEndRef = useRef(false);

  const loadSessions = useCallback(
    async (append: boolean) => {
      if (isLoadingSessionsRef.current || (append && reachedEndRef.current)) {
        return;
      }

      isLoadingSessionsRef.current = true;
      setIsLoadingSessions(true);

      try {
        const loadedSessions = await ChatService.listSessions(
          SESSION_PAGE_SIZE,
          append ? sessionOffsetRef.current : 0,
        );

        if (!loadedSessions || loadedSessions.length === 0) {
          reachedEndRef.current = true;
          return;
        }

        sessionOffsetRef.current = append
          ? sessionOffsetRef.current + loadedSessions.length
          : loadedSessions.length;

        setSessions(
          append
            ? [...useChat.getState().sessions, ...loadedSessions]
            : loadedSessions,
        );
      } catch (error) {
        showToastError(error);
      } finally {
        isLoadingSessionsRef.current = false;
        setIsLoadingSessions(false);
      }
    },
    [setSessions],
  );

  useEffect(() => {
    void loadSessions(false);
  }, [loadSessions]);

  const handleSessionsScroll = useCallback(
    (event: React.UIEvent<HTMLDivElement>) => {
      const target = event.currentTarget;
      if (target.scrollHeight - target.scrollTop - target.clientHeight <= 80) {
        void loadSessions(true);
      }
    },
    [loadSessions],
  );

  return (
    <Sidebar className="border-r-0" {...props}>
      <SidebarHeader>
        <h3 className="p-4">Nagare</h3>
        <NavMain items={data.navMain} />
      </SidebarHeader>
      <SidebarContent onScroll={handleSessionsScroll}>
        <NavListChat sessions={sessions} isLoading={isLoadingSessions} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}
